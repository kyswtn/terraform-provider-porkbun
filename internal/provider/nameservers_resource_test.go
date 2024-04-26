package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestNameserversResource(t *testing.T) {
	providerConfig, mockbun := getProviderConfigWithMockServer(t)
	mockbun.SetNameservers("example.com", []string{})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test create, read and update.
			{
				Config: providerConfig + `
					resource "porkbun_nameservers" "test" {
						domain = "example.com"
						nameservers = [
							"evan.ns.cloudflare.com",
							"sandy.ns.cloudflare.com"
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_nameservers.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("porkbun_nameservers.test", "nameservers.0", "evan.ns.cloudflare.com"),
					resource.TestCheckResourceAttr("porkbun_nameservers.test", "nameservers.1", "sandy.ns.cloudflare.com"),
				),
			},
			// Test import.
			{
				ResourceName:                         "porkbun_nameservers.test",
				ImportStateId:                        "example.com",
				ImportStateVerifyIdentifierAttribute: "domain",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}
