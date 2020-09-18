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

// Initiator is iSCSI initiator
type Initiator struct {
	FAILOVERMODE    string `json:"FAILOVERMODE"`
	HEALTHSTATUS    string `json:"HEALTHSTATUS"`
	ID              string `json:"ID"` // = IQN
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

func encodeIqn(iqn string) string {
	// must escape colon when using filter string.
	return strings.ReplaceAll(iqn, `:`, `\:`)
}

// GetInitiators search initiators.
// you must use encodeIqn when to search iqn.
// ex: initiators, err := d.GetInitiators(ctx, NewSearchQueryId(encodeIqn(iqn)))
func (d *Device) GetInitiators(ctx context.Context, query *SearchQuery) ([]Initiator, error) {
	spath := "/iscsi_initiator"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req = AddSearchQuery(req, query)

	var initiators []Initiator
	if err = d.requestWithRetry(req, &initiators, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	if len(initiators) == 0 {
		return nil, ErrInitiatorNotFound
	}

	return initiators, nil
}

// GetInitiator get initiator by id.
func (d *Device) GetInitiator(ctx context.Context, iqn string) (*Initiator, error) {
	spath := fmt.Sprintf("/iscsi_initiator/%s", iqn)

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	initiators := &Initiator{}
	if err = d.requestWithRetry(req, initiators, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return initiators, nil
}

// CreateInitiator create initiator.
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
		return nil, fmt.Errorf(ErrCreatePostValue+": %w", err)
	}

	req, err := d.newRequest(ctx, "POST", spath, bytes.NewBuffer(jb))
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	initiator := &Initiator{}
	if err = d.requestWithRetry(req, initiator, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return initiator, nil
}

// DeleteInitiator delete initiator.
func (d *Device) DeleteInitiator(ctx context.Context, iqn string) error {
	spath := fmt.Sprintf("/iscsi_initiator/%s", iqn)

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

// UpdateInitiatorParam is parameter for UpdateInitiator
type UpdateInitiatorParam struct {
	USECHAP    string `json:"USECHAP"`
	PARENTTYPE string `json:"PARENTTYPE"`
	TYPE       string `json:"TYPE"`
	ID         string `json:"ID"`
	PARENTID   string `json:"PARENTID"`
}

// UpdateInitiator update initiator information.
func (d *Device) UpdateInitiator(ctx context.Context, iqn string, initiatorParam UpdateInitiatorParam) (*Initiator, error) {
	spath := fmt.Sprintf("/iscsi_initiator/%s", iqn)

	jb, err := json.Marshal(initiatorParam)
	if err != nil {
		return nil, fmt.Errorf(ErrCreatePostValue+": %w", err)
	}

	req, err := d.newRequest(ctx, "PUT", spath, bytes.NewBuffer(jb))
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	initiator := &Initiator{}
	if err = d.requestWithRetry(req, initiator, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return initiator, nil
}

// GetInitiatorForce get initiator and create initiator if not exists.
func (d *Device) GetInitiatorForce(ctx context.Context, iqn string) (*Initiator, error) {
	initiators, err := d.GetInitiators(ctx, NewSearchQueryID(encodeIqn(iqn)))
	if err != nil {
		if err == ErrInitiatorNotFound {
			return d.CreateInitiator(ctx, iqn)
		}

		return nil, fmt.Errorf("failed to get initiators: %w", err)
	}

	if len(initiators) != 1 {
		return nil, errors.New("fount multiple initiators in same iqn")
	}

	return &initiators[0], nil
}
