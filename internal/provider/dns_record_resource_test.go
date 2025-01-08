package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDNSRecordResource(t *testing.T) {
	providerConfig, _ := getProviderConfigWithMockServer(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test create and read.
			{
				Config: providerConfig + `
					resource "porkbun_dns_record" "test" {
						domain = "example.com"
						type = "A"
						content = "1.2.3.4"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_dns_record.test", "id", "2010068919648951531"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test", "type", "A"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test", "content", "1.2.3.4"),
				),
			},
		},
	})
}
