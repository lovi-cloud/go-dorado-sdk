package dorado

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// Session is response of /sessions
type Session struct {
	IBaseToken string `json:"iBaseToken"`
	DeviceID   string `json:"deviceid"`
}

func (d *Device) getToken() (string, string, *http.CookieJar, error) {
	spath := "/sessions"

	param := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Scope    int    `json:"scope"`
	}{
		Username: d.Username,
		Password: d.Password,
		Scope:    0,
	}
	jb, err := json.Marshal(param)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to json.Marshal: %w", err)
	}
	urlStr := d.URL.String()
	d.HTTPClient.Jar = d.Jar
	resp, err := d.HTTPClient.Post(urlStr+spath, "application/json", bytes.NewBuffer(jb))
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to get token request: %w", err)
	}
	defer resp.Body.Close()

	body := &Session{}
	err = decodeBody(resp, body, d.Logger)
	if err != nil {
		return "", "", nil, fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	return body.IBaseToken, body.DeviceID, &d.HTTPClient.Jar, nil
}

// SetToken set iBaseToken from REST API.
func (c *Client) SetToken() error {
	var err error

	err = c.LocalDevice.setToken()
	if err != nil {
		return fmt.Errorf("failed to set token in local device: %w", err)
	}
	err = c.RemoteDevice.setToken()
	if err != nil {
		return fmt.Errorf("failed to set token in remote device: %w", err)
	}

	return nil
}

func (d *Device) setToken() error {
	for _, url := range d.Controllers {
		err := d.setBaseURL(url.String(), DefaultDeviceID)
		if err != nil {
			return fmt.Errorf("failed to set BaseURL: %w", err)
		}

		token, deviceID, jar, err := d.getToken()
		if err != nil {
			d.Logger.Printf("cannot get token, continue next controller (URL: %s): %s", url.String(), err)
			continue
		}
		d.DeviceID = deviceID
		d.Token = token
		d.HTTPClient.Jar = *jar

		err = d.setBaseURL(url.String(), deviceID)
		if err != nil {
			return fmt.Errorf("failed to set BaseURL: %w", err)
		}

		d.Logger.Printf("successlay setToken! (URL: %s)", url.String())
		return nil
	}

	return errors.New("cannot setToken in all controllers")
}
