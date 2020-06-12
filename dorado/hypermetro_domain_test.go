package dorado

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDevice_GetHyperMetroDomains(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/HyperMetroDomain", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w,
			`
{ 
 "data": [ 
        { 
            "CPSID": "", 
            "CPSNAME": "", 
            "CPTYPE": "3", 
            "DESCRIPTION": "", 
            "STANDBYCPSID": "", 
            "STANDBYCPSNAME": "", 
            "DOMAINTYPE": "1", 
            "ID": "8038bc14bd750100", 
            "NAME": "test", 
            "REMOTEDEVICES": "[{\"devId\":\"0\",\"devESN\":\"2102350BSE10F3000088\",\"devName\":\"33aa\"}]", 
            "RUNNINGSTATUS": "1", 
            "TYPE": 15362 
        } 
 ], 
 "error": { 
        "code": 0, 
        "description": "0" 
 } 
}`)
	})

	hmds, err := client.GetHyperMetroDomains(context.Background(), nil)
	if err != nil {
		t.Errorf("GetHyperMetroDomains return err: %s", err)
	}

	want := []HyperMetroDomain{
		{
			CPSID:         "",
			CPSNAME:       "",
			CPTYPE:        "3",
			DESCRIPTION:   "",
			DOMAINTYPE:    "1",
			ID:            "8038bc14bd750100",
			NAME:          "test",
			REMOTEDEVICES: "[{\"devId\":\"0\",\"devESN\":\"2102350BSE10F3000088\",\"devName\":\"33aa\"}]",
			RUNNINGSTATUS: "1",
			TYPE:          TypeHyperMetroDomain,
		},
	}

	if !reflect.DeepEqual(hmds, want) {
		t.Errorf("GetHyperMetroDomains return %+v, want %+v", hmds, want)
	}
}
