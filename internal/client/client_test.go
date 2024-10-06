package client

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setup(t *testing.T) (*http.ServeMux, *httptest.Server, *Client) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := NewClient("test", server.URL)

	return mux, server, client
}

func teardown(server *httptest.Server) {
	server.Close()
}

func assertRequest(t *testing.T, r *http.Request, method string) {
	t.Helper()

	if r.Method != method {
		t.Errorf("unexpected method: %s", r.Method)
	}
	if !strings.Contains(r.UserAgent(), "terraform-provider-pubsubhubbub") {
		t.Errorf("unexpected User-Agent: %s", r.UserAgent())
	}
}

func TestNewClient(t *testing.T) {
	// can instantiate without error
	var _ *Client = NewClient("test", "hub.example.com")
}

func TestClient_Subscribe(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assertRequest(t, r, "POST")

		w.WriteHeader(http.StatusAccepted)
	})

	err := client.Subscribe("http://example.com/callback", "http://example.com/topic", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestClient_Unsubscribe(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assertRequest(t, r, "POST")

		w.WriteHeader(http.StatusAccepted)
	})

	err := client.Unsubscribe("http://example.com/callback", "http://example.com/topic", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}
