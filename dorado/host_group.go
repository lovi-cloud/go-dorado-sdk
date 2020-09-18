package dorado

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// HostGroup is object of multiple host.
// storage - host mapping must have a host group.
// host group has only one host under our usage.
type HostGroup struct {
	DESCRIPTION       string `json:"DESCRIPTION"`
	ID                int    `json:"ID,string"`
	ISADD2MAPPINGVIEW bool   `json:"ISADD2MAPPINGVIEW,string"`
	NAME              string `json:"NAME"`
	TYPE              int    `json:"TYPE"`
}

// GetHostGroups get hostgroup objects by query.
func (d *Device) GetHostGroups(ctx context.Context, query *SearchQuery) ([]HostGroup, error) {
	spath := "/hostgroup"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req = AddSearchQuery(req, query)

	var hostGroups []HostGroup
	if err = d.requestWithRetry(req, &hostGroups, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	if len(hostGroups) == 0 {
		return nil, ErrHostGroupNotFound
	}

	return hostGroups, nil
}

// GetHostGroup get hostgroup object by id.
func (d *Device) GetHostGroup(ctx context.Context, hostgroupID int) (*HostGroup, error) {
	spath := fmt.Sprintf("/hostgroup/%d", hostgroupID)

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	hostGroup := &HostGroup{}
	if err = d.requestWithRetry(req, hostGroup, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return hostGroup, nil
}

// CreateHostGroup create hostgroup object.
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
		return nil, fmt.Errorf(ErrCreatePostValue+": %w", err)
	}
	req, err := d.newRequest(ctx, "POST", spath, bytes.NewBuffer(jb))
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	hostGroup := &HostGroup{}
	if err = d.requestWithRetry(req, hostGroup, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return hostGroup, nil
}

// DeleteHostGroup delete hostgroup object.
func (d *Device) DeleteHostGroup(ctx context.Context, hostGroupID int) error {
	spath := fmt.Sprintf("/hostgroup/%d", hostGroupID)

	req, err := d.newRequest(ctx, "DELETE", spath, nil)
	if err != nil {
		return fmt.Errorf(ErrCreatePostValue+": %w", err)
	}

	var i interface{} // this endpoint return N/A
	if err = d.requestWithRetry(req, i, DefaultHTTPRetryCount); err != nil {
		return fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return nil
}

// AssociateHost associate hostID to hostgroupID.
func (d *Device) AssociateHost(ctx context.Context, hostgroupID, hostID int) error {
	spath := "/hostgroup/associate"
	param := AssociateParam{
		ID:               strconv.Itoa(hostgroupID),
		ASSOCIATEOBJID:   strconv.Itoa(hostID),
		ASSOCIATEOBJTYPE: TypeHost,
	}

	jb, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf(ErrCreatePostValue+": %w", err)
	}

	req, err := d.newRequest(ctx, "POST", spath, bytes.NewBuffer(jb))
	if err != nil {
		return fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	var i interface{} // this endpoint return N/A
	if err = d.requestWithRetry(req, i, DefaultHTTPRetryCount); err != nil {
		return fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return nil
}

// DisAssociateHost delete associate hostID from hostgroupID.
func (d *Device) DisAssociateHost(ctx context.Context, hostgroupID, hostID int) error {
	spath := "/host/associate"

	req, err := d.newRequest(ctx, "DELETE", spath, nil)
	if err != nil {
		return fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	q := req.URL.Query()
	q.Add("ID", strconv.Itoa(hostgroupID))
	q.Add("ASSOCIATEOBJID", strconv.Itoa(hostID))
	q.Add("ASSOCIATEOBJTYPE", strconv.Itoa(TypeHost))
	q.Add("TYPE", strconv.Itoa(TypeHostGroup))
	req.URL.RawQuery = q.Encode()

	var i interface{} // this endpoint return N/A
	if err = d.requestWithRetry(req, i, DefaultHTTPRetryCount); err != nil {
		return fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return nil
}

// CreateHostGroupWithHost create hostgroup object and host (same name) object
// and associate host.
func (d *Device) CreateHostGroupWithHost(ctx context.Context, hostname string) (*HostGroup, *Host, error) {
	host, err := d.CreateHost(ctx, hostname)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Host: %w", err)
	}

	hostgroup, err := d.CreateHostGroup(ctx, hostname)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create hostgroup: %w", err)
	}

	err = d.AssociateHost(ctx, hostgroup.ID, host.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to associate to hostgroup: %w", err)
	}

	return hostgroup, host, nil
}

// DeleteHostGroupWithHost delete hostgroup object and host (same name) object.
func (d *Device) DeleteHostGroupWithHost(ctx context.Context, hostgroupID int) error {
	hostgroup, err := d.GetHostGroup(ctx, hostgroupID)
	if err != nil {
		return fmt.Errorf("failed to search hostgroup by ID: %w", err)
	}
	hosts, err := d.GetHosts(ctx, NewSearchQueryHostname(hostgroup.NAME))
	if err != nil {
		return fmt.Errorf("failed to search host: %w", err)
	}
	if len(hosts) != 1 {
		return errors.New("search result of host is not one")
	}
	host := hosts[0]

	err = d.DisAssociateHost(ctx, hostgroup.ID, host.ID)
	if err != nil {
		return fmt.Errorf("failed to deassociate hostgroup: %w", err)
	}
	err = d.DeleteHost(ctx, host.ID)
	if err != nil {
		return fmt.Errorf("failed to delete host: %w", err)
	}
	err = d.DeleteHostGroup(ctx, hostgroup.ID)
	if err != nil {
		return fmt.Errorf("failed to delete hostgroup: %w", err)
	}

	return nil
}

// GetHostGroupForce get hostgroup object and host object.
// create hostgroup and host object if not exists.
func (d *Device) GetHostGroupForce(ctx context.Context, hostname string) (*HostGroup, *Host, error) {
	// GetHostGroup and CreateHostGroup if not found.
	hostgroups, err := d.GetHostGroups(ctx, NewSearchQueryHostname(hostname))
	if err != nil {
		if err == ErrHostGroupNotFound {
			return d.CreateHostGroupWithHost(ctx, hostname)
		}

		// Unexpected Error
		return nil, nil, fmt.Errorf("failed to get hostgroup: %w", err)
	}

	if len(hostgroups) != 1 {
		// hostgroup is must be unique
		return nil, nil, fmt.Errorf("fount multiple hostgroup in same hostname (hostname: %s)", hostname)
	}
	hostgroup := hostgroups[0]

	hosts, err := d.GetHosts(ctx, NewSearchQueryHostname(hostname))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get host: %w", err)
	}

	// host : hostgroup is 1:1, if get not only one, data is incorrect!
	if len(hosts) != 1 {
		return nil, nil, fmt.Errorf("found multiple hosts associated hostgroup: %w", err)
	}
	host := hosts[0]

	if host.ISADD2HOSTGROUP == false {
		err = d.AssociateHost(ctx, hostgroup.ID, host.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to associate host to hostgroup: %w", err)
		}
	}

	return &hostgroup, &host, nil
}
