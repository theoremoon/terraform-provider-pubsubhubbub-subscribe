package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSubscribe(t *testing.T) {
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
