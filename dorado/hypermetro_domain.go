package dorado

import (
	"context"

	"github.com/pkg/errors"
)

type HyperMetroDomain struct {
	CPSID          string `json:"CPSID"`
	CPSNAME        string `json:"CPSNAME"`
	CPTYPE         string `json:"CPTYPE"`
	DESCRIPTION    string `json:"DESCRIPTION"`
	DOMAINTYPE     string `json:"DOMAINTYPE"`
	ID             string `json:"ID"`
	NAME           string `json:"NAME"`
	REMOTEDEVICES  string `json:"REMOTEDEVICES"`
	RUNNINGSTATUS  string `json:"RUNNINGSTATUS"`
	STANDBYCPSID   string `json:"STANDBYCPSID"`
	STANDBYCPSNAME string `json:"STANDBYCPSNAME"`
	TYPE           int    `json:"TYPE"`
}

func (d *Device) GetHyperMetroDomain(ctx context.Context) ([]HyperMetroDomain, error) {
	// NOTE(whywaita): implement only GET.
	// HyperMetroDomain is a few under our usage.

	spath := "HyperMetroDomain"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	var hyperMetroDomains []HyperMetroDomain
	if err = decodeBody(resp, &hyperMetroDomains); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return hyperMetroDomains, nil
}
