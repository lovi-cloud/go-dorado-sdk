package dorado

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	uuid "github.com/satori/go.uuid"
)

func (c *Client) CreateVolumeRaw(ctx context.Context, name uuid.UUID, capacityGB int, storagePoolName, hyperMetroDomainID string) (*HyperMetroPair, error) {
	// create volume (= hypermetro enabled lun)
	localLun, err := c.LocalDevice.CreateLUN(ctx, name, capacityGB, storagePoolName)
	if err != nil {
		return nil, fmt.Errorf("failed to create lun in local device: %w", err)
	}

	remoteLun, err := c.RemoteDevice.CreateLUN(ctx, name, capacityGB, storagePoolName)
	if err != nil {
		return nil, fmt.Errorf("failed to create lun in remote device: %w", err)
	}

	hyperMetroPair, err := c.CreateHyperMetroPair(ctx, hyperMetroDomainID, localLun.ID, remoteLun.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create HyperMetroPair: %w", err)
	}

	return hyperMetroPair, nil
}

func (c *Client) CreateVolumeFromSource(ctx context.Context, name uuid.UUID, capacityGB int, storagePoolName, hyperMetroDomainID string, sourceHyperMetroPairID string) (*HyperMetroPair, error) {
	source, err := c.GetHyperMetroPair(ctx, sourceHyperMetroPairID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source HyperMetroPair: %w", err)
	}

	localLun, err := c.LocalDevice.CreateLUNFromSource(ctx, source.LOCALOBJID, name, capacityGB, storagePoolName)
	if err != nil {
		return nil, fmt.Errorf("failed to crteate lun from source in local device: %w", err)
	}

	remoteLun, err := c.RemoteDevice.CreateLUNFromSource(ctx, source.REMOTEOBJID, name, capacityGB, storagePoolName)
	if err != nil {
		return nil, fmt.Errorf("failed to crteate lun from source in remote device: %w", err)
	}

	hyperMetroPair, err := c.CreateHyperMetroPair(ctx, hyperMetroDomainID, localLun.ID, remoteLun.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create HyperMetroPair from source: %w", err)
	}

	return hyperMetroPair, nil
}

func (d *Device) CreateLUNFromSource(ctx context.Context, sourceLUNID int, name uuid.UUID, capacityGB int, storagePoolName string) (*LUN, error) {
	snapshot, err := d.CreateSnapshotWithWait(ctx, sourceLUNID, name, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}
	defer func() {
		d.StopSnapshot(ctx, snapshot.ID)
		d.DeleteSnapshot(ctx, snapshot.ID)
	}()

	targetLUN, err := d.CreateLUNWithWait(ctx, name, capacityGB, storagePoolName)
	if err != nil {
		return nil, fmt.Errorf("failed to create raw LUN: %w", err)
	}

	luncopy, err := d.CreateLUNCopy(ctx, snapshot.ID, targetLUN.ID)
	if err != nil {
		d.DeleteLUN(ctx, targetLUN.ID)
		return nil, fmt.Errorf("failed to create luncopy object: %w", err)
	}
	defer d.DeleteLUNCopy(ctx, luncopy.ID)

	err = d.StartLUNCopyWithWait(ctx, luncopy.ID, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to copy lun: %w", err)
	}

	return d.GetLUN(ctx, targetLUN.ID)
}

