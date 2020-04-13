package dorado

import (
	"context"
	"strconv"

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

	// 1: delete HyperMetro Pair
	hmp, err := c.GetHyperMetroPair(ctx, hyperMetroPairId)
	if err != nil {
		return errors.Wrap(err, "failed to get HyperMetro Pair")
	}

	err = c.SuspendHyperMetroPair(ctx, hmp.ID)
	if err != nil {
		return errors.Wrap(err, "failed to suspend HyperMetroPair")
	}
	err = c.DeleteHyperMetroPair(ctx, hmp.ID)
	if err != nil {
		return errors.Wrap(err, "failed to delete HyperMetroPair")
	}

	// 2: delete LUN Group Associate
	llun, err := c.LocalDevice.GetLUN(ctx, hmp.LOCALOBJID)
	if err != nil {
		return errors.Wrap(err, "failed to get lun information")
	}
	if llun.ISADD2LUNGROUP == "true" {
		lLungroup, err := c.LocalDevice.GetLunGroupByLunId(ctx, hmp.LOCALOBJID)
		if err != nil {
			return errors.Wrap(err, "failed to get lungroup by associated lun")
		}
		err = c.LocalDevice.DisAssociateLun(ctx, lLungroup.ID, hmp.LOCALOBJID)
		if err != nil {
			return errors.Wrap(err, "failed to disassociate local lun")
		}
	}

	rlun, err := c.RemoteDevice.GetLUN(ctx, hmp.REMOTEOBJID)
	if err != nil {
		return errors.Wrap(err, "failed to get lun information")
	}
	if rlun.ISADD2LUNGROUP == "true" {
		rLungroup, err := c.RemoteDevice.GetLunGroupByLunId(ctx, hmp.REMOTEOBJID)
		if err != nil {
			return errors.Wrap(err, "failed to get lungroup by associated lun")
		}
		err = c.RemoteDevice.DisAssociateLun(ctx, rLungroup.ID, hmp.REMOTEOBJID)
		if err != nil {
			return errors.Wrap(err, "failed to disassociate remote lun")
		}
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

func (c *Client) ExtendVolume(ctx context.Context, hyperMetroPairId string, newVolumeSizeGb int) error {
	// 1: Suspend HyperMetro Pair
	// 2: Expand LUN
	// 3: Re-sync HyperMetro Pair

	// 1: Suspend HyperMetro Pair
	hmp, err := c.GetHyperMetroPair(ctx, hyperMetroPairId)
	if err != nil {
		return errors.Wrap(err, "failed to get HyperMetro Pair")
	}

	err = c.SuspendHyperMetroPair(ctx, hmp.ID)
	if err != nil {
		return errors.Wrap(err, "failed to suspend HyperMetroPair")
	}

	// 2: Expand LUN
	err = c.LocalDevice.ExpandLUN(ctx, hmp.LOCALOBJID, newVolumeSizeGb)
	if err != nil {
		return errors.Wrap(err, "failed to expand Local LUN")
	}
	err = c.RemoteDevice.ExpandLUN(ctx, hmp.REMOTEOBJID, newVolumeSizeGb)
	if err != nil {
		return errors.Wrap(err, "failed to expand Remote LUN")
	}

	// 3: Re-sync HyperMetro Pair
	err = c.SyncHyperMetroPair(ctx, hmp.ID)
	if err != nil {
		return errors.Wrap(err, "failed to re-sync HyperMetro Pair")
	}

	return nil
}

func (c *Client) AttachVolume(ctx context.Context, hyperMetroPairId, hostname, iqn string) error {
	volume, err := c.GetHyperMetroPair(ctx, hyperMetroPairId)
	if err != nil {
		return errors.Wrap(err, "failed to get volume information")
	}

	err = c.LocalDevice.AttachVolume(ctx, c.PortGroupName, hostname, iqn, volume.LOCALOBJID)
	if err != nil {
		return errors.Wrap(err, "failed to attach volume in Local Device")
	}
	err = c.RemoteDevice.AttachVolume(ctx, c.PortGroupName, hostname, iqn, volume.REMOTEOBJID)
	if err != nil {
		return errors.Wrap(err, "failed to attach volume in Remote Device")
	}

	return nil
}

func (d *Device) AttachVolume(ctx context.Context, portgroupName, hostname, iqn, lunId string) error {
	// wrapper function for client.AttachVolume
	portgroups, err := d.GetPortGroups(ctx, NewSearchQueryName(portgroupName))
	if err != nil {
		return errors.Wrap(err, "failed to get portgroup")
	}
	if len(portgroups) != 1 {
		return errors.New("found multiple portgroup in same PortGroup name")
	}
	portgroup := portgroups[0]

	hostgroup, host, err := d.GetHostGroupForce(ctx, hostname)
	if err != nil {
		return errors.Wrap(err, "failed to get hostgroup")
	}
	_, err = d.GetInitiatorForce(ctx, iqn)
	if err != nil {
		return errors.Wrap(err, "failed to get initiator")
	}
	initiatorUpdateParam := UpdateInitiatorParam{
		ID:         iqn,
		TYPE:       strconv.Itoa(TypeInitiator),
		USECHAP:    "false",
		PARENTID:   host.ID,
		PARENTTYPE: strconv.Itoa(TypeHost),
	}
	_, err = d.UpdateInitiator(ctx, iqn, initiatorUpdateParam) // set PARENTID (= host.ID)
	if err != nil {
		return errors.Wrap(err, "failed to set parameter for initiator")
	}

	lungroup, err := d.GetLunGroupForce(ctx, hostname)
	if err != nil {
		return errors.Wrap(err, "failed to get lungroup")
	}

	err = d.AssociateLun(ctx, lungroup.ID, lunId)
	if err != nil {
		return errors.Wrap(err, "failed to associate lun to lungroup")
	}

	mappingview, err := d.GetMappingViewForce(ctx, hostname)
	if err != nil {
		return errors.Wrap(err, "failed to get mappingview")
	}

	err = d.DoMapping(ctx, mappingview.ID, hostgroup.ID, lungroup.ID, portgroup.ID)
	if err != nil {
		return errors.Wrap(err, "failed associate object to mappingview")
	}

	return nil
}

func (c *Client) DetachVolume(ctx context.Context, hyperMetroPairId string) error {
	volume, err := c.GetHyperMetroPair(ctx, hyperMetroPairId)
	if err != nil {
		return errors.Wrap(err, "failed to get hypermetro pair")
	}

	err = c.LocalDevice.DetachVolume(ctx, volume.LOCALOBJID)
	if err != nil {
		return errors.Wrap(err, "failed to detach volume in Local Device")
	}
	err = c.RemoteDevice.DetachVolume(ctx, volume.REMOTEOBJID)
	if err != nil {
		return errors.Wrap(err, "failed to detach volume in Remote Device")
	}

	return nil
}

func (d *Device) DetachVolume(ctx context.Context, lunId string) error {
	lun, err := d.GetLUN(ctx, lunId)
	if err != nil {
		return errors.Wrap(err, "failed to get LUN")
	}

	lungroup, err := d.GetLunGroupByLunId(ctx, lun.ID)
	if err != nil {
		return errors.Wrap(err, "failed to get lungroup")
	}

	err = d.DisAssociateLun(ctx, lungroup.ID, lun.ID)
	if err != nil {
		return errors.Wrap(err, "failed to disassociate lun")
	}

	return nil
}
