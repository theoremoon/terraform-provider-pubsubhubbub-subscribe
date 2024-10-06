package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "pubsubhubbub-subscribe" {
  hub_url = "https://pubsubhubbub.example.com/"
}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"pubsubhubbub-subscribe": providerserver.NewProtocol6WithError(New("test")()),
	}
)