func (c *Client) DeleteVolume(ctx context.Context, hyperMetroPairID string) error {
	// 1: delete HyperMetro Pair
	// 2: delete LUN Group Associate
	// 3: delete LUN

	hmp, err := c.GetHyperMetroPair(ctx, hyperMetroPairID)
	if err != nil {
		return fmt.Errorf("failed to get HyperMetro Pair: %w", err)
	}

	// 2: delete LUN Group Associate
	llun, err := c.LocalDevice.GetLUN(ctx, hmp.LOCALOBJID)
	if err != nil {
		return fmt.Errorf("failed to get lun information: %w", err)
	}
	if llun.ISADD2LUNGROUP == true {
		lLungroup, err := c.LocalDevice.GetLunGroupByLunId(ctx, hmp.LOCALOBJID)
		if err != nil {
			return fmt.Errorf("failed to get lungroup by associated lun: %w", err)
		}
		err = c.LocalDevice.DisAssociateLun(ctx, lLungroup.ID, hmp.LOCALOBJID)
		if err != nil {
			return fmt.Errorf("failed to disassociate local lun: %w", err)
		}
	}

	rlun, err := c.RemoteDevice.GetLUN(ctx, hmp.REMOTEOBJID)
	if err != nil {
		return fmt.Errorf("failed to get lun information: %w", err)
	}
	if rlun.ISADD2LUNGROUP == true {
		rLungroup, err := c.RemoteDevice.GetLunGroupByLunId(ctx, hmp.REMOTEOBJID)
		if err != nil {
			return fmt.Errorf("failed to get lungroup by associated lun: %w", err)
		}
		err = c.RemoteDevice.DisAssociateLun(ctx, rLungroup.ID, hmp.REMOTEOBJID)
		if err != nil {
			return fmt.Errorf("failed to disassociate remote lun: %w", err)
		}
	}

	// 1: delete HyperMetro Pair
	if hmp.RUNNINGSTATUS != strconv.Itoa(StatusPause) {
		err = c.SuspendHyperMetroPair(ctx, hmp.ID)
		if err != nil {
			return fmt.Errorf("failed to suspend HyperMetroPair: %w", err)
		}
	}
	err = c.DeleteHyperMetroPair(ctx, hmp.ID)
	if err != nil {
		return fmt.Errorf("failed to delete HyperMetroPair: %w", err)
	}

	// 3: delete LUN
	err = c.LocalDevice.DeleteLUN(ctx, hmp.LOCALOBJID)
	if err != nil {
		return fmt.Errorf("failed to delete Local LUN: %w", err)
	}
	err = c.RemoteDevice.DeleteLUN(ctx, hmp.REMOTEOBJID)
	if err != nil {
		return fmt.Errorf("failed to delete Remote LUN: %w", err)
	}

	return nil
}

func (c *Client) ExtendVolume(ctx context.Context, hyperMetroPairID string, newVolumeSizeGb int) error {
	// 1: Suspend HyperMetro Pair
	// 2: Expand LUN
	// 3: Re-sync HyperMetro Pair

	// 1: Suspend HyperMetro Pair
	hmp, err := c.GetHyperMetroPair(ctx, hyperMetroPairID)
	if err != nil {
		return fmt.Errorf("failed to get HyperMetro Pair: %w", err)
	}

	err = c.SuspendHyperMetroPair(ctx, hmp.ID)
	if err != nil {
		return fmt.Errorf("failed to suspend HyperMetroPair: %w", err)
	}

	// 2: Expand LUN
	err = c.LocalDevice.ExpandLUN(ctx, hmp.LOCALOBJID, newVolumeSizeGb)
	if err != nil {
		return fmt.Errorf("failed to expand Local LUN: %w", err)
	}
	err = c.RemoteDevice.ExpandLUN(ctx, hmp.REMOTEOBJID, newVolumeSizeGb)
	if err != nil {
		return fmt.Errorf("failed to expand Remote LUN: %w", err)
	}

	// 3: Re-sync HyperMetro Pair
	err = c.SyncHyperMetroPair(ctx, hmp.ID)
	if err != nil {
		return fmt.Errorf("failed to re-sync HyperMetro Pair: %w", err)
	}

	return nil
}

