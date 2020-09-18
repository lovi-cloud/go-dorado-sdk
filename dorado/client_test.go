package dorado

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
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

	dummyIPs := []string{server.URL, server.URL}
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

// TestDevice_UnAuthorizedRetry test retry function
func TestDevice_UnAuthorizedRetry(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	responses := []func(w http.ResponseWriter, r *http.Request){
		func(w http.ResponseWriter, r *http.Request) {
			// response auth error
			fmt.Fprintln(w, `{
  "data": [],
  "error": {
    "code": -401,
    "description": "This operation fails to be performed because of the unauthorized REST.",
	"suggestion": "Before performing this operation, ensure that REST is authorized."
  }
}`)
		},
		func(w http.ResponseWriter, r *http.Request) {
			// response auth error
			fmt.Fprintln(w, `{
  "data": [],
  "error": {
    "code": -401,
    "description": "This operation fails to be performed because of the unauthorized REST.",
	"suggestion": "Before performing this operation, ensure that REST is authorized."
  }
}`)
		},
		func(w http.ResponseWriter, r *http.Request) {
			// response correct error (token refreshed)
			fmt.Fprintln(w, `{
  "data": [
       {
           "DESCRIPTION": "",
           "ENABLEINBANDCOMMAND": "true",
           "ID": "1",
           "INBANDLUNWWN": "",
           "NAME": "MappingView001",
           "TYPE": 245
       }
  ],
  "error": {
       "code": 0,
       "description": "0"
  }
}`)
		},
	}

	mux.HandleFunc("/sessions", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,
			`{
  "data": {
    "iBaseToken": "dummy_token",
    "deviceid": "xx"
  },
  "error": {
    "code": 0,
    "description": "0"
  }
}`)
	})

	responseCount := 0
	mux.HandleFunc("/mappingview", func(w http.ResponseWriter, r *http.Request) {
		responses[responseCount](w, r)
		responseCount++
	})

	mappingviews, err := client.LocalDevice.GetMappingViews(context.Background(), nil)
	if err != nil {
		t.Errorf("GetMappingViews return err: %s", err)
	}

	want := []MappingView{
		{
			DESCRIPTION:         "",
			ENABLEINBANDCOMMAND: true,
			ID:                  1,
			INBANDLUNWWN:        "",
			NAME:                "MappingView001",
			TYPE:                TypeMappingView,
		},
	}

	if !reflect.DeepEqual(mappingviews, want) {
		t.Errorf("GetMappingViews return %+v, want %+v", mappingviews, want)
	}
}

// TestDevice_UnAuthorizedRetryFailed test change DefaultHTTPRetryCount
func TestDevice_UnAuthorizedRetryFailed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	DefaultHTTPRetryCount = 1
	defer func() {
		DefaultHTTPRetryCount = 10
	}()

	responses := []func(w http.ResponseWriter, r *http.Request){
		func(w http.ResponseWriter, r *http.Request) {
			// response auth error
			fmt.Fprintln(w, `{
  "data": [],
  "error": {
    "code": -401,
    "description": "This operation fails to be performed because of the unauthorized REST.",
	"suggestion": "Before performing this operation, ensure that REST is authorized."
  }
}`)
		},
		func(w http.ResponseWriter, r *http.Request) {
			// response auth error
			fmt.Fprintln(w, `{
  "data": [],
  "error": {
    "code": -401,
    "description": "This operation fails to be performed because of the unauthorized REST.",
	"suggestion": "Before performing this operation, ensure that REST is authorized."
  }
}`)
		},
		func(w http.ResponseWriter, r *http.Request) {
			// response correct error (token refreshed)
			fmt.Fprintln(w, `{
  "data": [
       {
           "DESCRIPTION": "",
           "ENABLEINBANDCOMMAND": "true",
           "ID": "1",
           "INBANDLUNWWN": "",
           "NAME": "MappingView001",
           "TYPE": 245
       }
  ],
  "error": {
       "code": 0,
       "description": "0"
  }
}`)
		},
	}

	mux.HandleFunc("/sessions", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,
			`{
  "data": {
    "iBaseToken": "dummy_token",
    "deviceid": "xx"
  },
  "error": {
    "code": 0,
    "description": "0"
  }
}`)
	})

	responseCount := 0
	mux.HandleFunc("/mappingview", func(w http.ResponseWriter, r *http.Request) {
		responses[responseCount](w, r)
		responseCount++
	})

	mappingviews, err := client.LocalDevice.GetMappingViews(context.Background(), nil)
	if err == nil || mappingviews != nil {
		t.Errorf("GetMappingViews must return err, but err is nil")
	}

	if !errors.Is(err, ErrUnAuthorized) {
		t.Errorf("GetMappingViews must return err: %+v, but return err: %+v", ErrUnAuthorized, err)
	}
}
