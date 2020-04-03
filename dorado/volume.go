package dorado

import (
	"context"

	"github.com/pkg/errors"

	uuid "github.com/satori/go.uuid"
)

const (
	HyperMetroPairIsNotFound = "HyperMetroPair is not found"
)

func (c *Client) CreateVolume(ctx context.Context, name uuid.UUID, capacityGB int, storagePoolId, hyperMetroDomainId string) (*HyperMetroPair, error) {
	// create volume (= hypermetro enabled lun)
	localLun, err := c.LocalDevice.CreateLUN(ctx, name, capacityGB, storagePoolId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create lun in local device")
	}

	remoteLun, err := c.RemoteDevice.CreateLUN(ctx, name, capacityGB, storagePoolId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create lun in remote device")
	}

	hyperMetroPair, err := c.CreateHyperMetroPair(ctx, hyperMetroDomainId, localLun.ID, remoteLun.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HyperMetroPair")
	}

	return hyperMetroPair, nil
}

func (c *Client) DeleteVolume(ctx context.Context, hyperMetroPairId string) error {
	// 1: delete HyperMetro Pair
	// 2: delete LUN Group Associate
	// 3: delete LUN

	// TODO: 2

	// 1: delete HyperMetro Pair
	hyperMetroPair, err := c.GetHyperMetroPairs(ctx, CreateSearchId(hyperMetroPairId))
	if err != nil {
		return errors.Wrap(err, "failed to get HyperMetro Pair")
	}
	if len(hyperMetroPair) == 0 {
		return errors.Wrap(err, HyperMetroPairIsNotFound)
	}

	hmp := hyperMetroPair[0] // id is unique (maybe)
	err = c.SuspendHyperMetroPair(ctx, hmp.ID)
	if err != nil {
		return errors.Wrap(err, "failed to suspend HyperMetroPair")
	}
	err = c.DeleteHyperMetroPair(ctx, hmp.ID)
	if err != nil {
		return errors.Wrap(err, "failed to delete HyperMetroPair")
	}

	// 3: delete LUN
	err = c.LocalDevice.DeleteLUN(ctx, hmp.LOCALOBJID)
	if err != nil {
		return errors.Wrap(err, "failed to delete Local LUN")
	}
	err = c.RemoteDevice.DeleteLUN(ctx, hmp.REMOTEOBJID)
	if err != nil {
		return errors.Wrap(err, "failed to delete Remote LUN")
	}

	return nil
}
