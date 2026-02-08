package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirewallRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test rich rule
			{
				Config: providerConfig + `
				resource "linux_firewall_zone" "testzone" {
				  name = "testzone"
				}

				resource "linux_firewall_rule" "test_rich" {
				  zone      = linux_firewall_zone.testzone.name
				  rule_type = "rich"
				  rule      = "rule family=ipv4 source address=192.168.1.0/24 accept"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linux_firewall_rule.test_rich", "zone", "testzone"),
					resource.TestCheckResourceAttr("linux_firewall_rule.test_rich", "rule_type", "rich"),
					resource.TestCheckResourceAttr("linux_firewall_rule.test_rich", "rule", "rule family=ipv4 source address=192.168.1.0/24 accept"),
				),
			},
			// Test port rule
			{
				Config: providerConfig + `
				resource "linux_firewall_zone" "testzone" {
				  name = "testzone"
				}

				resource "linux_firewall_rule" "test_port" {
				  zone      = linux_firewall_zone.testzone.name
				  rule_type = "port"
				  port      = "8080"
				  protocol  = "tcp"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linux_firewall_rule.test_port", "zone", "testzone"),
					resource.TestCheckResourceAttr("linux_firewall_rule.test_port", "rule_type", "port"),
					resource.TestCheckResourceAttr("linux_firewall_rule.test_port", "port", "8080"),
					resource.TestCheckResourceAttr("linux_firewall_rule.test_port", "protocol", "tcp"),
				),
			},
			// Test service rule
			{
				Config: providerConfig + `
				resource "linux_firewall_zone" "testzone" {
				  name = "testzone"
				}

				resource "linux_firewall_rule" "test_service" {
				  zone      = linux_firewall_zone.testzone.name
				  rule_type = "service"
				  service   = "http"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("linux_firewall_rule.test_service", "zone", "testzone"),
					resource.TestCheckResourceAttr("linux_firewall_rule.test_service", "rule_type", "service"),
					resource.TestCheckResourceAttr("linux_firewall_rule.test_service", "service", "http"),
				),
			},
		},
	})
}
