package dorado

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

func TestDevice_GetAssociatedEthernetPort(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/eth_port/associate", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, ` 
{ 
 "data": [ 
        { 
            "BONDID": "18446744073709551615", 
            "BONDNAME": "", 
            "ERRORPACKETS": "0", 
            "ETHDUPLEX": "-1", 
            "ETHNEGOTIATE": "3", 
            "HEALTHSTATUS": "1", 
            "ID": "131328", 
            "INIORTGT": "4", 
            "IPV4ADDR": "", 
            "IPV4GATEWAY": "", 
            "IPV4MASK": "", 
            "IPV6ADDR": "", 
            "IPV6GATEWAY": "", 
            "IPV6MASK": "", 
            "ISCSINAME": "", 
            "ISCSITCPPORT": "0", 
            "LOCATION": "CTE0.A.IOM1.P0", 
            "LOGICTYPE": "0", 
            "LOSTPACKETS": "0", 
            "MACADDRESS": "04:f9:38:95:88:f1", 
            "MTU": "1500", 
            "NAME": "P0", 
            "OVERFLOWEDPACKETS": "0", 
            "PARENTID": "0A.1", 
            "PARENTTYPE": 209, 
            "PORTSWITCH": "true", 
            "RUNNINGSTATUS": "11", 
            "SPEED": "-1", 
            "STARTTIME": "1474289326", 
            "TYPE": 213, 
            "crcErrors": "0", 
            "dswId": "4294967295", 
            "dswLinkRight": "4294967295", 
            "frameErrors": "0", 
            "frameLengthErrors": "0", 
            "lightStatus": "0", 
            "maxSpeed": "1000", 
            "selectType": "0", 
            "zoneId": "4294967295",
            "workModeType": "10", 
            "workModeList": "[10, 11, 12]" 
        } 
 ], 
 "error": { 
        "code": 0, 
        "description": "0" 
 } 
}`)
	})

	query := &SearchQuery{
		AssociateObjID:   "0A.1",
		AssociateObjType: strconv.Itoa(TypePortGroup),
	}

	ether, err := client.LocalDevice.GetAssociatedEthernetPort(context.Background(), query)
	if err != nil {
		t.Errorf("GetAssociatedEthernetPort return err: %s", err)
	}

	want := []EthernetPort{
		{
			BONDID:            "18446744073709551615",
			BONDNAME:          "",
			ERRORPACKETS:      "0",
			ETHDUPLEX:         "-1",
			ETHNEGOTIATE:      "3",
			HEALTHSTATUS:      "1",
			ID:                "131328",
			INIORTGT:          "4",
			IPV4ADDR:          "",
			IPV4GATEWAY:       "",
			IPV4MASK:          "",
			IPV6ADDR:          "",
			IPV6GATEWAY:       "",
			IPV6MASK:          "",
			ISCSINAME:         "",
			ISCSITCPPORT:      "0",
			LOCATION:          "CTE0.A.IOM1.P0",
			LOGICTYPE:         "0",
			LOSTPACKETS:       "0",
			MACADDRESS:        "04:f9:38:95:88:f1",
			MTU:               "1500",
			NAME:              "P0",
			OVERFLOWEDPACKETS: "0",
			PARENTID:          "0A.1",
			PARENTTYPE:        209,
			PORTSWITCH:        "true",
			RUNNINGSTATUS:     "11",
			SPEED:             "-1",
			STARTTIME:         "1474289326",
			TYPE:              TypeEthernetPort,
			CrcErrors:         "0",
			DswID:             "4294967295",
			DswLinkRight:      "4294967295",
			FrameErrors:       "0",
			FrameLengthErrors: "0",
			LightStatus:       "0",
			MaxSpeed:          "1000",
			SelectType:        "0",
			WorkModeList:      "[10, 11, 12]",
			WorkModeType:      "10",
			ZoneID:            "4294967295",
		},
	}

	if !reflect.DeepEqual(ether, want) {
		t.Errorf("GetAssociatedEthernetPort return %+v, want %+v", ether, want)
	}
}
