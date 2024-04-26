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
					resource.TestCheckResourceAttr("porkbun_dns_record.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test", "type", "A"),
					resource.TestCheckResourceAttr("porkbun_dns_record.test", "content", "1.2.3.4"),
				),
			},
			// Test update and read with subdomain.
			// TODO: How would I pass ID from previous porkbun_dns_record.test into the new block so that I can test?
			// {
			// 	Config: providerConfig + `
			//                  resource "porkbun_dns_record" "test" {
			//                      domain = "example.com"
			//                               name = "www"
			//                      type = "CNAME"
			//                      content = "4.3.2.1"
			//                  }
			//              `,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("porkbun_dns_record.test", "domain", "example.com"),
			// 		resource.TestCheckResourceAttr("porkbun_dns_record.test", "name", "www"),
			// 		resource.TestCheckResourceAttr("porkbun_dns_record.test", "type", "CNAME"),
			// 		resource.TestCheckResourceAttr("porkbun_dns_record.test", "content", "4.3.2.1"),
			// 	),
			// },
		},
	})
}
