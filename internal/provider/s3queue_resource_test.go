// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccS3QueueResource(t *testing.T) {
	t.Skip("This test is not supported in CI environment")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccS3QueueConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clickhouseops_s3queue.new_table", "name", "new_table"),
				),
			},
		},
	})
}

const testAccS3QueueConfig = `
resource "clickhouseops_database" "new_database" {
	name = "new_database"
}

resource "clickhouseops_s3queue" "new_table" {
  name = "new_table"
  database_name = clickhouseops_database.new_database.name
  columns = [{
	name = "a"
	type = "String"
  },{
	name = "b"
	type = "String"
  }]
  path = "s3://localhost/path/"
  format = "CSV"
  nosign = true
  settings = [{
	name = "mode"
    value = "ordered"
  }]
}
`
