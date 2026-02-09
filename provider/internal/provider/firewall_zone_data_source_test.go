package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirewallZoneDataSource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' is set")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "linux_firewall_zone" "testzone" {
				  name = "testzone"
				  description = "Test zone for acceptance tests"
				  target = "default"
				}
				`,
			},
			{
				Config: providerConfig + `
				resource "linux_firewall_zone" "testzone" {
				  name = "testzone"
				  description = "Test zone for acceptance tests"
				  target = "default"
				}

				data "linux_firewall_zones" "zones" {}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.linux_firewall_zones.zones", "zones.#"),
				),
			},
		},
	})
}
