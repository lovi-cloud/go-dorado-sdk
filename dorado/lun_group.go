package dorado

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type LunGroup struct {
	CAPCITY            string `json:"CAPCITY"`
	DESCRIPTION        string `json:"DESCRIPTION"`
	ID                 string `json:"ID"`
	ISADD2MAPPINGVIEW  string `json:"ISADD2MAPPINGVIEW"`
	NAME               string `json:"NAME"`
	SMARTQOSPOLICYID   string `json:"SMARTQOSPOLICYID"`
	TYPE               int    `json:"TYPE"`
	ASSOCIATELUNIDLIST string `json:"ASSOCIATELUNIDLIST"`
}

const (
	ErrLunGroupNotFound = "LUN Group is not found"
)

func (d *Device) GetLunGroups(ctx context.Context, query *SearchQuery) ([]LunGroup, error) {
	spath := "/lungroup"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	req = AddSearchQuery(req, query)
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	lunGroups := []LunGroup{}
	if err = decodeBody(resp, &lunGroups); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	if len(lunGroups) == 0 {
		return nil, errors.New(ErrLunGroupNotFound)
	}

	return lunGroups, nil
}

func (d *Device) GetLunGroup(ctx context.Context, lungroupId string) (*LunGroup, error) {
	spath := fmt.Sprintf("/lungroup/%s", lungroupId)

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	lunGroup := &LunGroup{}
	if err = decodeBody(resp, lunGroup); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return lunGroup, nil
}

func (d *Device) CreateLunGroup(ctx context.Context, hostname string) (*LunGroup, error) {
	// Host : HostGroup : LunGroup is 1:1:1.
	// lungroup will create the same name as a host.
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

	lunGroup := &LunGroup{}
	if err = decodeBody(resp, lunGroup); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return lunGroup, nil
}

func (d *Device) DeleteLunGroup(ctx context.Context, lungroupId string) error {
	spath := fmt.Sprintf("/lungroup/%s", lungroupId)

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

func (d *Device) AssociateLun(ctx context.Context, lungroupId, lunId string) error {
	spath := "/lungroup/associate"
	param := AssociateParam{
		ID:               lungroupId,
		ASSOCIATEOBJID:   lunId,
		ASSOCIATEOBJTYPE: TypeLUN,
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
	if err = decodeBody(resp, i); err != nil {
		return errors.Wrap(err, ErrDecodeBody)
	}

	return nil
}

func (d *Device) DisAssociateLun(ctx context.Context, lungroupId, lunId string) error {
	spath := "/lungroup/associate"
	param := AssociateParam{
		ID:               lungroupId,
		ASSOCIATEOBJID:   lunId,
		ASSOCIATEOBJTYPE: TypeLUN,
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
	if err = decodeBody(resp, i); err != nil {
		return errors.Wrap(err, ErrDecodeBody)
	}

	return nil
}

func (d *Device) GetLunGroupByLunId(ctx context.Context, lunId string) (*LunGroup, error) {
	query := &SearchQuery{
		AssociateObjType: strconv.Itoa(TypeLUN),
		AssociateObjID:   lunId,
	}

	lungroups, err := d.GetLunGroups(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get lun group")
	}

	return &lungroups[0], nil
}

func (lg *LunGroup) IsAssociated() bool {
	list := lg.ASSOCIATELUNIDLIST

	if list != "" {
		return true
	}

	return false
}

func (d *Device) GetLunGroupForce(ctx context.Context, hostname string) (*LunGroup, error) {
	lungroups, err := d.GetLunGroups(ctx, NewSearchQueryHostname(hostname))
	if err != nil {
		if err.Error() == ErrLunGroupNotFound {
			return d.CreateLunGroup(ctx, hostname)
		}

		return nil, errors.Wrap(err, "failed to get lungroup")
	}

	if len(lungroups) != 1 {
		return nil, errors.New("found multiple lungroup in same hostname")
	}

	return &lungroups[0], nil
}