func (c *Client) AttachVolume(ctx context.Context, hyperMetroPairID, hostname, iqn string) error {
	volume, err := c.GetHyperMetroPair(ctx, hyperMetroPairID)
	if err != nil {
		return fmt.Errorf("failed to get volume information: %w", err)
	}

	err = c.LocalDevice.AttachVolume(ctx, c.PortGroupName, hostname, iqn, volume.LOCALOBJID)
	if err != nil {
		return fmt.Errorf("failed to attach volume in Local Device: %w", err)
	}
	err = c.RemoteDevice.AttachVolume(ctx, c.PortGroupName, hostname, iqn, volume.REMOTEOBJID)
	if err != nil {
		return fmt.Errorf("failed to attach volume in Remote Device: %w", err)
	}

	return nil
}

func (d *Device) AttachVolume(ctx context.Context, portgroupName, hostname, iqn string, lunID int) error {
	// wrapper function for client.AttachVolume
	portgroups, err := d.GetPortGroups(ctx, NewSearchQueryName(portgroupName))
	if err != nil {
		return fmt.Errorf("failed to get portgroup: %w", err)
	}
	if len(portgroups) != 1 {
		return errors.New("found multiple portgroup in same PortGroup name")
	}
	portgroup := portgroups[0]

	hostgroup, host, err := d.GetHostGroupForce(ctx, hostname)
	if err != nil {
		return fmt.Errorf("failed to get hostgroup: %w", err)
	}
	_, err = d.GetInitiatorForce(ctx, iqn)
	if err != nil {
		return fmt.Errorf("failed to get initiator: %w", err)
	}
	initiatorUpdateParam := UpdateInitiatorParam{
		ID:         iqn,
		TYPE:       strconv.Itoa(TypeInitiator),
		USECHAP:    "false",
		PARENTID:   strconv.Itoa(host.ID),
		PARENTTYPE: strconv.Itoa(TypeHost),
	}
	_, err = d.UpdateInitiator(ctx, iqn, initiatorUpdateParam) // set PARENTID (= host.ID)
	if err != nil {
		return fmt.Errorf("failed to set parameter for initiator: %w", err)
	}

	lungroup, err := d.GetLunGroupForce(ctx, hostname)
	if err != nil {
		return fmt.Errorf("failed to get lungroup: %w", err)
	}

	err = d.AssociateLun(ctx, lungroup.ID, lunID)
	if err != nil {
		return fmt.Errorf("failed to associate lun to lungroup: %w", err)
	}

	mappingview, err := d.GetMappingViewForce(ctx, hostname)
	if err != nil {
		return fmt.Errorf("failed to get mappingview: %w", err)
	}

	err = d.DoMapping(ctx, mappingview, hostgroup, lungroup, portgroup.ID)
	if err != nil {
		return fmt.Errorf("failed to associate object to mappingview: %w", err)
	}

	return nil
}

func (c *Client) DetachVolume(ctx context.Context, hyperMetroPairId string) error {
	volume, err := c.GetHyperMetroPair(ctx, hyperMetroPairId)
	if err != nil {
		return fmt.Errorf("failed to get hypermetro pair: %w", err)
	}

	err = c.LocalDevice.DetachVolume(ctx, volume.LOCALOBJID)
	if err != nil {
		return fmt.Errorf("failed to detach volume in Local Device: %w", err)
	}
	err = c.RemoteDevice.DetachVolume(ctx, volume.REMOTEOBJID)
	if err != nil {
		return fmt.Errorf("failed to detach volume in Remote Device: %w", err)
	}

	return nil
}

func (d *Device) DetachVolume(ctx context.Context, lunID int) error {
	lun, err := d.GetLUN(ctx, lunID)
	if err != nil {
		return fmt.Errorf("failed to get LUN: %w", err)
	}

	lungroup, err := d.GetLunGroupByLunId(ctx, lun.ID)
	if err != nil {
		return fmt.Errorf("failed to get lungroup: %w", err)
	}

	err = d.DisAssociateLun(ctx, lungroup.ID, lun.ID)
	if err != nil {
		return fmt.Errorf("failed to disassociate lun: %w", err)
	}

	// TODO: delete host, hostgroup and lungroup if nothing associate object

	return nil
}
