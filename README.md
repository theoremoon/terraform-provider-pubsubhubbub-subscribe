# terraform-provider-pubsubhubbub-subscribe

A terraform provider to register subscription to PubSubHubbub.

## example

```hcl
terraform {
  required_providers {
    pubsubhubbub-subscribe = {
    }
  }
}

provider "pubshubhubbub-subscribe" { 
  hub_url = "https://pubsubhubbub.appspot.com/"
}

resource "pubsubhubbub-subscribe_subscribe" "example" {
  topic_url = "https://example.com/feed"
  callback_url = "https://example.com/callback
  secret = "123456"
}
```
