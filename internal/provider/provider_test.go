package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/kyswtn/terraform-provider-porkbun/internal/mockbun"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"porkbun": providerserver.NewProtocol6WithError(New("test")()),
}

func getProviderConfigWithMockServer(t *testing.T) (string, *mockbun.Server) {
	t.Helper()

	mockbunServer := mockbun.New()
	t.Cleanup(mockbunServer.Close)

	config := fmt.Sprintf(`
    provider "porkbun" {
        api_key         = "apikey"
        secret_api_key  = "secretapikey"
        custom_base_url = "%s"
    }
    `, mockbunServer.URL)

	return config, mockbunServer
}
