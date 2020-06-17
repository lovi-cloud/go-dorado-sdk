package dorado

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDevice_GetLunGroups(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/lungroup", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w,
			`
{ 
 "data": [ 
        { 
            "CAPCITY": "2097152", 
            "DESCRIPTION": "", 
            "ID": "0", 
            "ISADD2MAPPINGVIEW": "true", 
            "NAME": "LUNGroup001",
            "TYPE": 256 
        } 
 ], 
 "error": { 
        "code": 0, 
        "description": "0" 
 } 
}`)
	})

	lungroups, err := client.LocalDevice.GetLunGroups(context.Background(), nil)
	if err != nil {
		t.Errorf("GetLunGroups return err: %s", err)
	}

	want := []LunGroup{
		{
			CAPCITY:            "2097152",
			DESCRIPTION:        "",
			ID:                 0,
			ISADD2MAPPINGVIEW:  true,
			NAME:               "LUNGroup001",
			SMARTQOSPOLICYID:   "",
			TYPE:               TypeLUNGroup,
			ASSOCIATELUNIDLIST: "",
		},
	}

	if !reflect.DeepEqual(lungroups, want) {
		t.Errorf("GetLunGroups return %+v, want %+v", lungroups, want)
	}
}
