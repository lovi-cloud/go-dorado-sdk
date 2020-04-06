package dorado

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// storage - host mapping must have a host group.
// host group has only one host under our usage.

type HostGroup struct {
	DESCRIPTION       string `json:"DESCRIPTION"`
	ID                string `json:"ID"`
	ISADD2MAPPINGVIEW string `json:"ISADD2MAPPINGVIEW"`
	NAME              string `json:"NAME"`
	TYPE              int    `json:"TYPE"`
}

const (
	ErrHostGroupNotFound = "HostGroup is not found"
)

func (d *Device) GetHostGroups(ctx context.Context, query *SearchQuery) ([]HostGroup, error) {
	spath := "/hostgroup"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	req = AddSearchQuery(req, query)
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	hostGroups := []HostGroup{}
	if err = decodeBody(resp, &hostGroups); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	if len(hostGroups) == 0 {
		return nil, errors.New(ErrHostGroupNotFound)
	}

	return hostGroups, nil
}

func (d *Device) GetHostGroup(ctx context.Context, hostgroupId string) (*HostGroup, error) {
	spath := fmt.Sprintf("/hostgroup/%s", hostgroupId)

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	hostGroup := &HostGroup{}
	if err = decodeBody(resp, hostGroup); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return hostGroup, nil
}

func (d *Device) CreateHostGroup(ctx context.Context, hostname string) (*HostGroup, error) {
	spath := "/hostgroup"
	param := struct {
		NAME        string `json:"NAME"`
		DESCRIPTION string `json:"DESCRIPTION"`
	}{
		NAME:        encodeHostName(hostname),
		DESCRIPTION: hostname,
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

	hg := HostGroup{}
	if err = decodeBody(resp, &hg); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return &hg, nil
}

func (d *Device) DeleteHostGroup(ctx context.Context, hostGroupId string) error {
	spath := fmt.Sprintf("/hostgroup/%s", hostGroupId)

	req, err := d.newRequest(ctx, "DELETE", spath, nil)
	if err != nil {
		return errors.Wrap(err, ErrCreatePostValue)
	}
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, ErrHTTPRequestDo)
	}

	var i interface{} // this endpoint return N/A
	if err = decodeBody(resp, &i); err != nil {
		return errors.Wrap(err, ErrDecodeBody)
	}

	return nil
}

type AssociateParam struct {
	ID               string `json:"ID"`
	ASSOCIATEOBJID   string `json:"ASSOCIATEOBJID"`
	ASSOCIATEOBJTYPE int    `json:"ASSOCIATEOBJTYPE"`
}

func (d *Device) AssociateHost(ctx context.Context, hostgroupId, hostId string) error {
	spath := "/hostgroup/associate"
	param := AssociateParam{
		ID:               hostgroupId,
		ASSOCIATEOBJID:   hostId,
		ASSOCIATEOBJTYPE: 21, // 21 is Host
	}
	jb, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err, ErrCreatePostValue)
	}

	req, err := d.newRequest(ctx, "POST", spath, bytes.NewBuffer(jb))
	if err != nil {
		return errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, ErrHTTPRequestDo)
	}

	var i interface{} // this endpoint return N/A
	if err = decodeBody(resp, &i); err != nil {
		return errors.Wrap(err, ErrDecodeBody)
	}

	return nil
}

func (d *Device) DisAssociateHost(ctx context.Context, hostgroupId, hostId string) error {
	spath := "/hostgroup/associate"
	param := AssociateParam{
		ID:               hostgroupId,
		ASSOCIATEOBJID:   hostId,
		ASSOCIATEOBJTYPE: 21, // 21 is Host
	}
	jb, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err, ErrCreatePostValue)
	}

	req, err := d.newRequest(ctx, "DELETE", spath, bytes.NewBuffer(jb))
	if err != nil {
		return errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, ErrHTTPRequestDo)
	}

	var i interface{} // this endpoint return N/A
	if err = decodeBody(resp, &i); err != nil {
		return errors.Wrap(err, ErrDecodeBody)
	}

	return nil
}

func (d *Device) CreateHostGroupWithHost(ctx context.Context, hostname string) (*HostGroup, error) {
	host, err := d.CreateHost(ctx, hostname)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Host")
	}

	hostgroup, err := d.CreateHostGroup(ctx, hostname)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create hostgroup")
	}

	err = d.AssociateHost(ctx, hostgroup.ID, host.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to associate to hostgroup")
	}

	return hostgroup, nil
}

func (d *Device) DeleteHostGroupWithHost(ctx context.Context, hostgroupId string) error {
	hostgroup, err := d.GetHostGroup(ctx, hostgroupId)
	if err != nil {
		return errors.Wrap(err, "failed to search hostgroup by ID")
	}
	hosts, err := d.GetHosts(ctx, CreateSearchName(hostgroup.NAME))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to search host that name is %s", hostgroup.NAME))
	}
	if len(hosts) != 1 {
		return errors.New("search result of host is not one")
	}
	host := hosts[0]

	err = d.DisAssociateHost(ctx, hostgroup.ID, host.ID)
	if err != nil {
		return errors.Wrap(err, "failed to deassociate hostgroup")
	}
	err = d.DeleteHost(ctx, host.ID)
	if err != nil {
		return errors.Wrap(err, "failed to delete host")
	}
	err = d.DeleteHostGroup(ctx, hostgroup.ID)
	if err != nil {
		return errors.Wrap(err, "failed to delete hostgroup")
	}

	return nil
}
