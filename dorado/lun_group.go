package dorado

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// LunGroup is group of LUN
type LunGroup struct {
	CAPCITY            string `json:"CAPCITY"`
	DESCRIPTION        string `json:"DESCRIPTION"`
	ID                 int    `json:"ID,string"`
	ISADD2MAPPINGVIEW  bool   `json:"ISADD2MAPPINGVIEW,string"`
	NAME               string `json:"NAME"`
	SMARTQOSPOLICYID   string `json:"SMARTQOSPOLICYID"`
	TYPE               int    `json:"TYPE"`
	ASSOCIATELUNIDLIST string `json:"ASSOCIATELUNIDLIST"`
}

// GetLunGroups get lun groups by query
func (d *Device) GetLunGroups(ctx context.Context, query *SearchQuery) ([]LunGroup, error) {
	spath := "/lungroup"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req = AddSearchQuery(req, query)

	var lunGroups []LunGroup
	if err = d.requestWithRetry(req, &lunGroups, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	if len(lunGroups) == 0 {
		return nil, ErrLunGroupNotFound
	}

	return lunGroups, nil
}

// GetLunGroup get lun group by id
func (d *Device) GetLunGroup(ctx context.Context, lungroupID int) (*LunGroup, error) {
	spath := fmt.Sprintf("/lungroup/%d", lungroupID)

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	lunGroup := &LunGroup{}
	if err = d.requestWithRetry(req, lunGroup, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return lunGroup, nil
}

// CreateLunGroup create lun group
// Host : HostGroup : LunGroup is 1:1:1.
// lun group will create the same name as a host.
func (d *Device) CreateLunGroup(ctx context.Context, hostname string) (*LunGroup, error) {
	spath := "/lungroup"
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

	lunGroup := &LunGroup{}
	if err = d.requestWithRetry(req, lunGroup, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return lunGroup, nil
}

// DeleteLunGroup delete lun group
func (d *Device) DeleteLunGroup(ctx context.Context, lungroupID int) error {
	spath := fmt.Sprintf("/lungroup/%d", lungroupID)

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

// AssociateLun associate lun to lun group
func (d *Device) AssociateLun(ctx context.Context, lungroupID, lunID int) error {
	spath := "/lungroup/associate"
	param := AssociateParam{
		ID:               strconv.Itoa(lungroupID),
		ASSOCIATEOBJID:   strconv.Itoa(lunID),
		ASSOCIATEOBJTYPE: TypeLUN,
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

// DisAssociateLun dis associate lun from lun group
func (d *Device) DisAssociateLun(ctx context.Context, lungroupID, lunID int) error {
	spath := "/lungroup/associate"
	param := &AssociateParam{
		ID:               strconv.Itoa(lungroupID),
		ASSOCIATEOBJID:   strconv.Itoa(lunID),
		ASSOCIATEOBJTYPE: TypeLUN,
	}

	req, err := d.newRequest(ctx, "DELETE", spath, nil)
	if err != nil {
		return fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req = AddAssociateParam(req, param)

	var i interface{} // this endpoint return N/A
	if err = d.requestWithRetry(req, i, DefaultHTTPRetryCount); err != nil {
		return fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return nil
}

// GetAssociateLunGroups get associated lun group by query
func (d *Device) GetAssociateLunGroups(ctx context.Context, query *SearchQuery) ([]LunGroup, error) {
	spath := "/lungroup/associate"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req = AddSearchQuery(req, query)

	var lunGroups []LunGroup
	if err = d.requestWithRetry(req, &lunGroups, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return lunGroups, nil
}

// GetLunGroupByLunID get associated lun group by lun id.
func (d *Device) GetLunGroupByLunID(ctx context.Context, lunID int) (*LunGroup, error) {
	query := &SearchQuery{
		AssociateObjType: strconv.Itoa(TypeLUN),
		AssociateObjID:   strconv.Itoa(lunID),
		Type:             strconv.Itoa(TypeLUNGroup),
	}

	lungroups, err := d.GetAssociateLunGroups(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get lun group: %w", err)
	}
	if len(lungroups) != 1 {
		return nil, fmt.Errorf("found multiple LUN Group in same lun id: %w", err)
	}

	return &lungroups[0], nil
}

// IsAssociated return boolean
func (lg *LunGroup) IsAssociated() bool {
	list := lg.ASSOCIATELUNIDLIST

	if list != "" {
		return true
	}

	return false
}

// GetLunGroupForce get lun group, and create lun group if not exist.
func (d *Device) GetLunGroupForce(ctx context.Context, hostname string) (*LunGroup, error) {
	lungroups, err := d.GetLunGroups(ctx, NewSearchQueryHostname(hostname))
	if err != nil {
		if err == ErrLunGroupNotFound {
			return d.CreateLunGroup(ctx, hostname)
		}

		return nil, fmt.Errorf("failed to get lungroup: %w", err)
	}

	if len(lungroups) != 1 {
		return nil, errors.New("found multiple lungroup in same hostname")
	}

	return &lungroups[0], nil
}
