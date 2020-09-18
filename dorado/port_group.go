package dorado

import (
	"context"
	"fmt"
	"strconv"
)

// PortGroup is group of Port (ex Ethernet, FiberChannel...)
type PortGroup struct {
	DESCRIPTION string `json:"DESCRIPTION"`
	ID          int    `json:"ID,string"`
	NAME        string `json:"NAME"`
	TYPE        int    `json:"TYPE"`
}

// GetPortGroups get port groups by query
func (d *Device) GetPortGroups(ctx context.Context, query *SearchQuery) ([]PortGroup, error) {
	spath := "/portgroup"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req = AddSearchQuery(req, query)

	var portGroups []PortGroup
	if err = d.requestWithRetry(req, &portGroups, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	if len(portGroups) == 0 {
		return nil, ErrPortGroupNotFound
	}

	return portGroups, nil
}

// GetPortGroup get port group by id
func (d *Device) GetPortGroup(ctx context.Context, portgroupID int) (*PortGroup, error) {
	spath := fmt.Sprintf("/portgroup/%d", portgroupID)

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	portGroup := &PortGroup{}
	if err = d.requestWithRetry(req, portGroup, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return portGroup, nil
}

// GetPortGroupsAssociate get port group that associated by mapping view id
func (d *Device) GetPortGroupsAssociate(ctx context.Context, mappingviewID int) ([]PortGroup, error) {
	spath := "/portgroup/associate"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	param := &AssociateParam{
		ASSOCIATEOBJID:   strconv.Itoa(mappingviewID),
		ASSOCIATEOBJTYPE: TypeMappingView,
	}
	req = AddAssociateParam(req, param)

	var portGroups []PortGroup
	if err = d.requestWithRetry(req, &portGroups, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return portGroups, nil
}

// IsAddToMappingViewPortGroup check to associated mapping view
func (d *Device) IsAddToMappingViewPortGroup(ctx context.Context, mappingViewID, portgroupID int) (bool, error) {
	portgroups, err := d.GetPortGroupsAssociate(ctx, mappingViewID)
	if err != nil {
		return false, fmt.Errorf("failed to get portgroups: %w", err)
	}

	for _, p := range portgroups {
		if p.ID == portgroupID {
			return true, nil
		}
	}

	return false, nil
}
