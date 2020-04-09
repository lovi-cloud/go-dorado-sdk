package dorado

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

type Client struct {
	LocalDevice  *Device
	RemoteDevice *Device

	Logger *log.Logger
}

type Device struct {
	URL        *url.URL
	HTTPClient *http.Client
	DeviceId   string
	Token      string
	Jar        *cookiejar.Jar

	Username string
	Password string
}

type Result struct {
	Data  interface{} `json:"data"`
	Error ErrorResp   `json:"error"`
}

type ErrorResp struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion"`
}

var (
	userAgent = fmt.Sprintf("DoradoGoClient")
)

const (
	DefaultDeviceId = "xx"
)

func NewClient(ctx context.Context, localIp, remoteIp, username, password string, logger *log.Logger) (*Client, error) {
	// validate input value
	if len(username) == 0 {
		return nil, errors.New("username is required")
	}
	if len(password) == 0 {
		return nil, errors.New("password is required")
	}

	l := log.New(ioutil.Discard, "", log.LstdFlags)
	if logger == nil {
		logger = l
	}

	tlsConfig := tls.Config{
		InsecureSkipVerify: true,
	}
	transport := *http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tlsConfig
	httpClient := &http.Client{Transport: &transport}

	localDevice, err := newDevice(ctx, localIp, username, password, httpClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Local Device")
	}

	remoteDevice, err := newDevice(ctx, remoteIp, username, password, httpClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Remote Device")
	}

	c := &Client{
		LocalDevice:  localDevice,
		RemoteDevice: remoteDevice,
		Logger:       logger,
	}

	return c, nil
}

func newDevice(ctx context.Context, ipStr, username, password string, httpClient *http.Client) (*Device, error) {
	urlStr := fmt.Sprintf("https://%s/deviceManager/rest/%s", ipStr, DefaultDeviceId)
	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse url")
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cookiejar")
	}

	d := &Device{
		URL:        parsedURL,
		HTTPClient: httpClient,
		Username:   username,
		Password:   password,
		Jar:        jar,
	}

	token, deviceId, err := d.getToken()
	if err != nil {
		return nil, errors.Wrap(err, "failed to getToken")
	}
	d.DeviceId = deviceId
	d.Token = token
	deviceURL := fmt.Sprintf("https://%s/deviceManager/rest/%s", ipStr, deviceId)
	parsedDeviceURL, err := url.ParseRequestURI(deviceURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse url")
	}
	d.URL = parsedDeviceURL

	return d, nil
}

func (d *Device) newRequest(ctx context.Context, method, spath string, body io.Reader) (*http.Request, error) {
	u := *d.URL
	u.Path = path.Join(d.URL.Path, spath)

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new HTTP Request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	req.Header.Set("iBaseToken", d.Token)
	d.HTTPClient.Jar = d.Jar

	return req, nil
}

type Session struct {
	IBaseToken string `json:"iBaseToken"`
	DeviceId   string `json:"deviceid"`
}

func (d *Device) getToken() (string, string, error) {
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
		return "", "", errors.Wrap(err, "failed to json.Marshal")
	}
	urlStr := d.URL.String()
	d.HTTPClient.Jar = d.Jar
	resp, err := d.HTTPClient.Post(urlStr+spath, "application/json", bytes.NewBuffer(jb))
	if err != nil {
		return "", "", errors.Wrap(err, "failed to get token request")
	}
	defer resp.Body.Close()

	body := &Session{}
	err = decodeBody(resp, body)
	if err != nil {
		return "", "", errors.Wrap(err, ErrDecodeBody)
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to get token: %w", body)
	}

	return body.IBaseToken, body.DeviceId, nil
}
