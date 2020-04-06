package dorado

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Host struct {
	DESCRIPTION     string `json:"DESCRIPTION"`
	HEALTHSTATUS    string `json:"HEALTHSTATUS"`
	ID              string `json:"ID"`
	INITIATORNUM    string `json:"INITIATORNUM"`
	IP              string `json:"IP"`
	ISADD2HOSTGROUP string `json:"ISADD2HOSTGROUP"`
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
	if len(hostname) > MAX_NAME_LENGTH {
		hash := md5.Sum([]byte(hostname))
		return hex.EncodeToString(hash[:])[:MAX_NAME_LENGTH]
	}

	return hostname
}

func (d *Device) GetHosts(ctx context.Context, query *SearchQuery) ([]Host, error) {
	spath := "/host"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	req = AddSearchQuery(req, query)
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	hosts := []Host{}
	if err = decodeBody(resp, &hosts); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return hosts, nil
}

func (d *Device) GetHost(ctx context.Context, hostId string) (*Host, error) {
	spath := fmt.Sprintf("/host/%s", hostId)

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	host := &Host{}
	if err = decodeBody(resp, host); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return host, nil
}

func (d *Device) CreateHost(ctx context.Context, hostname string) (*Host, error) {
	spath := "/host"
	param := struct {
		NAME            string `json:"NAME"`
		TYPE            string `json:"TYPE"`
		OPERATIONSYSTEM string `json:"OPERATIONSYSTEM"`
		DESCRIPTION     string `json:"DESCRIPTION"`
	}{
		NAME:            encodeHostName(hostname),
		TYPE:            "21", // NOTE(whywaita): I don't know nothing. this value from OpenStack cinder-driver
		OPERATIONSYSTEM: "0",
		DESCRIPTION:     hostname,
	}
	jb, err := json.Marshal(param)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreatePostValue)
	}
	req, err := d.newRequest(ctx, "POST", spath, bytes.NewBuffer(jb))
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	host := &Host{}
	if err = decodeBody(resp, host); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return host, nil
}

func (d *Device) DeleteHost(ctx context.Context, hostId string) error {
	spath := fmt.Sprintf("/host/%s", hostId)

	req, err := d.newRequest(ctx, "DELETE", spath, nil)
	if err != nil {
		return errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, ErrHTTPRequestDo)
	}

	var i interface{} // this endpoint return N/A
	if err = decodeBody(resp, i); err != nil {
		return errors.Wrap(err, ErrDecodeBody)
	}

	return nil
}
