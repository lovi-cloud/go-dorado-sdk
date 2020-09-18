package dorado

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// TargetPort is target port (ex: iSCSI)
type TargetPort struct {
	ETHPORTID string `json:"ETHPORTID"`
	ID        string `json:"ID"`
	TPGT      string `json:"TPGT"`
	TYPE      int    `json:"TYPE"`
}

// GetTargetPort get target ports by query
func (d *Device) GetTargetPort(ctx context.Context, query *SearchQuery) ([]TargetPort, error) {
	spath := "/iscsi_tgt_port"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req = AddSearchQuery(req, query)

	var targetPorts []TargetPort
	if err = d.requestWithRetry(req, &targetPorts, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	if len(targetPorts) == 0 {
		return nil, ErrTargetPortNotFound
	}

	return targetPorts, nil
}

// GetTargetIQNs get target IQN
func (d *Device) GetTargetIQNs(ctx context.Context) ([]string, error) {
	targetports, err := d.GetTargetPort(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get target ports: %w", err)
	}

	var targetIqns []string
	for _, targetport := range targetports {
		iqn, err := d.parseTargetPortID(targetport.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse target port ID: %w", err)
		}
		if strings.HasPrefix(iqn, "iqn") == false {
			return nil, fmt.Errorf("invalid target IQN: %s", iqn)
		}

		targetIqns = append(targetIqns, iqn)
	}

	return targetIqns, nil
}

// parseTargetPortID parse TargetPort.ID. return target IQN.
// ex: 0+iqn.2006-08.com.huawei:oceanstor:name1:192.0.2.10,t,0x0001
func (d *Device) parseTargetPortID(id string) (string, error) {
	s := strings.Split(id, "+")
	if len(s) != 2 {
		return "", errors.New("splited length is not 2 (separator is +)")
	}

	s2 := strings.Split(s[1], ",")
	if len(s2) != 3 {
		return "", errors.New("splited length is not 2 (separator is ,)")
	}
	iqn := s2[0]

	return iqn, nil
}
