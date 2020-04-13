package dorado

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type HyperMetroPairParam struct {
	RECONVERYPOLICY string `json:"RECONVERYPOLICY"`
	DOMAINID        string `json:"DOMAINID"`
	SPEED           int    `json:"SPEED"`
	HCRESOURCETYPE  string `json:"HCRESOURCETYPE"`
	REMOTEOBJID     string `json:"REMOTEOBJID"`
	LOCALOBJID      string `json:"LOCALOBJID"`
	ISFIRSTSYNC     bool   `json:"ISFIRSTSYNC"`
}

type HyperMetroPair struct {
	CAPACITYBYTE             string `json:"CAPACITYBYTE"`
	CGID                     string `json:"CGID"`
	CGNAME                   string `json:"CGNAME"`
	DOMAINID                 string `json:"DOMAINID"`
	DOMAINNAME               string `json:"DOMAINNAME"`
	ENDTIME                  string `json:"ENDTIME"`
	HCRESOURCETYPE           string `json:"HCRESOURCETYPE"`
	HDRINGID                 string `json:"HDRINGID"`
	HEALTHSTATUS             string `json:"HEALTHSTATUS"`
	ID                       string `json:"ID"`
	ISINCG                   string `json:"ISINCG"`
	ISISOLATION              string `json:"ISISOLATION"`
	ISISOLATIONTHRESHOLDTIME string `json:"ISISOLATIONTHRESHOLDTIME"`
	ISPRIMARY                string `json:"ISPRIMARY"`
	LINKSTATUS               string `json:"LINKSTATUS"`
	LOCALDATASTATE           string `json:"LOCALDATASTATE"`
	LOCALHOSTACCESSSTATE     string `json:"LOCALHOSTACCESSSTATE"`
	LOCALOBJID               string `json:"LOCALOBJID"`
	LOCALOBJNAME             string `json:"LOCALOBJNAME"`
	RECOVERYPOLICY           string `json:"RECOVERYPOLICY"`
	REMOTEDATASTATE          string `json:"REMOTEDATASTATE"`
	REMOTEHOSTACCESSSTATE    string `json:"REMOTEHOSTACCESSSTATE"`
	REMOTEOBJID              string `json:"REMOTEOBJID"`
	REMOTEOBJNAME            string `json:"REMOTEOBJNAME"`
	RESOURCEWWN              string `json:"RESOURCEWWN"`
	RUNNINGSTATUS            string `json:"RUNNINGSTATUS"`
	SPEED                    string `json:"SPEED"`
	STARTTIME                string `json:"STARTTIME"`
	SYNCDIRECTION            string `json:"SYNCDIRECTION"`
	SYNCLEFTTIME             string `json:"SYNCLEFTTIME"`
	SYNCPROGRESS             string `json:"SYNCPROGRESS"`
	TYPE                     int    `json:"TYPE"`
	WRITESECONDARYTIMEOUT    string `json:"WRITESECONDARYTIMEOUT"`
}

func (c *Client) GetHyperMetroPairs(ctx context.Context, query *SearchQuery) ([]HyperMetroPair, error) {
	spath := "/HyperMetroPair"

	req, err := c.LocalDevice.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	if query == nil {
		query = &SearchQuery{
			Range: "[0-4095]", // NOTE(whywaita): if set range, response become fast and not duplicated
		}
	}
	req = AddSearchQuery(req, query)

	resp, err := c.LocalDevice.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	hyperMetroPairs := []HyperMetroPair{}
	if err = decodeBody(resp, &hyperMetroPairs); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return hyperMetroPairs, nil
}

func (c *Client) GetHyperMetroPair(ctx context.Context, hyperMetroPairId string) (*HyperMetroPair, error) {
	spath := fmt.Sprintf("/HyperMetroPair/%s", hyperMetroPairId)

	req, err := c.LocalDevice.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := c.LocalDevice.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	hyperMetroPair := &HyperMetroPair{}
	if err = decodeBody(resp, hyperMetroPair); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return hyperMetroPair, nil
}

func (c *Client) CreateHyperMetroPair(ctx context.Context, hyperMetroDomainId, localLunId, remoteLunId string) (*HyperMetroPair, error) {
	spath := "/HyperMetroPair"
	param := &HyperMetroPairParam{
		RECONVERYPOLICY: "1",
		DOMAINID:        hyperMetroDomainId,
		SPEED:           2,
		HCRESOURCETYPE:  "1",
		REMOTEOBJID:     remoteLunId,
		LOCALOBJID:      localLunId,
		ISFIRSTSYNC:     false,
	}

	jb, err := json.Marshal(param)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreatePostValue)
	}
	req, err := c.LocalDevice.newRequest(ctx, "POST", spath, bytes.NewBuffer(jb))
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := c.LocalDevice.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	hyperMetroPair := &HyperMetroPair{}
	if err = decodeBody(resp, hyperMetroPair); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return hyperMetroPair, nil
}

func (c *Client) DeleteHyperMetroPair(ctx context.Context, hyperMetroPairId string) error {
	// must be suspend HyperMetro Pair before call this method.
	spath := fmt.Sprintf("/HyperMetroPair/%s", hyperMetroPairId)

	req, err := c.LocalDevice.newRequest(ctx, "DELETE", spath, nil)
	if err != nil {
		return errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := c.LocalDevice.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, ErrHTTPRequestDo)
	}

	var i interface{} // this endpoint return N/A
	if err = decodeBody(resp, i); err != nil {
		return errors.Wrap(err, ErrDecodeBody)
	}

	return nil
}

func (c *Client) SuspendHyperMetroPair(ctx context.Context, hyperMetroPairId string) error {
	spath := "/HyperMetroPair/disable_hcpair"
	param := struct {
		ID   string `json:"ID"`
		TYPE string `json:"TYPE"`
	}{
		ID:   hyperMetroPairId,
		TYPE: strconv.Itoa(TypeHyperMetroPair),
	}
	jb, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err, ErrCreatePostValue)
	}

	req, err := c.LocalDevice.newRequest(ctx, "PUT", spath, bytes.NewBuffer(jb))
	if err != nil {
		return errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := c.LocalDevice.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, ErrHTTPRequestDo)
	}

	var i interface{} // this endpoint return N/A
	if err = decodeBody(resp, i); err != nil {
		return errors.Wrap(err, ErrDecodeBody)
	}

	return nil
}

func (c *Client) SyncHyperMetroPair(ctx context.Context, hyperMetroPairId string) error {
	spath := "/HyperMetroPair/synchronize_hcpair"
	param := struct {
		ID   string `json:"ID"`
		TYPE string `json:"TYPE"`
	}{
		ID:   hyperMetroPairId,
		TYPE: strconv.Itoa(TypeHyperMetroPair),
	}
	jb, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err, ErrCreatePostValue)
	}

	req, err := c.LocalDevice.newRequest(ctx, "PUT", spath, bytes.NewBuffer(jb))
	if err != nil {
		return errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := c.LocalDevice.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, ErrHTTPRequestDo)
	}

	var i interface{} // this endpoint return N/A
	if err = decodeBody(resp, i); err != nil {
		return errors.Wrap(err, ErrDecodeBody)
	}

	return nil
}
