package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func setup(t *testing.T) (*http.ServeMux, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	return mux, server
}

func TestSubscribe(t *testing.T) {
	mux, server := setup(t)
	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	providerConfig := fmt.Sprintf(`
provider "pubsubhubbub-subscribe" {
  hub_url = "%s"
}`, server.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					resource "pubsubhubbub-subscribe_subscribe" "tf-test2" {
					  callback_url = "https://example.com/callback"
					  topic_url = "https://example.com/topic"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pubsubhubbub-subscribe_subscribe.tf-test2", "callback_url", "https://example.com/callback"),
					resource.TestCheckResourceAttr("pubsubhubbub-subscribe_subscribe.tf-test2", "topic_url", "https://example.com/topic"),
				),
			},
			{
				// cleanup
				Config: providerConfig,
			},
		},
	})
}
