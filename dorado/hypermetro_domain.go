package dorado

import (
	"context"
	"fmt"
)

// NOTE(whywaita): implement only GET.
// HyperMetroDomain is a few under our usage.

// HyperMetroDomain is domain of HyperMetro
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

// GetHyperMetroDomains get HyperMetroDomain objects.
func (c *Client) GetHyperMetroDomains(ctx context.Context, query *SearchQuery) ([]HyperMetroDomain, error) {
	// HyperMetroDomain is a same value between a local device and a remote device.
	return c.LocalDevice.GetHyperMetroDomains(ctx, query)
}

// GetHyperMetroDomains get HyperMetroDomain objects in device.
func (d *Device) GetHyperMetroDomains(ctx context.Context, query *SearchQuery) ([]HyperMetroDomain, error) {
	spath := "/HyperMetroDomain"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req = AddSearchQuery(req, query)

	var hyperMetroDomains []HyperMetroDomain
	if err = d.requestWithRetry(req, &hyperMetroDomains, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	if len(hyperMetroDomains) == 0 {
		return nil, ErrHyperMetroDomainNotFound
	}

	return hyperMetroDomains, nil
}
