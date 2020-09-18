package dorado

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
)

// Host is object of hypervisor (a.k.a. compute node) in dorado.
type Host struct {
	DESCRIPTION     string `json:"DESCRIPTION"`
	HEALTHSTATUS    string `json:"HEALTHSTATUS"`
	ID              int    `json:"ID,string"`
	INITIATORNUM    string `json:"INITIATORNUM"`
	IP              string `json:"IP"`
	ISADD2HOSTGROUP bool   `json:"ISADD2HOSTGROUP,string"`
	LOCATION        string `json:"LOCATION"`
	MODEL           string `json:"MODEL"`
	NAME            string `json:"NAME"`
	NETWORKNAME     string `json:"NETWORKNAME"`
	OPERATIONSYSTEM string `json:"OPERATIONSYSTEM"`
	PARENTID        string `json:"PARENTID"`
	PARENTNAME      string `json:"PARENTNAME"`
	PARENTTYPE      int    `json:"PARENTTYPE"`
	RUNNINGSTATUS   string `json:"RUNNINGSTATUS"`
	TYPE            int    `json:"TYPE"`
}

func encodeHostName(hostname string) string {
	// this function binding by huawei_utils.encode_host_name(id) in OpenStack cinder-driver.
	if len(hostname) > MaxNameLength {
		hash := md5.Sum([]byte(hostname))
		return hex.EncodeToString(hash[:])[:MaxNameLength]
	}

	return hostname
}

// GetHosts get host objects query by SearchQuery.
func (d *Device) GetHosts(ctx context.Context, query *SearchQuery) ([]Host, error) {
	spath := "/host"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req = AddSearchQuery(req, query)

	var hosts []Host
	if err = d.requestWithRetry(req, &hosts, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	if len(hosts) == 0 {
		return nil, ErrHostNotFound
	}

	return hosts, nil
}

// GetHost get host object by host ID.
func (d *Device) GetHost(ctx context.Context, hostID int) (*Host, error) {
	spath := fmt.Sprintf("/host/%d", hostID)

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	host := &Host{}
	if err = d.requestWithRetry(req, host, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return host, nil
}

// CreateHost create host object.
func (d *Device) CreateHost(ctx context.Context, hostname string) (*Host, error) {
	spath := "/host"
	param := struct {
		NAME            string `json:"NAME"`
		TYPE            string `json:"TYPE"`
		OPERATIONSYSTEM string `json:"OPERATIONSYSTEM"`
		DESCRIPTION     string `json:"DESCRIPTION"`
	}{
		NAME:            encodeHostName(hostname),
		TYPE:            strconv.Itoa(TypeHost),
		OPERATIONSYSTEM: "0",
		DESCRIPTION:     hostname,
	}
	jb, err := json.Marshal(param)
	if err != nil {
		return nil, fmt.Errorf(ErrCreatePostValue+": %w", err)
	}
	req, err := d.newRequest(ctx, "POST", spath, bytes.NewBuffer(jb))
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	host := &Host{}
	if err = d.requestWithRetry(req, host, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return host, nil
}

// DeleteHost delete host object.
func (d *Device) DeleteHost(ctx context.Context, hostID int) error {
	spath := fmt.Sprintf("/host/%d", hostID)

	req, err := d.newRequest(ctx, "DELETE", spath, nil)
	if err != nil {
		return fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	var i interface{} // this endpoint return N/A
	if err = d.requestWithRetry(req, i, DefaultHTTPRetryCount); err != nil {
		return fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return nil
}
