package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestNameserversDataSource(t *testing.T) {
	providerConfig, mockbun := getProviderConfigWithMockServer(t)
	mockbun.SetNameservers("example.com", []string{
		"evan.ns.cloudflare.com",
		"sandy.ns.cloudflare.com",
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "porkbun_nameservers" "test" {
						domain = "example.com"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.porkbun_nameservers.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("data.porkbun_nameservers.test", "nameservers.0", "evan.ns.cloudflare.com"),
					resource.TestCheckResourceAttr("data.porkbun_nameservers.test", "nameservers.1", "sandy.ns.cloudflare.com"),
				),
			},
		},
	})
}
