package dorado

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	uuid "github.com/satori/go.uuid"
)

// LUN is raw block storage object.
type LUN struct {
	ALLOCCAPACITY               string `json:"ALLOCCAPACITY"`
	ALLOCTYPE                   string `json:"ALLOCTYPE"`
	ASSOCIATEMETADATA           string `json:"ASSOCIATEMETADATA"`
	CAPABILITY                  string `json:"CAPABILITY"`
	CAPACITY                    int    `json:"CAPACITY,string"`
	CAPACITYALARMLEVEL          string `json:"CAPACITYALARMLEVEL"`
	CLONEIDS                    string `json:"CLONEIDS"`
	COMPRESSION                 string `json:"COMPRESSION"`
	COMPRESSIONSAVEDCAPACITY    string `json:"COMPRESSIONSAVEDCAPACITY"`
	COMPRESSIONSAVEDRATIO       string `json:"COMPRESSIONSAVEDRATIO"`
	DEDUPSAVEDCAPACITY          string `json:"DEDUPSAVEDCAPACITY"`
	DEDUPSAVEDRATIO             string `json:"DEDUPSAVEDRATIO"`
	DESCRIPTION                 string `json:"DESCRIPTION"`
	DRSENABLE                   string `json:"DRS_ENABLE"`
	ENABLECOMPRESSION           string `json:"ENABLECOMPRESSION"`
	ENABLEISCSITHINLUNTHRESHOLD string `json:"ENABLEISCSITHINLUNTHRESHOLD"`
	ENABLESMARTDEDUP            string `json:"ENABLESMARTDEDUP"`
	EXPOSEDTOINITIATOR          string `json:"EXPOSEDTOINITIATOR"`
	EXTENDIFSWITCH              string `json:"EXTENDIFSWITCH"`
	HEALTHSTATUS                string `json:"HEALTHSTATUS"`
	HYPERCDPSCHEDULEDISABLE     string `json:"HYPERCDPSCHEDULEDISABLE"`
	ID                          int    `json:"ID,string"`
	IOCLASSID                   string `json:"IOCLASSID"`
	IOPRIORITY                  string `json:"IOPRIORITY"`
	ISADD2LUNGROUP              bool   `json:"ISADD2LUNGROUP,string"`
	ISCHECKZEROPAGE             string `json:"ISCHECKZEROPAGE"`
	ISCLONE                     bool   `json:"ISCLONE,string"`
	ISCSITHINLUNTHRESHOLD       string `json:"ISCSITHINLUNTHRESHOLD"`
	LUNCOPYIDS                  string `json:"LUNCOPYIDS"`
	LUNMigrationOrigin          string `json:"LUNMigrationOrigin"`
	MIRRORPOLICY                string `json:"MIRRORPOLICY"`
	MIRRORTYPE                  string `json:"MIRRORTYPE"`
	NAME                        string `json:"NAME"`
	OWNINGCONTROLLER            string `json:"OWNINGCONTROLLER"`
	PARENTID                    int    `json:"PARENTID,string"`
	PARENTNAME                  string `json:"PARENTNAME"`
	PREFETCHPOLICY              string `json:"PREFETCHPOLICY"`
	PREFETCHVALUE               string `json:"PREFETCHVALUE"`
	REMOTELUNID                 string `json:"REMOTELUNID"`
	REMOTEREPLICATIONIDS        string `json:"REMOTEREPLICATIONIDS"`
	REPLICATIONCAPACITY         string `json:"REPLICATION_CAPACITY"`
	RUNNINGSTATUS               string `json:"RUNNINGSTATUS"`
	RUNNINGWRITEPOLICY          string `json:"RUNNINGWRITEPOLICY"`
	SECTORSIZE                  string `json:"SECTORSIZE"`
	SNAPSHOTIDS                 string `json:"SNAPSHOTIDS"`
	SNAPSHOTSCHEDULEID          string `json:"SNAPSHOTSCHEDULEID"`
	SUBTYPE                     string `json:"SUBTYPE"`
	THINCAPACITYUSAGE           string `json:"THINCAPACITYUSAGE"`
	TOTALSAVEDCAPACITY          string `json:"TOTALSAVEDCAPACITY"`
	TOTALSAVEDRATIO             string `json:"TOTALSAVEDRATIO"`
	TYPE                        int    `json:"TYPE"`
	USAGETYPE                   string `json:"USAGETYPE"`
	WORKINGCONTROLLER           string `json:"WORKINGCONTROLLER"`
	WORKLOADTYPEID              string `json:"WORKLOADTYPEID"`
	WORKLOADTYPENAME            string `json:"WORKLOADTYPENAME"`
	WRITEPOLICY                 string `json:"WRITEPOLICY"`
	WWN                         string `json:"WWN"`
	HyperCdpScheduleID          string `json:"hyperCdpScheduleId"`
	LunCgID                     string `json:"lunCgId"`
	RemoteLunWwn                string `json:"remoteLunWwn"`
	TakeOverLunWwn              string `json:"takeOverLunWwn"`
}

