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
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

type Client struct {
	LocalDevice  *Device
	RemoteDevice *Device

	PortGroupName string

	Logger *log.Logger
}

type Device struct {
	IPAddress  *net.TCPAddr // TODO: implement for dual controller
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

func NewClient(localIp, remoteIp, username, password, portgroupName string, logger *log.Logger) (*Client, error) {
	client, err := NewClientDefaultToken(localIp, remoteIp, username, password, portgroupName, logger)
	if err != nil {
		return nil, err
	}

	err = client.SetToken()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewClientDefaultToken(localIp, remoteIp, username, password, portgroupName string, logger *log.Logger) (*Client, error) {
	// validate input value
	if len(username) == 0 {
		return nil, errors.New("username is required")
	}
	if len(password) == 0 {
		return nil, errors.New("password is required")
	}

	if logger == nil {
		l := log.New(ioutil.Discard, "", log.LstdFlags)
		logger = l
	}

	tlsConfig := tls.Config{
		InsecureSkipVerify: true,
	}
	transport := *http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tlsConfig
	httpClient := &http.Client{Transport: &transport}

	localDevice, err := newDevice(localIp, username, password, httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create Local Device: %w", err)
	}

	remoteDevice, err := newDevice(remoteIp, username, password, httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create Remote Device: %w", err)
	}

	c := &Client{
		LocalDevice:   localDevice,
		RemoteDevice:  remoteDevice,
		PortGroupName: portgroupName,
		Logger:        logger,
	}

	return c, nil
}

func newDevice(ipStr, username, password string, httpClient *http.Client) (*Device, error) {
	parsed, err := net.ResolveTCPAddr("tcp", ipStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input ipstr: %w", err)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookiejar: %w", err)
	}

	d := &Device{
		IPAddress:  parsed,
		HTTPClient: httpClient,
		Username:   username,
		Password:   password,
		Jar:        jar,
	}

	err = d.setBaseURL("https://"+ipStr, DefaultDeviceId)
	if err != nil {
		return nil, fmt.Errorf("failed to set BaseURL: %w", err)
	}

	return d, nil
}

func (d *Device) setBaseURL(baseHost, token string) error {
	urlStr := fmt.Sprintf("%s/deviceManager/rest/%s", baseHost, token)
	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}

	d.URL = parsedURL
	return nil
}

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
	token, deviceId, err := d.getToken()
	if err != nil {
		return fmt.Errorf("failed to getToken: %w", err)
	}
	d.DeviceId = deviceId
	d.Token = token

	err = d.setBaseURL("https://"+d.IPAddress.String(), deviceId)
	if err != nil {
		return fmt.Errorf("failed to set BaseURL: %w", err)
	}

	return nil
}

func (d *Device) newRequest(ctx context.Context, method, spath string, body io.Reader) (*http.Request, error) {
	u := *d.URL
	u.Path = path.Join(d.URL.Path, spath)

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new HTTP Request: %w", err)
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
		return "", "", fmt.Errorf("failed to json.Marshal: %w", err)
	}
	urlStr := d.URL.String()
	d.HTTPClient.Jar = d.Jar
	resp, err := d.HTTPClient.Post(urlStr+spath, "application/json", bytes.NewBuffer(jb))
	if err != nil {
		return "", "", fmt.Errorf("failed to get token request: %w", err)
	}
	defer resp.Body.Close()

	body := &Session{}
	err = decodeBody(resp, body)
	if err != nil {
		return "", "", fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	return body.IBaseToken, body.DeviceId, nil
}
