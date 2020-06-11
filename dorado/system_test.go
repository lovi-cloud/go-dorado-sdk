package dorado

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDevice_GetSystem(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/system/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w,
			`
{
  "data": {
    "ID": "TEST-DORADO-ID",
    "NAME": "test-dorado",
    "TYPE": 201
  },
  "error": {
    "code": 0,
    "description": "0"
  }
}`)
	})

	system, err := client.LocalDevice.GetSystem(context.Background())
	if err != nil {
		t.Errorf("GetSystem return err: %v", err)
	}

	want := &System{
		ID:   "TEST-DORADO-ID",
		NAME: "test-dorado",
		TYPE: 201,
	}

	if !reflect.DeepEqual(system, want) {
		t.Errorf("GetSystem return %+v, want %+v", system, want)
	}
}
