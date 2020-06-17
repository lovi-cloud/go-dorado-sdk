package dorado

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDevice_GetPortGroups(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/portgroup", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w,
			`
{ 
 "data": [ 
        { 
            "DESCRIPTION": "", 
            "ID": "0", 
            "NAME": "PortGroup001", 
            "TYPE": 257 
        } 
 ], 
 "error": { 
        "code": 0, 
        "description": "0" 
 } 
}`)
	})

	portgroups, err := client.LocalDevice.GetPortGroups(context.Background(), nil)
	if err != nil {
		t.Errorf("GetPortGroups return err: %s", err)
	}

	want := []PortGroup{
		{
			DESCRIPTION: "",
			ID:          0,
			NAME:        "PortGroup001",
			TYPE:        TypePortGroup,
		},
	}

	if !reflect.DeepEqual(portgroups, want) {
		t.Errorf("GetPortGroups return %+v, want %+v", portgroups, want)
	}
}
