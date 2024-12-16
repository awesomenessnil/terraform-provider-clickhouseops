// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMergeTreeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "clickhouseops_database" "new_database" {
	name = "new_db"
	comment = "new db test comment"
}		

resource "clickhouseops_mergetree" "new_table1" {
	name = "test_merge_tree"
	database_name = clickhouseops_database.new_database.name
	columns = [{
		name = "a"
		type = "String"
	},{
		name = "b"
		type = "String"
	}]
	order_by = ["a"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clickhouseops_mergetree.new_table1", "name", "test_merge_tree"),
					resource.TestCheckResourceAttr("clickhouseops_mergetree.new_table1", "order_by.0", "a"),
				),
			},
		},
	})
}
