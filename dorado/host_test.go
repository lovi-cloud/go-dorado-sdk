package dorado

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDevice_GetHosts(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/host", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w,
			`
{ 
 "data": [ 
        { 
            "DESCRIPTION": "", 
            "HEALTHSTATUS": "1", 
            "ID": "0", 
            "INITIATORNUM": "0", 
            "IP": "", 
            "ISADD2HOSTGROUP": "true", 
            "LOCATION": "", 
            "MODEL": "", 
            "NAME": "Host001", 
            "NETWORKNAME": "", 
            "OPERATIONSYSTEM": "0", 
            "PARENTID": "1", 
            "PARENTNAME": "hostgroup1", 
            "PARENTTYPE": 14, 
            "RUNNINGSTATUS": "1", 
            "TYPE": 21 
        }, 
        { 
            "DESCRIPTION": "", 
            "HEALTHSTATUS": "1", 
            "ID": "1", 
            "INITIATORNUM": "0", 
            "IP": "", 
            "ISADD2HOSTGROUP": "false", 
            "LOCATION": "", 
            "MODEL": "", 
            "NAME": "Host002", 
            "NETWORKNAME": "", 
            "OPERATIONSYSTEM": "0", 
            "RUNNINGSTATUS": "1", 
            "TYPE": 21 
        } 
 ], 
 "error": { 
        "code": 0, 
        "description": "0" 
 } 
}`)
	})

	hosts, err := client.LocalDevice.GetHosts(context.Background(), nil)
	if err != nil {
		t.Errorf("GetHosts return err: %s", err)
	}

	want := []Host{
		{
			DESCRIPTION:     "",
			HEALTHSTATUS:    "1",
			ID:              0,
			INITIATORNUM:    "0",
			IP:              "",
			ISADD2HOSTGROUP: true,
			LOCATION:        "",
			MODEL:           "",
			NAME:            "Host001",
			NETWORKNAME:     "",
			OPERATIONSYSTEM: "0",
			PARENTID:        "1",
			PARENTNAME:      "hostgroup1",
			PARENTTYPE:      TypeHostGroup,
			RUNNINGSTATUS:   "1",
			TYPE:            TypeHost,
		},
		{
			DESCRIPTION:     "",
			HEALTHSTATUS:    "1",
			ID:              1,
			INITIATORNUM:    "0",
			IP:              "",
			ISADD2HOSTGROUP: false,
			LOCATION:        "",
			MODEL:           "",
			NAME:            "Host002",
			NETWORKNAME:     "",
			OPERATIONSYSTEM: "0",
			PARENTID:        "",
			PARENTNAME:      "",
			PARENTTYPE:      0,
			RUNNINGSTATUS:   "1",
			TYPE:            TypeHost,
		},
	}

	if !reflect.DeepEqual(hosts, want) {
		t.Errorf("GetHosts return %+v, want %+v", hosts, want)
	}
}
