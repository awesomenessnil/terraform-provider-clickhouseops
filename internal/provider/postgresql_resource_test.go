// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPostgreSQLResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPostgreSQLResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clickhouseops_postgresql.new_table", "name", "new_table"),
				),
			},
		},
	})
}

const (
	testAccPostgreSQLResourceConfig = `
resource "clickhouseops_database" "new_database" {
	name = "new_database"
}

resource "clickhouseops_postgresql" "new_table" {
	name = "new_table"
	database_name = clickhouseops_database.new_database.name
	columns = [{
		name = "a"
		type = "String"
	}]
	postgresql_host = "localhost"
	postgresql_port = 5432
	postgresql_database_name = "postgres"
	postgresql_table_name = "test"
	postgresql_username = "user"
    postgresql_password = "password"
}
`
)
