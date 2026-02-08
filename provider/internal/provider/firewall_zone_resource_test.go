package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirewallZoneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "linux_firewall_zone" "test" {
				  name        = "testzone"
				  description = "Test zone"
				  target      = "default"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linux_firewall_zone.test", "name", "testzone"),
					resource.TestCheckResourceAttr("linux_firewall_zone.test", "description", "Test zone"),
					resource.TestCheckResourceAttr("linux_firewall_zone.test", "target", "default"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "linux_firewall_zone.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
