package dorado

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDevice_GetMappingViews(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/mappingview", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w,
			`
{ 
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
