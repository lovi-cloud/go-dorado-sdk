package dorado

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDevice_GetHostGroups(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/hostgroup", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w,
			`
{ 
 "data": [ 
        { 
            "DESCRIPTION": "", 
            "ID": "0", 
            "ISADD2MAPPINGVIEW": "false", 
            "NAME": "hostgroup1", 
            "TYPE": 14 
        }, 
        { 
            "DESCRIPTION": "", 
            "ID": "1", 
            "ISADD2MAPPINGVIEW": "false", 
            "NAME": "HostGroup002", 
            "TYPE": 14 
        } 
 ], 
 "error": { 
        "code": 0, 
        "description": "0" 
 } 
}`)
	})

	hostgroups, err := client.LocalDevice.GetHostGroups(context.Background(), nil)
	if err != nil {
		t.Errorf("GetHostGroups return err: %s", err)
	}

	want := []HostGroup{
		{
			DESCRIPTION:       "",
			ID:                0,
			ISADD2MAPPINGVIEW: false,
			NAME:              "hostgroup1",
			TYPE:              TypeHostGroup,
		},
		{
			DESCRIPTION:       "",
			ID:                1,
			ISADD2MAPPINGVIEW: false,
			NAME:              "HostGroup002",
			TYPE:              TypeHostGroup,
		},
	}

	if !reflect.DeepEqual(hostgroups, want) {
		t.Errorf("GetHostGroups return %+v, want %+v", hostgroups, want)
	}
}
