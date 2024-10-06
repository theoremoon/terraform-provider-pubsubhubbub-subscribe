resource "pubsubhubbub-subscribe_subscribe" "example" {
  topic_url = "https://example.com/topic"
  callback_url = "https://example.com/callback"
  secret = "123"
}
