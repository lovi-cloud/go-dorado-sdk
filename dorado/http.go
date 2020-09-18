package dorado

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func decodeBody(resp *http.Response, out interface{}, logger *log.Logger) error {
	defer resp.Body.Close()
	jb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	r := &Result{
		Data: out,
	}
	if err := json.Unmarshal(jb, r); err != nil {
		logger.Printf("Dorado response: %v", string(jb))
		return fmt.Errorf("failed to unmarshal response JSON: %w", err)
	}

	if r.Error.Error() != nil {
		logger.Printf("Dorado return error: %v", r.Error.Error())
		return r.Error.Error()
	}

	out = r.Data
	return nil
}

func (e ErrorResp) Error() error {
	switch e.Code {
	case 0:
		// no error
		return nil
	case ErrorCodeUnAuthorized, ErrorCodeUserIsOffline:
		// please retry
		return ErrUnAuthorized
	}

	return fmt.Errorf("Dorado Internal Error: %s (code: %d) Suggestion: %s", e.Description, e.Code, e.Suggestion)
}

// requestWithRetry do HTTP Request and retry if return UnAuthorized token.
// set false in retried when call from outer.
func (d *Device) requestWithRetry(req *http.Request, out interface{}, retryCount int) error {
	resp, err := d.request(req)
	if err != nil {
		return fmt.Errorf("failed to request: %w", err)
	}

	err = decodeBody(resp, out, d.Logger)
	if err == ErrUnAuthorized && retryCount > 0 {
		// retry after refresh token
		// need update iBaseToken and ismsession in Cookie
		err = d.setToken()
		if err != nil {
			return fmt.Errorf("failed to setToken: %w", err)
		}

		spath := strings.TrimPrefix(req.URL.Path, d.URL.Path)
		var jb []byte
		if req.GetBody != nil {
			b, err := req.GetBody()
			if err != nil {
				return fmt.Errorf("failed to GetBody: %w", err)
			}

			jb, err = ioutil.ReadAll(b) // NOTE(whywaita): need to fix many memory allocation if occurred problem
			if err != nil {
				return fmt.Errorf("failed to ReadAll: %w", err)
			}
		}

		newReq, err := d.newRequest(req.Context(), req.Method, spath, bytes.NewBuffer(jb))
		if err != nil {
			return fmt.Errorf("failed to create new http request: %w", err)
		}

		return d.requestWithRetry(newReq, out, retryCount-1)
	}

	if err != nil {
		return fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	return nil
}

func (d *Device) request(req *http.Request) (*http.Response, error) {
	d.Logger.Printf("Do Request %+v\n", req)

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	return resp, nil
}
