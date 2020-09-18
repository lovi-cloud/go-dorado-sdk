package dorado

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// MappingView is mapping object for lun
type MappingView struct {
	DESCRIPTION         string `json:"DESCRIPTION"`
	ENABLEINBANDCOMMAND bool   `json:"ENABLEINBANDCOMMAND,string"`
	ID                  int    `json:"ID,string"`
	INBANDLUNWWN        string `json:"INBANDLUNWWN"`
	NAME                string `json:"NAME"`
	TYPE                int    `json:"TYPE"`
}

// GetMappingViews get mapping view objects by query
func (d *Device) GetMappingViews(ctx context.Context, query *SearchQuery) ([]MappingView, error) {
	spath := "/mappingview"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req = AddSearchQuery(req, query)

	mappingviews := []MappingView{}
	if err = d.requestWithRetry(req, &mappingviews, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	if len(mappingviews) == 0 {
		return nil, ErrMappingViewNotFound
	}

	return mappingviews, nil
}

// GetMappingView get mapping view object by id
func (d *Device) GetMappingView(ctx context.Context, mappingviewID int) (*MappingView, error) {
	spath := fmt.Sprintf("/mappingview/%d", mappingviewID)

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	mappingview := &MappingView{}
	if err = d.requestWithRetry(req, mappingview, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return mappingview, nil
}

// CreateMappingView create mapping view object
func (d *Device) CreateMappingView(ctx context.Context, hostname string) (*MappingView, error) {
	spath := "/mappingview"
	param := struct {
		TYPE string `json:"TYPE"`
		NAME string `json:"NAME"`
	}{
		TYPE: strconv.Itoa(TypeMappingView),
		NAME: encodeHostName(hostname),
	}
	jb, err := json.Marshal(param)
	if err != nil {
		return nil, fmt.Errorf(ErrCreatePostValue+": %w", err)
	}

	req, err := d.newRequest(ctx, "POST", spath, bytes.NewBuffer(jb))
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	mappingview := &MappingView{}
	if err = d.requestWithRetry(req, mappingview, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return mappingview, nil
}

// DeleteMappingView delete mapping view object
func (d *Device) DeleteMappingView(ctx context.Context, mappingviewID int) error {
	spath := fmt.Sprintf("/mappingview/%d", mappingviewID)

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

// AssociateMappingView associate object to mapping view
func (d *Device) AssociateMappingView(ctx context.Context, param AssociateParam) error {
	spath := "/mappingview/create_associate"

	jb, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf(ErrCreatePostValue+": %w", err)
	}
	req, err := d.newRequest(ctx, "PUT", spath, bytes.NewBuffer(jb))
	if err != nil {
		return fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	var i interface{} // this endpoint return N/A
	if err = d.requestWithRetry(req, i, DefaultHTTPRetryCount); err != nil {
		return fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return nil
}

// DisAssociateMappingView disassociate object from mapping view
func (d *Device) DisAssociateMappingView(ctx context.Context, param AssociateParam) error {
	spath := "/mappingview/remove_associate"

	jb, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf(ErrCreatePostValue+": %w", err)
	}
	req, err := d.newRequest(ctx, "PUT", spath, bytes.NewBuffer(jb))
	if err != nil {
		return fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	var i interface{} // this endpoint return N/A
	if err = d.requestWithRetry(req, i, DefaultHTTPRetryCount); err != nil {
		return fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}
	return nil
}

// GetMappingViewForce get mapping view object and create if not exist
func (d *Device) GetMappingViewForce(ctx context.Context, hostname string) (*MappingView, error) {
	mappingviews, err := d.GetMappingViews(ctx, NewSearchQueryHostname(hostname))
	if err != nil {
		if err == ErrMappingViewNotFound {
			return d.CreateMappingView(ctx, hostname)
		}

		return nil, fmt.Errorf("failed to get mapping view: %w", err)
	}

	if len(mappingviews) != 1 {
		return nil, errors.New("fount multiple mapping view in same hostname")
	}

	return &mappingviews[0], nil
}

// DoMapping do mapping hostgroup/lungroup/portgroup to mappingview id
func (d *Device) DoMapping(ctx context.Context, mappingview *MappingView, hostgroup *HostGroup, lungroup *LunGroup, portgroupID int) error {
	param := AssociateParam{
		ID:   strconv.Itoa(mappingview.ID),
		TYPE: strconv.Itoa(TypeMappingView),
	}

	if hostgroup.ISADD2MAPPINGVIEW == false {
		param.ASSOCIATEOBJTYPE = TypeHostGroup
		param.ASSOCIATEOBJID = strconv.Itoa(hostgroup.ID)
		err := d.AssociateMappingView(ctx, param)
		if err != nil {
			return fmt.Errorf("failed to associate hostgroup: %w", err)
		}
	}

	if lungroup.ISADD2MAPPINGVIEW == false {
		param.ASSOCIATEOBJTYPE = TypeLUNGroup
		param.ASSOCIATEOBJID = strconv.Itoa(lungroup.ID)
		err := d.AssociateMappingView(ctx, param)
		if err != nil {
			return fmt.Errorf("failed to associate lungroup: %w", err)
		}
	}

	isExist, err := d.IsAddToMappingViewPortGroup(ctx, mappingview.ID, portgroupID)
	if err != nil {
		return fmt.Errorf("failed to get portgroup: %w", err)
	}
	if isExist == false {
		param.ASSOCIATEOBJTYPE = TypePortGroup
		param.ASSOCIATEOBJID = strconv.Itoa(portgroupID)
		err := d.AssociateMappingView(ctx, param)
		if err != nil {
			return fmt.Errorf("failed to associate portgroup: %w", err)
		}
	}

	return nil
}