// AssociateMetaData is one of LUN parameter
type AssociateMetaData struct {
	HostLUNID int `json:"HostLUNID"`
}

// ParamCreateLUN is parameter for CreateLUN
type ParamCreateLUN struct {
	WRITEPOLICY        string `json:"WRITEPOLICY"`
	PREFETCHVALUE      string `json:"PREFETCHVALUE"`
	ALLOCTYPE          int    `json:"ALLOCTYPE"`
	PARENTID           string `json:"PARENTID"`
	MIRRORPOLICY       string `json:"MIRRORPOLICY"`
	DATATRANSFERPOLICY string `json:"DATATRANSFERPOLICY"`
	DESCRIPTION        string `json:"DESCRIPTION"`
	CAPACITY           int    `json:"CAPACITY"`
	NAME               string `json:"NAME"`
	WORKLOADTYPEID     string `json:"WORKLOADTYPEID"`
	PREFETCHPOLICY     string `json:"PREFETCHPOLICY"`
}

// ParamCreateCloneLUN is parameter for CreateCloneLUN
type ParamCreateCloneLUN struct {
	NAME          string `json:"NAME"`
	CLONESOURCEID int    `json:"CLONESOURCEID"`
	ISCLONE       bool   `json:"ISCLONE"`
}

// PrefixVolumeDescription is prefix of volume Description
var PrefixVolumeDescription = "volume-"

// EncodeLunName encode name for LUN Name
func EncodeLunName(u uuid.UUID) string {
	// MAX LUN name length is 31, but uuid
	// this function binding by huawei_utils.encode_name(id) in OpenStack cinder-driver.
	// ref: https://github.com/openstack/cinder/blob/006a4f48174c04c8720175f69271a317906867e9/cinder/volume/drivers/huawei/huawei_utils.py#L38
	values := strings.Split(u.String(), "-")
	prefix := values[0] + "-"

	hash := md5.Sum(u.Bytes())
	name := prefix + hex.EncodeToString(hash[:])[:MaxNameLength-len(prefix)]
	return name
}

