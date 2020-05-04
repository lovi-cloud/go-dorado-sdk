package dorado

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type PortGroup struct {
	DESCRIPTION string `json:"DESCRIPTION"`
	ID          string `json:"ID"`
	NAME        string `json:"NAME"`
	TYPE        int    `json:"TYPE"`
}

const (
	PortGroupNotFound = "PortGroup is not found"
)

func (d *Device) GetPortGroups(ctx context.Context, query *SearchQuery) ([]PortGroup, error) {
	spath := "/portgroup"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req = AddSearchQuery(req, query)

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	portgroups := []PortGroup{}
	if err = decodeBody(resp, &portgroups); err != nil {
		return nil, fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	if len(portgroups) == 0 {
		return nil, errors.New(PortGroupNotFound)
	}

	return portgroups, nil
}

func (d *Device) GetPortGroup(ctx context.Context, portgroupId string) (*PortGroup, error) {
	spath := fmt.Sprintf("/portgroup/%s", portgroupId)

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	portgroup := &PortGroup{}
	if err = decodeBody(resp, portgroup); err != nil {
		return nil, fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	return portgroup, nil
}

func (d *Device) GetPortGroupsAssociate(ctx context.Context, mappingviewId string) ([]PortGroup, error) {
	spath := "/portgroup/associate"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	param := &AssociateParam{
		ASSOCIATEOBJID:   mappingviewId,
		ASSOCIATEOBJTYPE: TypeMappingView,
	}
	req = AddAssociateParam(req, param)
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	portgroups := []PortGroup{}
	if err = decodeBody(resp, &portgroups); err != nil {
		return nil, fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	return portgroups, nil
}

func (d *Device) IsAddToMappingViewPortGroup(ctx context.Context, mappingViewId, portgroupId string) (bool, error) {
	portgroups, err := d.GetPortGroupsAssociate(ctx, mappingViewId)
	if err != nil {
		return false, fmt.Errorf("failed to get portgroups: %w", err)
	}

	for _, p := range portgroups {
		if p.ID == portgroupId {
			return true, nil
		}
	}

	return false, nil
}
