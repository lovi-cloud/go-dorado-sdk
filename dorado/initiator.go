package dorado

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Initiator struct {
	FAILOVERMODE    string `json:"FAILOVERMODE"`
	HEALTHSTATUS    string `json:"HEALTHSTATUS"`
	ID              string `json:"ID"`
	ISFREE          string `json:"ISFREE"`
	MULTIPATHTYPE   string `json:"MULTIPATHTYPE"`
	OPERATIONSYSTEM string `json:"OPERATIONSYSTEM"`
	PATHTYPE        string `json:"PATHTYPE"`
	RUNNINGSTATUS   string `json:"RUNNINGSTATUS"`
	SPECIALMODETYPE string `json:"SPECIALMODETYPE"`
	TYPE            int    `json:"TYPE"`
	USECHAP         string `json:"USECHAP"`
	PARENTID        string `json:"PARENTID,omitempty"`
	PARENTNAME      string `json:"PARENTNAME,omitempty"`
	PARENTTYPE      int    `json:"PARENTTYPE,omitempty"`
}

const (
	ErrInitiatorNotFound = "Initiator is not found"
)

func encodeIqn(iqn string) string {
	// must escape colon when using filter string.
	return strings.ReplaceAll(iqn, `:`, `\:`)
}

func (d *Device) GetInitiators(ctx context.Context, query *SearchQuery) ([]Initiator, error) {
	spath := "/iscsi_initiator"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	req = AddSearchQuery(req, query)
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	initiators := []Initiator{}
	if err = decodeBody(resp, initiators); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	if len(initiators) == 0 {
		return nil, errors.New(ErrInitiatorNotFound)
	}

	return initiators, nil
}

func (d *Device) GetInitiator(ctx context.Context, iqn string) (*Initiator, error) {
	spath := fmt.Sprintf("/iscsi_initiator/%s", iqn)

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	initiators := &Initiator{}
	if err = decodeBody(resp, initiators); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return initiators, nil
}

func (d *Device) CreateInitiator(ctx context.Context, iqn string) (*Initiator, error) {
	spath := "/iscsi_initiator"
	param := struct {
		USECHAP string `json:"USECHAP"`
		TYPE    string `json:"TYPE"`
		ID      string `json:"ID"`
	}{
		USECHAP: "false",
		TYPE:    strconv.Itoa(TypeInitiator),
		ID:      iqn,
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

	initiator := &Initiator{}
	if err = decodeBody(resp, initiator); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return initiator, nil
}

func (d *Device) DeleteInitiator(ctx context.Context, iqn string) error {
	spath := fmt.Sprintf("/iscsi_initiator/%s", iqn)

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

type UpdateInitiatorParam struct {
	USECHAP    string `json:"USECHAP"`
	PARENTTYPE string `json:"PARENTTYPE"`
	TYPE       string `json:"TYPE"`
	ID         string `json:"ID"`
	PARENTID   string `json:"PARENTID"`
}

func (d *Device) UpdateInitiator(ctx context.Context, iqn string, initiatorParam UpdateInitiatorParam) (*Initiator, error) {
	spath := fmt.Sprintf("/iscsi_initiator/%s", iqn)

	jb, err := json.Marshal(initiatorParam)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreatePostValue)
	}

	req, err := d.newRequest(ctx, "PUT", spath, bytes.NewBuffer(jb))
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	i := &Initiator{}
	if err = decodeBody(resp, i); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return i, nil
}

func (d *Device) GetInitiatorForce(ctx context.Context, iqn string) (*Initiator, error) {
	initiators, err := d.GetInitiators(ctx, NewSearchQueryName(encodeIqn(iqn)))
	if err != nil {
		if err.Error() == ErrInitiatorNotFound {
			return d.CreateInitiator(ctx, iqn)
		}

		return nil, errors.Wrap(err, "failed to get initiators")
	}

	if len(initiators) != 1 {
		return nil, errors.New("fount multiple initiators in same iqn")
	}

	return &initiators[0], nil
}
