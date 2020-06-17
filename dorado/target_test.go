package dorado

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDevice_GetTargetPort(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/iscsi_tgt_port", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w,
			`
{
  "data": [
    {
      "ETHPORTID": "1920220101",
      "ID": "0+iqn.2006-08.com.huawei:oceanstor:2100a400e255e226::20100:192.0.2.20,t,0x0101",
      "TPGT": "257",
      "TYPE": 249
    },
    {
      "ETHPORTID": "0000102",
      "ID": "0+iqn.2006-08.com.huawei:oceanstor:2100a400e255e226::20101:0.0.0.0,t,0x0102",
      "TPGT": "258",
      "TYPE": 249
    }
  ],
  "error": {
    "code": 0,
    "description": "0"
  }
}
`)
	})

	targets, err := client.LocalDevice.GetTargetPort(context.Background(), nil)
	if err != nil {
		t.Errorf("GetTargetPort return err: %s", err)
	}

	want := []TargetPort{
		{
			ETHPORTID: "1920220101",
			ID:        "0+iqn.2006-08.com.huawei:oceanstor:2100a400e255e226::20100:192.0.2.20,t,0x0101",
			TPGT:      "257",
			TYPE:      249,
		}, {
			ETHPORTID: "0000102",
			ID:        "0+iqn.2006-08.com.huawei:oceanstor:2100a400e255e226::20101:0.0.0.0,t,0x0102",
			TPGT:      "258",
			TYPE:      249,
		},
	}

	if !reflect.DeepEqual(targets, want) {
		t.Errorf("GetTargetPort return %+v, want %+v", targets, want)
	}
}
