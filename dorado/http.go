package dorado

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ErrCreateRequest   = "failed to create request"
	ErrHTTPRequestDo   = "failed to HTTP request"
	ErrDecodeBody      = "failed to decodeBody"
	ErrCreatePostValue = "failed to create post value"
)

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	r := &Result{
		Data: out,
	}
	err := decoder.Decode(r)
	if err != nil {
		return fmt.Errorf("failed to create json decoder: %w", err)
	}
	if r.Error.Error() != nil {
		return r.Error.Error()
	}

	out = r.Data
	return nil
}

func (e *ErrorResp) Error() error {
	if e == nil {
		return nil
	}

	if e.Code == 0 {
		// no error
		return nil
	}
	return fmt.Errorf("Dorado Internal Error: %s (code: %d) Suggestion: %s", e.Description, e.Code, e.Suggestion)
}