// GetLUNs get lun objects by query
func (d *Device) GetLUNs(ctx context.Context, query *SearchQuery) ([]LUN, error) {
	spath := "/lun"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req = AddSearchQuery(req, query)

	var luns []LUN
	if err = d.requestWithRetry(req, &luns, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	if len(luns) == 0 {
		return nil, ErrLunNotFound
	}

	return luns, nil
}

// GetLUN get lun object by id
func (d *Device) GetLUN(ctx context.Context, lunID int) (*LUN, error) {
	spath := fmt.Sprintf("/lun/%d", lunID)

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	lun := &LUN{}
	if err = d.requestWithRetry(req, lun, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return lun, nil
}

// CreateLUN create lun object
func (d *Device) CreateLUN(ctx context.Context, u uuid.UUID, capacityGB int, storagePoolName string) (*LUN, error) {
	storagePools, err := d.GetStoragePools(ctx, NewSearchQueryName(storagePoolName))
	if err != nil {
		return nil, fmt.Errorf("failed to get storagepool: %w", err)
	}

	if len(storagePools) != 1 {
		return nil, errors.New("found multiple storagepool in same name")
	}
	storagePoolID := storagePools[0].ID

	p := ParamCreateLUN{
		NAME:               EncodeLunName(u),
		PARENTID:           strconv.Itoa(storagePoolID),
		DESCRIPTION:        PrefixVolumeDescription + u.String(),
		CAPACITY:           capacityGB * CapacityUnit,
		WRITEPOLICY:        "1",
		PREFETCHVALUE:      "0",
		ALLOCTYPE:          1,
		MIRRORPOLICY:       "1",
		DATATRANSFERPOLICY: "0",
		WORKLOADTYPEID:     "0",
		PREFETCHPOLICY:     "3",
	}

	return d.createLUN(ctx, p)
}

func (d *Device) createLUN(ctx context.Context, param interface{}) (*LUN, error) {
	spath := "/lun"

	jb, err := json.Marshal(param)
	if err != nil {
		return nil, fmt.Errorf(ErrCreatePostValue+": %w", err)
	}
	req, err := d.newRequest(ctx, "POST", spath, bytes.NewBuffer(jb))
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	lun := &LUN{}
	if err = d.requestWithRetry(req, lun, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return lun, nil
}

// CreateLUNWithWait create LUN and waiting ready
func (d *Device) CreateLUNWithWait(ctx context.Context, u uuid.UUID, capacityGB int, storagePoolName string) (*LUN, error) {
	lun, err := d.CreateLUN(ctx, u, capacityGB, storagePoolName)
	if err != nil {
		return nil, fmt.Errorf("failed to create LUN: %w", err)
	}

	// wait 10 seconds
	for i := 0; i < 10; i++ {
		isReady, err := d.lunIsReady(ctx, lun.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to wait that LUN is ready: %w", err)
		}

		if isReady == true {
			return d.GetLUN(ctx, lun.ID)
		}

		time.Sleep(1 * time.Second)
	}

	return nil, ErrTimeoutWait
}

func (d *Device) lunIsReady(ctx context.Context, LUNID int) (bool, error) {
	lun, err := d.GetLUN(ctx, LUNID)
	if err != nil {
		return false, fmt.Errorf("failed to get LUN (ID: %d): %w", LUNID, err)
	}

	if lun.HEALTHSTATUS == strconv.Itoa(StatusHealth) &&
		lun.RUNNINGSTATUS == strconv.Itoa(StatusVolumeReady) &&
		lun.ISCLONE == false {
		return true, nil
	}

	return false, nil
}

// DeleteLUN delete lun object (also include data)
func (d *Device) DeleteLUN(ctx context.Context, lunID int) error {
	spath := fmt.Sprintf("/lun/%d", lunID)
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

// ExpandLUN expand lun capacity
func (d *Device) ExpandLUN(ctx context.Context, lunID int, newLunSizeGb int) error {
	spath := "/lun/expand"
	param := struct {
		ID       string `json:"ID"`
		TYPE     int    `json:"TYPE"`
		CAPACITY uint64 `json:"CAPACITY"`
	}{
		ID:       strconv.Itoa(lunID),
		TYPE:     TypeLUN,
		CAPACITY: uint64(newLunSizeGb * CapacityUnit),
	}
	jb, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf(ErrCreatePostValue+": %w", err)
	}

	req, err := d.newRequest(ctx, "PUT", spath, bytes.NewBuffer(jb))
	if err != nil {
		return fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	var i interface{} // this endpoint return N/A
	if err = d.requestWithRetry(req, i, DefaultHTTPRetryCount); err != nil {
		return fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return nil
}

// GetAssociateLUNs get lun objects that associated object (ex: host)
func (d *Device) GetAssociateLUNs(ctx context.Context, query *SearchQuery) ([]LUN, error) {
	spath := "/lun/associate"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req = AddSearchQuery(req, query)

	var luns []LUN
	if err = d.requestWithRetry(req, &luns, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	if len(luns) == 0 {
		return nil, ErrLunNotFound
	}

	return luns, nil
}

// GetHostAssociatedLUNs get LUNs associated specific host
func (d *Device) GetHostAssociatedLUNs(ctx context.Context, hostID int) ([]LUN, error) {
	query := &SearchQuery{
		AssociateObjType: strconv.Itoa(TypeHost),
		AssociateObjID:   strconv.Itoa(hostID),
	}

	return d.GetAssociateLUNs(ctx, query)
}

// GetHostLUNID get LUN ID per host.
func (d *Device) GetHostLUNID(ctx context.Context, lunID, hostID int) (int, error) {
	luns, err := d.GetHostAssociatedLUNs(ctx, hostID)
	if err != nil {
		return 0, fmt.Errorf("failed to get associated LUNs: %w", err)
	}

	for _, lun := range luns {
		if lun.ID == lunID {
			jsonStr := lun.ASSOCIATEMETADATA
			hostLunID := AssociateMetaData{}
			err := json.Unmarshal([]byte(jsonStr), &hostLunID)
			if err != nil {
				return 0, fmt.Errorf("failed to parse ASSOCIATEMETADATA: %w", err)
			}

			return hostLunID.HostLUNID, nil
		}
	}

	return 0, fmt.Errorf("LUN (ID: %d) is not associated host (ID: %d)", lunID, hostID)
}

// CreateCloneLUN create clone LUN
func (d *Device) CreateCloneLUN(ctx context.Context, lunID int, lunName uuid.UUID) (*LUN, error) {
	param := ParamCreateCloneLUN{
		CLONESOURCEID: lunID,
		ISCLONE:       true,
		NAME:          EncodeLunName(lunName),
	}

	lun, err := d.createLUN(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("failed to create LUN: %w", err)
	}

	return lun, nil
}

// SplitCloneLUN start to split LUN Clone
func (d *Device) SplitCloneLUN(ctx context.Context, cloneLUNID int) error {
	spath := "/lunclone_split_switch"
	param := struct {
		ID          int  `json:"ID"`
		SPLITACTION int  `json:"SPLITACTION"`
		ISCLONE     bool `json:"ISCLONE"`
		SPLITSPEED  int  `json:"SPLITSPEED"`
	}{
		ID:          cloneLUNID,
		SPLITACTION: 1,
		ISCLONE:     true,
		SPLITSPEED:  4,
	}
	jb, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf(ErrCreatePostValue+": %w", err)
	}

	req, err := d.newRequest(ctx, "PUT", spath, bytes.NewBuffer(jb))
	if err != nil {
		return fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	var i interface{} // this endpoint return N/A
	if err = d.requestWithRetry(req, i, DefaultHTTPRetryCount); err != nil {
		return fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return nil
}
