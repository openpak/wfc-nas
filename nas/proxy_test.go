package nas

import (
	"io"
	"net/http"
	"net/http/httptest"
	"openpak/wfc-common/service"
	"strings"
	"testing"
)

// Exercise native duplicate-Host framing and streamed POST bodies through both
// compatibility and NAS hops, using the same keepalive connection across hosts.
func TestForwardingAcrossTwoHopsPreservesPostBody(t *testing.T) {
	const body = `<Envelope><Body>native SOAP data</Body></Envelope>`
	key := strings.Repeat("k", 32)
	t.Setenv("WFC_INTERNAL_KEY", key)
	owner := httptest.NewServer(service.Guard(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			raw, err := io.ReadAll(r.Body)
			if err != nil || string(raw) != body {
				t.Errorf("body changed: %q %v", raw, err)
			}
			w.Write(raw)
			return
		}
		http.NotFound(w, r)
	})))
	defer owner.Close()
	proxy, err := service.Proxy(owner.URL, key)
	if err != nil {
		t.Fatal(err)
	}
	nasServer := httptest.NewUnstartedServer(proxy)
	nasServer.Listener = &httpListener{Listener: nasServer.Listener}
	nasServer.Start()
	defer nasServer.Close()
	forward, err := service.Proxy(nasServer.URL, key)
	if err != nil {
		t.Fatal(err)
	}
	gateway := httptest.NewServer(forward)
	defer gateway.Close()
	client := gateway.Client()
	for i := 0; i < 10; i++ {
		for _, host := range []string{"sake.gs.nintendowifi.net", "race.gs.nintendowifi.net", "gamestats.gs.nintendowifi.net", "sake.gs.nintendowifi.net"} {
			method := "GET"
			payload := ""
			if host == "sake.gs.nintendowifi.net" {
				method = "POST"
				payload = body
			}
			request, err := http.NewRequest(method, gateway.URL+"/wire", strings.NewReader(payload))
			if err != nil {
				t.Fatal(err)
			}
			request.Host = host
			response, err := client.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			raw, err := io.ReadAll(response.Body)
			response.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			if method == "POST" && (response.StatusCode != 200 || string(raw) != body) {
				t.Fatalf("POST reply: %d %q", response.StatusCode, raw)
			}
		}
	}
}
