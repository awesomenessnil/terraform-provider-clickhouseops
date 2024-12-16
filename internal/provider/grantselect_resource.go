// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/katzucurry/terraform-provider-clickhouseops/internal/common"
)

var (
	_ resource.Resource                = &GrantSelect{}
	_ resource.ResourceWithConfigure   = &GrantSelect{}
	_ resource.ResourceWithImportState = &GrantSelect{}
)

func NewGrantSelect() resource.Resource {
	return &GrantSelect{}
}

type GrantSelect struct {
	db *sql.DB
}

type GrantSelectModel struct {
	ID           types.String   `tfsdk:"id"`
	DatabaseName types.String   `tfsdk:"database_name"`
	TableName    types.String   `tfsdk:"table_name"`
	ColumnsName  []types.String `tfsdk:"columns_name"`
	ClusterName  types.String   `tfsdk:"cluster_name"`
	Assignee     types.String   `tfsdk:"assignee"`
}

func (r *GrantSelect) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_grantselect"
}

func (r *GrantSelect) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Clickhouse grant select privilige to a user or a role (Assignee)",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"database_name": schema.StringAttribute{
				MarkdownDescription: "Name of the database where table you want to grant select permissions is located",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"table_name": schema.StringAttribute{
				MarkdownDescription: "Name of the table you want to grant select permissions",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"columns_name": schema.ListAttribute{
				MarkdownDescription: "List of columns is the user or role restricted on the target table",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"cluster_name": schema.StringAttribute{
				MarkdownDescription: "Clickhouse cluster name",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"assignee": schema.StringAttribute{
				MarkdownDescription: "User or Role you want grant permissions",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *GrantSelect) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	db, ok := req.ProviderData.(*sql.DB)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sql.DB, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.db = db
}

/*
	Clickhouse Grant Syntax for reference

GRANT [ON CLUSTER cluster_name] privilege[(column_name [,...])] [,...] ON {db.table|db.*|*.*|table|*} TO {user | role | CURRENT_USER} [,...] [WITH GRANT OPTION] [WITH REPLACE OPTION].
*/
const ddlCreateGrantSelectTemplate = `
GRANT {{if not .ClusterName.IsNull}} ON CLUSTER '{{.ClusterName.ValueString}}' {{end}}SELECT{{$size := size .ColumnsName}}{{with .ColumnsName}}({{range $i, $e := .}}"{{$e.ValueString}}"{{if lt $i $size}},{{end}}{{end}}){{end}} ON "{{.DatabaseName.ValueString}}".{{if not .TableName.IsNull}}"{{.TableName.ValueString}}"{{else}}*{{end}} TO '{{.Assignee.ValueString}}'
`

func (r *GrantSelect) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *GrantSelectModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query, err := common.RenderTemplate(ddlCreateGrantSelectTemplate, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Granting Permissions",
			"Could not render DDL, unexpected error: "+err.Error(),
		)
		return
	}

	_, err = r.db.Exec(*query)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Granting Permissions",
			"Could not execute DDL, unexpected error: "+*query+err.Error(),
		)
		return
	}
	var columns []string
	for _, v := range data.ColumnsName {
		columns = append(columns, v.ValueString())
	}
	data.ID = types.StringValue(data.ClusterName.ValueString() + ":" + strings.Join(columns, ":") + ":" + data.Assignee.ValueString())

	tflog.Trace(ctx, "Created a GrantSelect Resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *GrantSelect) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *GrantSelectModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *GrantSelect) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *GrantSelectModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

const ddlDestroyGrantSelectTemplate = `
REVOKE {{if not .ClusterName.IsNull}} ON CLUSTER '{{.ClusterName.ValueString}}' {{end}}SELECT{{$size := size .ColumnsName}}{{with .ColumnsName}}({{range $i, $e := .}}"{{$e.ValueString}}"{{if lt $i $size}},{{end}}{{end}}){{end}} ON "{{.DatabaseName.ValueString}}".{{if not .TableName.IsNull}}"{{.TableName.ValueString}}"{{else}}*{{end}} FROM '{{.Assignee.ValueString}}'
`

func (r *GrantSelect) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *GrantSelectModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	query, err := common.RenderTemplate(ddlDestroyGrantSelectTemplate, data)
	if err != nil {
		resp.Diagnostics.AddError("", ""+err.Error())
		return
	}

	_, err = r.db.Exec(*query)
	if err != nil {
		resp.Diagnostics.AddError("", ""+err.Error())
		return
	}
}

func (r *GrantSelect) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
