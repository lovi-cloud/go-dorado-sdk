package dorado

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDevice_GetLUNCopys(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/luncopy", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w,
			`
{
  "data": [
    {
      "BASELUN": "4294967295",
      "COPYPROGRESS": "-1",
      "COPYSPEED": "4",
      "COPYSTARTTIME": "1586176352",
      "COPYSTOPTIME": "1586176491",
      "DESCRIPTION": "",
      "HEALTHSTATUS": "1",
      "ID": "53",
      "LUNCOPYTYPE": "1",
      "NAME": "LUNCopy_148_151",
      "RUNNINGSTATUS": "40",
      "SOURCELUN": "INVALID;148;INVALID;INVALID;INVALID;",
      "SOURCELUNCAPACITY": "22020096",
      "SOURCELUNCAPACITYBYTE": "22548578304",
      "SOURCELUNNAME": "77bea474-b5fb1fed02ce0ed195abb1",
      "SOURCELUNWWN": "6a400e210055e22650d557a000000094",
      "SUBTYPE": "0",
      "TARGETLUN": "INVALID;151;INVALID;INVALID;INVALID;",
      "TYPE": 219
    }
  ],
  "error": {
    "code": 0,
    "description": "0"
  }
}`)
	})

	luncopys, err := client.LocalDevice.GetLUNCopys(context.Background(), nil)
	if err != nil {
		t.Errorf("GetLUNCopys return err: %s", err)
	}

	want := []LunCopy{
		{
			BASELUN:               "4294967295",
			COPYPROGRESS:          "-1",
			COPYSPEED:             "4",
			COPYSTARTTIME:         "1586176352",
			COPYSTOPTIME:          "1586176491",
			DESCRIPTION:           "",
			HEALTHSTATUS:          "1",
			ID:                    53,
			LUNCOPYTYPE:           "1",
			NAME:                  "LUNCopy_148_151",
			RUNNINGSTATUS:         "40",
			SOURCELUN:             "INVALID;148;INVALID;INVALID;INVALID;",
			SOURCELUNCAPACITY:     "22020096",
			SOURCELUNCAPACITYBYTE: "22548578304",
			SOURCELUNNAME:         "77bea474-b5fb1fed02ce0ed195abb1",
			SOURCELUNWWN:          "6a400e210055e22650d557a000000094",
			SUBTYPE:               "0",
			TARGETLUN:             "INVALID;151;INVALID;INVALID;INVALID;",
			TYPE:                  TypeLUNCopy,
		},
	}

	if !reflect.DeepEqual(luncopys, want) {
		t.Errorf("GetLUNCopys return %+v, want %+v", luncopys, want)
	}
}
