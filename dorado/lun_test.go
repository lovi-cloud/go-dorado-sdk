package dorado

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDevice_GetLUNs(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/lun", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w,
			`
{ 
 "data": 
 [ 
        { 
            "ALLOCCAPACITY": "0", 
            "ALLOCTYPE": "1", 
            "CAPACITY": "2097152", 
            "CLONESOURCEID": "2", 
            "CLONESOURCENAME": "lun1", 
            "CLONESOURCETYPE": "1", 
            "COMPRESSION": "0", 
            "COMPRESSIONSAVEDCAPACITY": "0", 
            "COMPRESSIONSAVEDRATIO": "0", 
            "CREATECLONETIME": "1507732128", 
            "DEDUPSAVEDCAPACITY": "0", 
            "DEDUPSAVEDRATIO": "0", 
            "DESCRIPTION": "",
            "ENABLECOMPRESSION":"true",
            "ENABLEDEDUP":"true",
            "ENABLEISCSITHINLUNTHRESHOLD": "false", 
            "EXPOSEDTOINITIATOR": "false", 
            "EXTENDIFSWITCH": "false", 
            "HEALTHSTATUS": "1", 
            "ID": "6", 
            "IOCLASSID": "", 
            "IOPRIORITY": "1", 
            "ISADD2LUNGROUP": "false", 
            "ISCHECKZEROPAGE": "false", 
            "ISCLONE": "true", 
            "ISCSITHINLUNTHRESHOLD": "90", 
            "MIRRORPOLICY": "1", 
            "MIRRORTYPE": "0", 
            "NAME": "lun1_Clone_1710110628431", 
            "OWNINGCONTROLLER": "0B", 
            "PARENTID": "0", 
            "PARENTNAME": "sp", 
            "PREFETCHPOLICY": "0", 
            "PREFETCHVALUE": "0", 
            "REPLICATION_CAPACITY": "0", 
            "RUNNINGSTATUS": "27", 
            "RUNNINGWRITEPOLICY": "1", 
            "SECTORSIZE": "512",
            "SUBTYPE": "0", 
            "TOTALSAVEDCAPACITY": "0", 
            "TOTALSAVEDRATIO": "0", 
            "TYPE": 11, 
            "USAGETYPE": "0", 
            "WORKLOADTYPENAME": "", 
            "WRITEPOLICY": "1" 
        } 
 ], 
 "error": { 
        "code": 0, 
        "description": "0" 
 } 
}`)
	})

	luns, err := client.LocalDevice.GetLUNs(context.Background(), nil)
	if err != nil {
		t.Errorf("GetLUNs return err: %s", err)
	}

	want := []LUN{
		{
			ALLOCCAPACITY:               "0",
			ALLOCTYPE:                   "1",
			CAPACITY:                    2097152,
			COMPRESSION:                 "0",
			COMPRESSIONSAVEDCAPACITY:    "0",
			COMPRESSIONSAVEDRATIO:       "0",
			DEDUPSAVEDCAPACITY:          "0",
			DEDUPSAVEDRATIO:             "0",
			DESCRIPTION:                 "",
			ENABLECOMPRESSION:           "true",
			ENABLEISCSITHINLUNTHRESHOLD: "false",
			EXPOSEDTOINITIATOR:          "false",
			EXTENDIFSWITCH:              "false",
			HEALTHSTATUS:                "1",
			ID:                          6,
			IOCLASSID:                   "",
			IOPRIORITY:                  "1",
			ISADD2LUNGROUP:              false,
			ISCHECKZEROPAGE:             "false",
			ISCLONE:                     true,
			ISCSITHINLUNTHRESHOLD:       "90",
			MIRRORPOLICY:                "1",
			MIRRORTYPE:                  "0",
			NAME:                        "lun1_Clone_1710110628431",
			OWNINGCONTROLLER:            "0B",
			PARENTID:                    0,
			PARENTNAME:                  "sp",
			PREFETCHPOLICY:              "0",
			PREFETCHVALUE:               "0",
			REPLICATIONCAPACITY:         "0",
			RUNNINGSTATUS:               "27",
			RUNNINGWRITEPOLICY:          "1",
			SECTORSIZE:                  "512",
			SNAPSHOTIDS:                 "",
			SNAPSHOTSCHEDULEID:          "",
			SUBTYPE:                     "0",
			THINCAPACITYUSAGE:           "",
			TOTALSAVEDCAPACITY:          "0",
			TOTALSAVEDRATIO:             "0",
			TYPE:                        TypeLUN,
			USAGETYPE:                   "0",
			WORKINGCONTROLLER:           "",
			WORKLOADTYPEID:              "",
			WORKLOADTYPENAME:            "",
			WRITEPOLICY:                 "1",
			WWN:                         "",
			HyperCdpScheduleID:          "",
			LunCgID:                     "",
			RemoteLunWwn:                "",
			TakeOverLunWwn:              "",
		},
	}

	if !reflect.DeepEqual(luns, want) {
		t.Errorf("GetLUNs return %+v, want %+v", luns, want)
	}
}
