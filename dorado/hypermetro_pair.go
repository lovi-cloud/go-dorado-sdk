package dorado

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// HyperMetroPairParam is parameter of CreateHyperMetroPair
type HyperMetroPairParam struct {
	RECONVERYPOLICY string `json:"RECONVERYPOLICY"`
	DOMAINID        string `json:"DOMAINID"`
	SPEED           int    `json:"SPEED"`
	HCRESOURCETYPE  string `json:"HCRESOURCETYPE"`
	REMOTEOBJID     string `json:"REMOTEOBJID"`
	LOCALOBJID      string `json:"LOCALOBJID"`
	ISFIRSTSYNC     bool   `json:"ISFIRSTSYNC"`
}

// HyperMetroPair is object of LUN (synced by HyperMetro)
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
	LOCALOBJID               int    `json:"LOCALOBJID,string"`
	LOCALOBJNAME             string `json:"LOCALOBJNAME"`
	RECOVERYPOLICY           string `json:"RECOVERYPOLICY"`
	REMOTEDATASTATE          string `json:"REMOTEDATASTATE"`
	REMOTEHOSTACCESSSTATE    string `json:"REMOTEHOSTACCESSSTATE"`
	REMOTEOBJID              int    `json:"REMOTEOBJID,string"`
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

// Error const
const (
	ErrHyperMetroPairNotFound = "HyperMetroPair is not found"
)

// GetHyperMetroPairs get HyperMetro objects by query
func (c *Client) GetHyperMetroPairs(ctx context.Context, query *SearchQuery) ([]HyperMetroPair, error) {
	spath := "/HyperMetroPair"

	req, err := c.LocalDevice.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	if query == nil {
		query = &SearchQuery{
			Range: "[0-4095]", // NOTE(whywaita): if set range, response become fast and not duplicated
		}
	}
	req = AddSearchQuery(req, query)

	resp, err := c.LocalDevice.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	hyperMetroPairs := []HyperMetroPair{}
	if err = decodeBody(resp, &hyperMetroPairs); err != nil {
		return nil, fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	if len(hyperMetroPairs) == 0 {
		return nil, errors.New(ErrHyperMetroPairNotFound)
	}

	return hyperMetroPairs, nil
}

// GetHyperMetroPair get HyperMetro object by id
func (c *Client) GetHyperMetroPair(ctx context.Context, hyperMetroPairID string) (*HyperMetroPair, error) {
	spath := fmt.Sprintf("/HyperMetroPair/%s", hyperMetroPairID)

	req, err := c.LocalDevice.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	resp, err := c.LocalDevice.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	hyperMetroPair := &HyperMetroPair{}
	if err = decodeBody(resp, hyperMetroPair); err != nil {
		return nil, fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	return hyperMetroPair, nil
}

// CreateHyperMetroPair create HyperMetroPair.
func (c *Client) CreateHyperMetroPair(ctx context.Context, hyperMetroDomainID string, localLunID, remoteLunID int) (*HyperMetroPair, error) {
	spath := "/HyperMetroPair"
	param := &HyperMetroPairParam{
		RECONVERYPOLICY: "1",
		DOMAINID:        hyperMetroDomainID,
		SPEED:           2,
		HCRESOURCETYPE:  "1",
		REMOTEOBJID:     strconv.Itoa(remoteLunID),
		LOCALOBJID:      strconv.Itoa(localLunID),
		ISFIRSTSYNC:     false,
	}

	jb, err := json.Marshal(param)
	if err != nil {
		return nil, fmt.Errorf(ErrCreatePostValue+": %w", err)
	}
	req, err := c.LocalDevice.newRequest(ctx, "POST", spath, bytes.NewBuffer(jb))
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	resp, err := c.LocalDevice.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	hyperMetroPair := &HyperMetroPair{}
	if err = decodeBody(resp, hyperMetroPair); err != nil {
		return nil, fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	return hyperMetroPair, nil
}

// DeleteHyperMetroPair delete HyperMetroPair.
// must be suspend HyperMetro Pair before call this method.
func (c *Client) DeleteHyperMetroPair(ctx context.Context, hyperMetroPairID string) error {
	spath := fmt.Sprintf("/HyperMetroPair/%s", hyperMetroPairID)

	req, err := c.LocalDevice.newRequest(ctx, "DELETE", spath, nil)
	if err != nil {
		return fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	resp, err := c.LocalDevice.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	var i interface{} // this endpoint return N/A
	if err = decodeBody(resp, i); err != nil {
		return fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	return nil
}

// SuspendHyperMetroPair suspend HyperMetro sync.
func (c *Client) SuspendHyperMetroPair(ctx context.Context, hyperMetroPairID string) error {
	spath := "/HyperMetroPair/disable_hcpair"
	param := struct {
		ID   string `json:"ID"`
		TYPE string `json:"TYPE"`
	}{
		ID:   hyperMetroPairID,
		TYPE: strconv.Itoa(TypeHyperMetroPair),
	}
	jb, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf(ErrCreatePostValue+": %w", err)
	}

	req, err := c.LocalDevice.newRequest(ctx, "PUT", spath, bytes.NewBuffer(jb))
	if err != nil {
		return fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	resp, err := c.LocalDevice.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	var i interface{} // this endpoint return N/A
	if err = decodeBody(resp, i); err != nil {
		return fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	return nil
}

// SyncHyperMetroPair start to sync HyperMetro.
func (c *Client) SyncHyperMetroPair(ctx context.Context, hyperMetroPairID string) error {
	spath := "/HyperMetroPair/synchronize_hcpair"
	param := struct {
		ID   string `json:"ID"`
		TYPE string `json:"TYPE"`
	}{
		ID:   hyperMetroPairID,
		TYPE: strconv.Itoa(TypeHyperMetroPair),
	}
	jb, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf(ErrCreatePostValue+": %w", err)
	}

	req, err := c.LocalDevice.newRequest(ctx, "PUT", spath, bytes.NewBuffer(jb))
	if err != nil {
		return fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	resp, err := c.LocalDevice.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	var i interface{} // this endpoint return N/A
	if err = decodeBody(resp, i); err != nil {
		return fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	return nil
}
