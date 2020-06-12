package dorado

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestClient_GetHyperMetroPairs(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/HyperMetroPair", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w,
			`
{ 
 "data": [ 
        { 
            "CAPACITYBYTE": "10737418240", 
            "CGID": "", 
            "CGNAME": "", 
            "DOMAINID": "e4c2d1eaf02c0100", 
            "DOMAINNAME": "shuanyu", 
            "ENDTIME": "1", 
            "HCRESOURCETYPE": "1", 
            "HEALTHSTATUS": "1", 
            "HDRINGID": "--", 
            "ID": "e4c2d1eaf02c0001", 
            "ISINCG": "false", 
            "ISISOLATION":"false", 
            "ISISOLATIONTHRESHOLDTIME":"1000", 
            "ISPRIMARY": "true", 
            "LINKSTATUS": "1", 
            "LOCALDATASTATE": "1", 
            "LOCALHOSTACCESSSTATE": "3", 
            "LOCALOBJID": "216", 
            "LOCALOBJNAME": "lsm94_LUN2170000", 
            "RECOVERYPOLICY": "1", 
            "REMOTEDATASTATE": "2", 
            "REMOTEHOSTACCESSSTATE": "1", 
            "REMOTEOBJID": "514", 
            "REMOTEOBJNAME": "LSM_LUN5150000", 
            "RESOURCEWWN": "6e4c2d1100eaf02c04fa5287000000d8", 
            "RUNNINGSTATUS": "41", 
            "SPEED": "2", 
            "STARTTIME": "-1", 
            "SYNCDIRECTION": "1", 
            "SYNCLEFTTIME": "-1", 
            "SYNCPROGRESS": "1", 
            "TYPE": 15361, 
            "WRITESECONDARYTIMEOUT": "30" 
        } 
 ], 
 "error": { 
        "code": 0, 
        "description": "0" 
 } 
}`)
	})

	hmps, err := client.GetHyperMetroPairs(context.Background(), nil)
	if err != nil {
		t.Errorf("GetHyperMetroPairs return err: %s", err)
	}

	want := []HyperMetroPair{
		{
			CAPACITYBYTE:             "10737418240",
			CGID:                     "",
			CGNAME:                   "",
			DOMAINID:                 "e4c2d1eaf02c0100",
			DOMAINNAME:               "shuanyu",
			ENDTIME:                  "1",
			HCRESOURCETYPE:           "1",
			HDRINGID:                 "--",
			HEALTHSTATUS:             "1",
			ID:                       "e4c2d1eaf02c0001",
			ISINCG:                   "false",
			ISISOLATION:              "false",
			ISISOLATIONTHRESHOLDTIME: "1000",
			ISPRIMARY:                "true",
			LINKSTATUS:               "1",
			LOCALDATASTATE:           "1",
			LOCALHOSTACCESSSTATE:     "3",
			LOCALOBJID:               216,
			LOCALOBJNAME:             "lsm94_LUN2170000",
			RECOVERYPOLICY:           "1",
			REMOTEDATASTATE:          "2",
			REMOTEHOSTACCESSSTATE:    "1",
			REMOTEOBJID:              514,
			REMOTEOBJNAME:            "LSM_LUN5150000",
			RESOURCEWWN:              "6e4c2d1100eaf02c04fa5287000000d8",
			RUNNINGSTATUS:            "41",
			SPEED:                    "2",
			STARTTIME:                "-1",
			SYNCDIRECTION:            "1",
			SYNCLEFTTIME:             "-1",
			SYNCPROGRESS:             "1",
			TYPE:                     TypeHyperMetroPair,
			WRITESECONDARYTIMEOUT:    "30",
		},
	}

	if !reflect.DeepEqual(hmps, want) {
		t.Errorf("GetHyperMetroPairs return %+v, want %+v", hmps, want)
	}
}
