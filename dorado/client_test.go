package dorado

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	baseURLTestPath = "/deviceManager/rest/xx"
)

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLTestPath+"/", http.StripPrefix(baseURLTestPath, mux))

	server := httptest.NewServer(apiHandler)

	dummyIPs := []string{"192.0.2.1:80", "192.0.2.2:80"} // not used
	client, err := NewClientDefaultToken(dummyIPs, dummyIPs, "username", "password", "portgroup", nil)
	if err != nil {
		log.Fatalf("failed to create dorado.Client: %s", err)
	}

	err = client.LocalDevice.setBaseURL(server.URL, DefaultDeviceID)
	if err != nil {
		log.Fatalf("failed to set baseURL in local devive: %s", err)
	}
	err = client.RemoteDevice.setBaseURL(server.URL, DefaultDeviceID)
	if err != nil {
		log.Fatalf("failed to set baseURL in remote devive: %s", err)
	}

	return client, mux, serverURL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}
