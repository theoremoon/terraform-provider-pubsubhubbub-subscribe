package client

import (
	"fmt"
	"net/http"

	"github.com/motemen/go-wsse"
)

// transport wraps wsse.Transport to set X-WSSE header.
// Additionally, it sets User-Agent header
type transport struct {
	Transport wsse.Transport

	version string
}

func newTransport(version string) *transport {
	return &transport{
		version: version,
	}
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ua := fmt.Sprintf("terraform-provider-pubsubhubbub/%s", t.version)
	req.Header.Set("User-Agent", ua)

	return t.Transport.RoundTrip(req)
}
