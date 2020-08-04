package dorado

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	case ErrorCodeUnAuthorized:
		// please retry
		return ErrUnAuthorized
	}

	return fmt.Errorf("Dorado Internal Error: %s (code: %d) Suggestion: %s", e.Description, e.Code, e.Suggestion)
}

// requestWithRetry do HTTP Request and retry if return UnAuthorized token.
// set false in retried when call from outer.
func (d *Device) requestWithRetry(req *http.Request, out interface{}, retried bool) error {
	resp, err := d.request(req)
	if err != nil {
		return fmt.Errorf("failed to request: %w", err)
	}

	err = decodeBody(resp, out, d.Logger)
	if err == ErrUnAuthorized && retried == false {
		// retry after refresh token
		err = d.setToken()
		if err != nil {
			return fmt.Errorf("failed to setToken: %w", err)
		}
		req.Header.Set("iBaseToken", d.Token)
		d.HTTPClient.Jar = d.Jar
		return d.requestWithRetry(req, out, true)
	}

	if err != nil {
		return fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	return nil
}

func (d *Device) request(req *http.Request) (*http.Response, error) {
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	return resp, nil
}
