package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	client *http.Client
	hubURL string
}

func NewClient(version, hubURL string) *Client {
	return &Client{
		client: &http.Client{
			Transport: newTransport(version),
		},
		hubURL: hubURL,
	}
}

func (c *Client) subscriptionRequest(mode, callbackURL, topicURL, hmacSecret string) error {
	data := url.Values{}
	data.Set("hub.callback", callbackURL)
	data.Set("hub.mode", mode)
	data.Set("hub.topic", topicURL)
	data.Set("hub.secret", hmacSecret)

	req, err := http.NewRequest("POST", c.hubURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// The hub MUST respond to a subscription request with an HTTP [RFC2616] 202 "Accepted" response to indicate that the request was received
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) Subscribe(callbackURL, topicURL, hmacSecret string) error {
	return c.subscriptionRequest("subscribe", callbackURL, topicURL, hmacSecret)
}

func (c *Client) Unsubscribe(callbackURL, topicURL, hmacSecret string) error {
	return c.subscriptionRequest("unsubscribe", callbackURL, topicURL, hmacSecret)
}
