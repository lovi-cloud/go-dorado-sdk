package dorado

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDevice_GetSystem(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/system/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w,
			`
{ 
 "data": { 
        "CACHEWRITEQUOTA": "333", 
        "CONFIGMODEL": "1", 
        "DESCRIPTION": "", 
        "DOMAINNAME": "", 
        "FREEDISKSCAPACITY": "40920093408", 
        "HEALTHSTATUS": "1", 
        "HOTSPAREDISKSCAPACITY": "0", 
        "ID": "210235843910E6000009", 
        "LOCATION": "", 
        "MEMBERDISKSCAPACITY": "28500963638", 
        "NAME": "Huawei.Storage", 
        "PRODUCTMODE": "68", 
        "PRODUCTVERSION": "V300R005C00", 
        "RUNNINGSTATUS": "1", 
        "SECTORSIZE": "512", 
        "STORAGEPOOLCAPACITY": "6694109184", 
        "STORAGEPOOLFREECAPACITY": "6344409088", 
        "STORAGEPOOLHOSTSPARECAPACITY": "2607546367", 
        "STORAGEPOOLRAWCAPACITY": "17127200853", 
        "STORAGEPOOLUSEDCAPACITY": "349700096", 
        "THICKLUNSALLOCATECAPACITY": "320339968", 
        "THICKLUNSUSEDCAPACITY": "-1", 
        "THINLUNSALLOCATECAPACITY": "0", 
        "THINLUNSMAXCAPACITY": "0", 
        "THINLUNSUSEDCAPACITY": "-1", 
        "TOTALCAPACITY": "69421057046", 
        "TYPE": 201, 
        "UNAVAILABLEDISKSCAPACITY": "1559321616", 
        "USEDCAPACITY": "349700096", 
        "VASA_ALTERNATE_NAME": "Huawei.Storage", 
        "VASA_SUPPORT_BLOCK": "", 
        "VASA_SUPPORT_FILESYSTEM": "NFS", 
        "VASA_SUPPORT_PROFILE": "FileSystemProfile", 
        "WRITETHROUGHSW": "true", 
        "WRITETHROUGHTIME": "72", 
        "mappedLunsCountCapacity": "0", 
        "patchVersion": "", 
        "unMappedLunsCountCapacity": "306184192", 
        "userFreeCapacity": "52991726636", 
        "wwn": "210030d17eb4f761" 
 }, 
 "error": { 
        "code": 0, 
        "description": "0" 
 } 
}`)
	})

	system, err := client.LocalDevice.GetSystem(context.Background())
	if err != nil {
		t.Errorf("GetSystem return err: %v", err)
	}

	want := &System{
		CACHEWRITEQUOTA:              "333",
		CONFIGMODEL:                  "1",
		DESCRIPTION:                  "",
		DOMAINNAME:                   "",
		FREEDISKSCAPACITY:            "40920093408",
		HEALTHSTATUS:                 "1",
		HOTSPAREDISKSCAPACITY:        "0",
		ID:                           "210235843910E6000009",
		LOCATION:                     "",
		MEMBERDISKSCAPACITY:          "28500963638",
		NAME:                         "Huawei.Storage",
		PRODUCTMODE:                  "68",
		PRODUCTVERSION:               "V300R005C00",
		RUNNINGSTATUS:                "1",
		SECTORSIZE:                   "512",
		STORAGEPOOLCAPACITY:          "6694109184",
		STORAGEPOOLFREECAPACITY:      "6344409088",
		STORAGEPOOLHOSTSPARECAPACITY: "2607546367",
		STORAGEPOOLRAWCAPACITY:       "17127200853",
		STORAGEPOOLUSEDCAPACITY:      "349700096",
		THICKLUNSALLOCATECAPACITY:    "320339968",
		THICKLUNSUSEDCAPACITY:        "-1",
		THINLUNSALLOCATECAPACITY:     "0",
		THINLUNSMAXCAPACITY:          "0",
		THINLUNSUSEDCAPACITY:         "-1",
		TOTALCAPACITY:                "69421057046",
		TYPE:                         201,
		UNAVAILABLEDISKSCAPACITY:     "1559321616",
		USEDCAPACITY:                 "349700096",
		VASAALTERNATENAME:            "Huawei.Storage",
		VASASUPPORTBLOCK:             "",
		VASASUPPORTFILESYSTEM:        "NFS",
		VASASUPPORTPROFILE:           "FileSystemProfile",
		WRITETHROUGHSW:               "true",
		WRITETHROUGHTIME:             "72",
		MappedLunsCountCapacity:      "0",
		PatchVersion:                 "",
		UnMappedLunsCountCapacity:    "306184192",
		UserFreeCapacity:             "52991726636",
		Wwn:                          "210030d17eb4f761",
	}

	if !reflect.DeepEqual(system, want) {
		t.Errorf("GetSystem return %+v, want %+v", system, want)
	}
}
