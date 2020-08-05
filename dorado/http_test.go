package dorado

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

var (
	testLogger = log.New(os.Stdout, "[go-dorado-sdk testing]", log.LstdFlags)
)

func TestDecodeBody_Interface(t *testing.T) {
	input := `
{
  "data": [],
  "error": {
    "code": 1077949002,
    "description": "The operation is not supported.",
    "suggestion": "Contact technical support engineers."
  }
}`
	resp := &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(input))}
	want := ErrorResp{
		Code:        1077949002,
		Description: "The operation is not supported.",
		Suggestion:  "Contact technical support engineers.",
	}

	var i interface{}
	err := decodeBody(resp, i, testLogger)
	if err == nil {
		t.Errorf("decodeBody return error is nil, want to return error response")
	}

	if err.Error() != want.Error().Error() {
		t.Errorf("decodeBody return %+v, want %+v", err, want)
	}
}

func TestDecodeBody_ErrorInvalidParameter(t *testing.T) {
	input := `{
  "data": {},
  "error": {
    "code": 1077674272,
    "description": "The entered HyperMetro parameters are invalid.",
    "suggestion": "Enter valid parameters."
  }
}`
	resp := &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(input))}
	want := ErrorResp{
		Code:        1077674272,
		Description: "The entered HyperMetro parameters are invalid.",
		Suggestion:  "Enter valid parameters.",
	}

	hyperMetroPair := &HyperMetroPair{}
	err := decodeBody(resp, &hyperMetroPair, testLogger)
	if err == nil {
		t.Errorf("decodeBody return error is nil, want to return error response")
	}

	if err.Error() != want.Error().Error() {
		t.Errorf("decodeBody return %+v, want %+v", err, want)
	}
}

func TestDecodeBody_Slice(t *testing.T) {
	input := `
{
  "data": [],
  "error": {
    "code": 1077949002,
    "description": "The operation is not supported.",
    "suggestion": "Contact technical support engineers."
  }
}`
	resp := &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(input))}
	want := ErrorResp{
		Code:        1077949002,
		Description: "The operation is not supported.",
		Suggestion:  "Contact technical support engineers.",
	}

	hyperMetroPairs := []HyperMetroPair{}
	err := decodeBody(resp, &hyperMetroPairs, testLogger)
	if err == nil {
		t.Errorf("decodeBody return error is nil, want to return error response")
	}

	if err.Error() != want.Error().Error() {
		t.Errorf("decodeBody return %+v, want %+v", err, want)
	}
}

func TestDecodeBody_SliceBadInput(t *testing.T) {
	input := `
{
  "data": {},
  "error": {
    "code": 1077949002,
    "description": "The operation is not supported.",
    "suggestion": "Contact technical support engineers."
  }
}`
	resp := &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(input))}

	hyperMetroPairs := []HyperMetroPair{} // input "data" is {} (not slice), but catch slice
	err := decodeBody(resp, &hyperMetroPairs, testLogger)
	if err == nil {
		t.Errorf("decodeBody return error is nil, want to return error response")
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if !errors.As(err, &unmarshalTypeError) {
		t.Errorf("decodeBody return %+v, want %+v", err, unmarshalTypeError)
	}
}
