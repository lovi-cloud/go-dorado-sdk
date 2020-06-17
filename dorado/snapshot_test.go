package dorado

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDevice_GetSnapshots(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/snapshot", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w,
			`
{ 
  "data": [ 
 { 
            "CASCADEDLEVEL": "0", 
            "CASCADEDNUM": "1", 
            "CONSUMEDCAPACITY": "0", 
            "DESCRIPTION": "", 
            "EXPOSEDTOINITIATOR": "false", 
            "HEALTHSTATUS": "1", 
            "ID": "12", 
            "IOCLASSID": "", 
            "IOPRIORITY": "1", 
            "SOURCELUNCAPACITY":"2097152", 
            "ISSCHEDULEDSNAP":"0", 
            "NAME": "gz0000_Snap_1707251229076", 
            "PARENTID": "7", 
            "PARENTNAME": "gz0000", 
            "PARENTTYPE": 11, 
            "ROLLBACKENDTIME": "-1", 
            "ROLLBACKRATE": "-1", 
            "ROLLBACKSPEED": "-1", 
            "ROLLBACKSTARTTIME": "-1", 
            "ROLLBACKTARGETOBJID": "4294967295", 
            "ROLLBACKTARGETOBJNAME": "--", 
            "RUNNINGSTATUS": "43", 
            "SOURCELUNID": "7", 
            "SOURCELUNNAME": "gz0000", 
            "SUBTYPE": "0", 
            "TIMESTAMP": "1501013445", 
            "TYPE": 27, 
            "USERCAPACITY": "2097152", 
            "WORKINGCONTROLLER": "0B", 
            "WWN": "6a1b2c3100f4d5e600089fd70000000c", 
            "replicationCapacity": "0", 
            "WORKLOADTYPEID": "14", 
            "WORKLOADTYPENAME": "Databases" 
        }, { 
            "CASCADEDLEVEL": "1", 
            "CASCADEDNUM": "0", 
            "CONSUMEDCAPACITY": "0", 
            "DESCRIPTION": "", 
            "EXPOSEDTOINITIATOR": "false", 
            "HEALTHSTATUS": "1", 
            "ID": "15", 
            "IOCLASSID": "", 
            "IOPRIORITY": "1", 
            "SOURCELUNCAPACITY":"2097152", 
            "ISSCHEDULEDSNAP":"0", 
            "NAME": "gz0000_Snap_1707251229263", 
            "PARENTID": "12", 
            "PARENTNAME": "gz0000_Snap_1707251229076", 
            "PARENTTYPE": 27, 
            "ROLLBACKENDTIME": "-1", 
            "ROLLBACKRATE": "-1", 
            "ROLLBACKSPEED": "-1", 
            "ROLLBACKSTARTTIME": "-1", 
            "ROLLBACKTARGETOBJID": "4294967295", 
            "ROLLBACKTARGETOBJNAME": "--", 
            "RUNNINGSTATUS": "45", 
            "SOURCELUNID": "7", 
            "SOURCELUNNAME": "gz0000", 
            "SUBTYPE": "0", 
            "TIMESTAMP": "-1", 
            "TYPE": 27, 
            "USERCAPACITY": "2097152", 
            "WORKINGCONTROLLER": "0B", 
            "WWN": "6a1b2c3100f4d5e60008b5710000000f", 
            "replicationCapacity": "0", 
            "WORKLOADTYPEID": "13", 
            "WORKLOADTYPENAME": "Datab2ases" 
        } 
 ], 
 "error": { 
        "code": 0, 
        "description": "0" 
 } 
}`)
	})

	snapshots, err := client.LocalDevice.GetSnapshots(context.Background(), nil)
	if err != nil {
		t.Errorf("GetSnapshos return err: %s", err)
	}

	want := []Snapshot{
		{
			CASCADEDLEVEL:         "0",
			CASCADEDNUM:           "1",
			CONSUMEDCAPACITY:      "0",
			DESCRIPTION:           "",
			EXPOSEDTOINITIATOR:    "false",
			HEALTHSTATUS:          "1",
			ID:                    12,
			IOCLASSID:             "",
			IOPRIORITY:            "1",
			SOURCELUNCAPACITY:     "2097152",
			ISSCHEDULEDSNAP:       "0",
			NAME:                  "gz0000_Snap_1707251229076",
			PARENTID:              7,
			PARENTNAME:            "gz0000",
			PARENTTYPE:            TypeLUN,
			ROLLBACKENDTIME:       "-1",
			ROLLBACKRATE:          "-1",
			ROLLBACKSPEED:         "-1",
			ROLLBACKSTARTTIME:     "-1",
			ROLLBACKTARGETOBJID:   "4294967295",
			ROLLBACKTARGETOBJNAME: "--",
			RUNNINGSTATUS:         "43",
			SOURCELUNID:           "7",
			SOURCELUNNAME:         "gz0000",
			SUBTYPE:               "0",
			TIMESTAMP:             "1501013445",
			TYPE:                  TypeSnapshot,
			USERCAPACITY:          "2097152",
			WORKINGCONTROLLER:     "0B",
			WORKLOADTYPEID:        "14",
			WORKLOADTYPENAME:      "Databases",
			WWN:                   "6a1b2c3100f4d5e600089fd70000000c",
			ReplicationCapacity:   "0",
		}, {
			CASCADEDLEVEL:         "1",
			CASCADEDNUM:           "0",
			CONSUMEDCAPACITY:      "0",
			DESCRIPTION:           "",
			EXPOSEDTOINITIATOR:    "false",
			HEALTHSTATUS:          "1",
			ID:                    15,
			IOCLASSID:             "",
			IOPRIORITY:            "1",
			SOURCELUNCAPACITY:     "2097152",
			ISSCHEDULEDSNAP:       "0",
			NAME:                  "gz0000_Snap_1707251229263",
			PARENTID:              12,
			PARENTNAME:            "gz0000_Snap_1707251229076",
			PARENTTYPE:            TypeSnapshot,
			ROLLBACKENDTIME:       "-1",
			ROLLBACKRATE:          "-1",
			ROLLBACKSPEED:         "-1",
			ROLLBACKSTARTTIME:     "-1",
			ROLLBACKTARGETOBJID:   "4294967295",
			ROLLBACKTARGETOBJNAME: "--",
			RUNNINGSTATUS:         "45",
			SOURCELUNID:           "7",
			SOURCELUNNAME:         "gz0000",
			SUBTYPE:               "0",
			TIMESTAMP:             "-1",
			TYPE:                  TypeSnapshot,
			USERCAPACITY:          "2097152",
			WORKINGCONTROLLER:     "0B",
			WORKLOADTYPEID:        "13",
			WORKLOADTYPENAME:      "Datab2ases",
			WWN:                   "6a1b2c3100f4d5e60008b5710000000f",
			ReplicationCapacity:   "0",
		},
	}

	if !reflect.DeepEqual(snapshots, want) {
		t.Errorf("GetSnapshots return %+v, want %+v", snapshots, want)
	}
}
