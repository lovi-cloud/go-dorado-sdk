package dorado

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDevice_GetInitiators(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/iscsi_initiator", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w,
			`
{ 
 "data": [ 
        { 
            "HEALTHSTATUS": "1", 
            "ID": "111111111111111111", 
            "ISFREE": "false", 
            "MULTIPATHTYPE": "0", 
            "OPERATIONSYSTEM": "255", 
            "PARENTID": "0", 
            "PARENTNAME": "Host001", 
            "PARENTTYPE": 21, 
            "RUNNINGSTATUS": "28", 
            "TYPE": 222, 
            "USECHAP": "false", 
            "FAILOVERMODE": "3", 
            "SPECIALMODETYPE": "2", 
            "PATHTYPE": "1" 
        } 
 ], 
 "error": { 
        "code": 0, 
        "description": "0" 
 } 
}`)
	})

	initiators, err := client.LocalDevice.GetInitiators(context.Background(), nil)
	if err != nil {
		t.Errorf("GetInitiator return err: %s", err)
	}

	want := []Initiator{
		{
			FAILOVERMODE:    "3",
			HEALTHSTATUS:    "1",
			ID:              "111111111111111111",
			ISFREE:          "false",
			MULTIPATHTYPE:   "0",
			OPERATIONSYSTEM: "255",
			PATHTYPE:        "1",
			RUNNINGSTATUS:   "28",
			SPECIALMODETYPE: "2",
			TYPE:            TypeInitiator,
			USECHAP:         "false",
			PARENTID:        "0",
			PARENTNAME:      "Host001",
			PARENTTYPE:      TypeHost,
		},
	}

	if !reflect.DeepEqual(initiators, want) {
		t.Errorf("GetInitiators return %+v, want %+v", initiators, want)
	}
}
