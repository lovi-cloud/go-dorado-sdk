package dorado

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestDecodeBody(t *testing.T) {
	errResp := `
{
  "data": [],
  "error": {
    "code": 1077949002,
    "description": "The operation is not supported.",
    "suggestion": "Contact technical support engineers."
  }
}`
	resp := http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(errResp))}

	var i interface{}
	err := decodeBody(&resp, i)
	if err == nil {
		t.Errorf("decodeBody return error is nil, want to return error response")
	}

	want := ErrorResp{
		Code:        1077949002,
		Description: "The operation is not supported.",
		Suggestion:  "Contact technical support engineers.",
	}

	if err.Error() != want.Error().Error() {
		t.Errorf("decodeBody return %+v, want %+v", err, want)
	}
}
