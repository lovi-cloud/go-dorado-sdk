package dorado

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDevice_GetStoragePools(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/storagepool", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w,
			`
{ 
  "data": [ 
 { 
      "COMPRESSEDCAPACITY": "0", 
      "COMPRESSINVOLVEDCAPACITY": "0", 
      "COMPRESSIONRATE": "{\"numerator\":\"10\", \"denominator\":\"10\",\"logic\":\"=\"}", 
      "DATASPACE": "11996393472", 
      "DEDUPEDCAPACITY": "0", 
      "DEDUPINVOLVEDCAPACITY": "0", 
      "DEDUPLICATIONRATE": "{\"numerator\":\"10\", \"denominator\":\"10\",\"logic\":\"=\"}", 
      "DESCRIPTION": "", 
      "ENDINGUPTHRESHOLD": "90", 
      "HEALTHSTATUS": "1", 
      "ID": "3", 
      "LUNCONFIGEDCAPACITY": "62914560", 
      "NAME": "pool", 
      "PARENTID": "0", 
      "PARENTNAME": "domain", 
      "PARENTTYPE": 266, 
      "PROVISIONINGLIMIT": "-1", 
      "PROVISIONINGLIMITSWITCH": "false", 
      "REDUCTIONINVOLVEDCAPACITY": "0", 
      "REPLICATIONCAPACITY": "0", 
      "RUNNINGSTATUS": "27", 
      "SAVECAPACITYRATE": "{\"numerator\":\"1000\", \"denominator\":\"10\",\"logic\":\">=\"}", 
      "SPACEREDUCTIONRATE": "{\"numerator\":\"1000\", \"denominator\":\"1000\",\"logic\":\"=\"}", 
      "THINPROVISIONSAVEPERCENTAGE": "{\"numerator\":\"1000\", \"denominator\":\"10\",\"logic\":\"=\"}", 
      "TIER0CAPACITY": "11996393472", 
      "TIER0DISKTYPE": "3", 
      "TIER0RAIDLV": "5", 
      "TOTALLUNWRITECAPACITY": "0", 
      "TYPE": 216, 
      "USAGETYPE": "1", 
      "USERCONSUMEDCAPACITY": "0", 
      "USERCONSUMEDCAPACITYPERCENTAGE": "0", 
      "USERCONSUMEDCAPACITYTHRESHOLD": "80", 
      "USERCONSUMEDCAPACITYWITHOUTMETA": "0", 
      "USERFREECAPACITY": "11996393472", 
      "USERTOTALCAPACITY": "11996393472", 
      "USERWRITEALLOCCAPACITY": "0", 
      "autoDeleteSwitch": "0", 
      "poolProtectHighThreshold": "30", 
      "poolProtectLowThreshold": "20", 
      "protectSize": "0", 
      "totalSizeWithoutSnap": "41943040" 
    } 
  ], 
  "error": { 
 "code": 0, 
 "description": "0" 
  } 
}`)
	})

	storagepools, err := client.LocalDevice.GetStoragePools(context.Background(), nil)
	if err != nil {
		t.Errorf("GetStoragePools return err: %s", err)
	}

	want := []StoragePools{
		{
			COMPRESSEDCAPACITY:              "0",
			COMPRESSINVOLVEDCAPACITY:        "0",
			COMPRESSIONRATE:                 "{\"numerator\":\"10\", \"denominator\":\"10\",\"logic\":\"=\"}",
			DATASPACE:                       "11996393472",
			DEDUPEDCAPACITY:                 "0",
			DEDUPINVOLVEDCAPACITY:           "0",
			DEDUPLICATIONRATE:               "{\"numerator\":\"10\", \"denominator\":\"10\",\"logic\":\"=\"}",
			DESCRIPTION:                     "",
			ENDINGUPTHRESHOLD:               "90",
			HEALTHSTATUS:                    "1",
			ID:                              3,
			LUNCONFIGEDCAPACITY:             "62914560",
			NAME:                            "pool",
			PARENTID:                        "0",
			PARENTNAME:                      "domain",
			PARENTTYPE:                      266,
			PROVISIONINGLIMIT:               "-1",
			PROVISIONINGLIMITSWITCH:         "false",
			REDUCTIONINVOLVEDCAPACITY:       "0",
			REPLICATIONCAPACITY:             "0",
			RUNNINGSTATUS:                   "27",
			SAVECAPACITYRATE:                "{\"numerator\":\"1000\", \"denominator\":\"10\",\"logic\":\">=\"}",
			SPACEREDUCTIONRATE:              "{\"numerator\":\"1000\", \"denominator\":\"1000\",\"logic\":\"=\"}",
			THINPROVISIONSAVEPERCENTAGE:     "{\"numerator\":\"1000\", \"denominator\":\"10\",\"logic\":\"=\"}",
			TIER0CAPACITY:                   "11996393472",
			TIER0DISKTYPE:                   "3",
			TIER0RAIDLV:                     "5",
			TOTALLUNWRITECAPACITY:           "0",
			TYPE:                            216,
			USAGETYPE:                       "1",
			USERCONSUMEDCAPACITY:            "0",
			USERCONSUMEDCAPACITYPERCENTAGE:  "0",
			USERCONSUMEDCAPACITYTHRESHOLD:   "80",
			USERCONSUMEDCAPACITYWITHOUTMETA: "0",
			USERFREECAPACITY:                "11996393472",
			USERTOTALCAPACITY:               "11996393472",
			USERWRITEALLOCCAPACITY:          "0",
			AutoDeleteSwitch:                "0",
			PoolProtectHighThreshold:        "30",
			PoolProtectLowThreshold:         "20",
			ProtectSize:                     "0",
			TotalSizeWithoutSnap:            "41943040",
		},
	}

	if !reflect.DeepEqual(storagepools, want) {
		t.Errorf("GetStoragePools return %+v, want %+v", storagepools, want)
	}
}
