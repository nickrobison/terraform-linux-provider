package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZpoolDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "linux_zpool" "pool1" {
				  name = "tank"
				}
				`,
			},
			{
				Config: providerConfig + `
				data "linux_zpools" "pools" {}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zpools.pools", "zpools.#", "1"),
				),
			},
		},
	})
}
