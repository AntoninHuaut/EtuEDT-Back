package server

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/AntoninHuaut/EtuEDT-Back/internal/config"
)

func TestSourceOverridePublicContract(t *testing.T) {
	cfg := testConfig()
	cfg.Universities[0].SourceOverride = &config.SourceOverrideConfig{Provider: "reparser", BaseURL: "https://private-upstream.example.test"}
	srv := newTestServer(t, cfg)
	defer srv.Close()
	for _, path := range []string{"/v3/univs", "/v3/univs/1", "/v3/univs/1/rooms", "/v3/univs/1/rooms/10", "/openapi.json"} {
		resp := get(t, srv, path)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("%s: status %d", path, resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		for _, hidden := range []string{"sourceOverride", "private-upstream", "reparser", "experimental"} {
			if strings.Contains(strings.ToLower(string(body)), strings.ToLower(hidden)) {
				t.Errorf("%s exposes %s", path, hidden)
			}
		}
	}
	resp := get(t, srv, "/v3/univs/1/rooms/999/events")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("unknown room: %d", resp.StatusCode)
	}
	resp = get(t, srv, "/v3/univs/1/free-rooms")
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("free rooms: %d", resp.StatusCode)
	}
}
