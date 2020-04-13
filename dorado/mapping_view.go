package dorado

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type MappingView struct {
	DESCRIPTION         string `json:"DESCRIPTION"`
	ENABLEINBANDCOMMAND string `json:"ENABLEINBANDCOMMAND"`
	ID                  string `json:"ID"`
	INBANDLUNWWN        string `json:"INBANDLUNWWN"`
	NAME                string `json:"NAME"`
	TYPE                int    `json:"TYPE"`
}

const (
	ErrMappingViewNotFound = "mapping view is not found"
)

func (d *Device) GetMappingViews(ctx context.Context, query *SearchQuery) ([]MappingView, error) {
	spath := "/mappingview"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	req = AddSearchQuery(req, query)

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	mappingviews := []MappingView{}
	if err = decodeBody(resp, &mappingviews); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	if len(mappingviews) == 0 {
		return nil, errors.New(ErrMappingViewNotFound)
	}

	return mappingviews, nil
}

func (d *Device) GetMappingView(ctx context.Context, mappingviewId string) (*MappingView, error) {
	spath := fmt.Sprintf("/mappingview/%s", mappingviewId)

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	mappingview := &MappingView{}
	if err = decodeBody(resp, mappingview); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return mappingview, nil
}

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

	mappingview := &MappingView{}
	if err = decodeBody(resp, mappingview); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return mappingview, nil
}

func (d *Device) DeleteMappingView(ctx context.Context, mappingviewId string) error {
	spath := fmt.Sprintf("/mappingview/%s", mappingviewId)

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

func (d *Device) AssociateMappingView(ctx context.Context, param AssociateParam) error {
	spath := "/mappingview/create_associate"

	jb, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err, ErrCreatePostValue)
	}
	req, err := d.newRequest(ctx, "PUT", spath, bytes.NewBuffer(jb))
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

func (d *Device) DisAssociateMappingView(ctx context.Context, param AssociateParam) error {
	spath := "mappingview/remove_associate"

	jb, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err, ErrCreatePostValue)
	}
	req, err := d.newRequest(ctx, "PUT", spath, bytes.NewBuffer(jb))
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

func (d *Device) GetMappingViewForce(ctx context.Context, hostname string) (*MappingView, error) {
	mappingviews, err := d.GetMappingViews(ctx, NewSearchQueryHostname(hostname))
	if err != nil {
		if err.Error() == ErrMappingViewNotFound {
			return d.CreateMappingView(ctx, hostname)
		}

		return nil, errors.Wrap(err, "failed to get mapping view")
	}

	if len(mappingviews) != 1 {
		return nil, errors.New("fount multiple mapping view in same hostname")
	}

	return &mappingviews[0], nil
}

func (d *Device) DoMapping(ctx context.Context, mappingview *MappingView, hostgroup *HostGroup, lungroup *LunGroup, portgroupId string) error {
	param := AssociateParam{
		ID:   mappingview.ID,
		TYPE: strconv.Itoa(TypeMappingView),
	}

	if hostgroup.ISADD2MAPPINGVIEW == "false" {
		param.ASSOCIATEOBJTYPE = TypeHostGroup
		param.ASSOCIATEOBJID = hostgroup.ID
		err := d.AssociateMappingView(ctx, param)
		if err != nil {
			return errors.Wrap(err, "failed to associate hostgroup")
		}
	}

	if lungroup.ISADD2MAPPINGVIEW == "false" {
		param.ASSOCIATEOBJTYPE = TypeLUNGroup
		param.ASSOCIATEOBJID = lungroup.ID
		err := d.AssociateMappingView(ctx, param)
		if err != nil {
			return errors.Wrap(err, "failed to associate lungroup")
		}
	}

	isExist, err := d.IsAddToMappingViewPortGroup(ctx, mappingview.ID, portgroupId)
	if err != nil {
		return errors.Wrap(err, "failed to get portgroup")
	}
	if isExist == false {
		param.ASSOCIATEOBJTYPE = TypePortGroup
		param.ASSOCIATEOBJID = portgroupId
		err := d.AssociateMappingView(ctx, param)
		if err != nil {
			return errors.Wrap(err, "failed to associate portgroup")
		}
	}

	return nil
}
