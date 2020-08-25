package dorado

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"sync"

	"github.com/pkg/errors"
)

// Client is client for go-dorado-sdk
type Client struct {
	LocalDevice  *Device
	RemoteDevice *Device

	PortGroupName string

	Logger *log.Logger
}

// Device is device of dorado
type Device struct {
	Controllers []*url.URL
	URL         *url.URL
	HTTPClient  *http.Client
	DeviceID    string
	Token       string
	Jar         *cookiejar.Jar
	Logger      *log.Logger

	Username string
	Password string
}

// Result is response of REST API
type Result struct {
	Data  interface{} `json:"data"`
	Error ErrorResp   `json:"error"`
}

// ErrorResp is error response of REST API
type ErrorResp struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion"`
}

var (
	userAgent = fmt.Sprintf("DoradoGoClient")
)

// httpMu lock to create *http.Request while to call setToken
var (
	httpMu sync.RWMutex
)

// NewClient create go-dorado-sdk client and set iBaseToken create by REST API.
func NewClient(localIPs, remoteIPs []string, username, password, portgroupName string, logger *log.Logger) (*Client, error) {
	client, err := NewClientDefaultToken(localIPs, remoteIPs, username, password, portgroupName, logger)
	if err != nil {
		return nil, err
	}

	err = client.SetToken()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewClientDefaultToken create go-dorado-sdk client.
// this function not call REST API.
func NewClientDefaultToken(localIPs, remoteIPs []string, username, password, portgroupName string, logger *log.Logger) (*Client, error) {
	// validate input value
	if len(username) == 0 {
		return nil, errors.New("username is required")
	}
	if len(password) == 0 {
		return nil, errors.New("password is required")
	}
	if len(localIPs) == 0 || len(remoteIPs) == 0 {
		return nil, errors.New("IPs is required")
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

	localDevice, err := newDevice(localIPs, username, password, httpClient, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Local Device: %w", err)
	}

	remoteDevice, err := newDevice(remoteIPs, username, password, httpClient, logger)
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

func newDevice(ips []string, username, password string, httpClient *http.Client, logger *log.Logger) (*Device, error) {
	var parsedURLs []*url.URL
	for _, ipStr := range ips {
		parsed, err := url.Parse(ipStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse input urls: %w", err)
		}

		parsedURLs = append(parsedURLs, parsed)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookiejar: %w", err)
	}

	d := &Device{
		Controllers: parsedURLs,
		HTTPClient:  httpClient,
		Username:    username,
		Password:    password,
		Jar:         jar,
		Logger:      logger,
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

func (d *Device) newRequest(ctx context.Context, method, spath string, body io.Reader) (*http.Request, error) {
	httpMu.RLock()
	u := *d.URL
	u.Path = path.Join(d.URL.Path, spath)
	httpMu.RUnlock()

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
