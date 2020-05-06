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

	"github.com/pkg/errors"

	uuid "github.com/satori/go.uuid"
)

type LUN struct {
	ALLOCCAPACITY               string `json:"ALLOCCAPACITY"`
	ALLOCTYPE                   string `json:"ALLOCTYPE"`
	ASSOCIATEMETADATA           string `json:"ASSOCIATEMETADATA"`
	CAPABILITY                  string `json:"CAPABILITY"`
	CAPACITY                    string `json:"CAPACITY"`
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
	ID                          string `json:"ID"`
	IOCLASSID                   string `json:"IOCLASSID"`
	IOPRIORITY                  string `json:"IOPRIORITY"`
	ISADD2LUNGROUP              string `json:"ISADD2LUNGROUP"`
	ISCHECKZEROPAGE             string `json:"ISCHECKZEROPAGE"`
	ISCLONE                     string `json:"ISCLONE"`
	ISCSITHINLUNTHRESHOLD       string `json:"ISCSITHINLUNTHRESHOLD"`
	LUNCOPYIDS                  string `json:"LUNCOPYIDS"`
	LUNMigrationOrigin          string `json:"LUNMigrationOrigin"`
	MIRRORPOLICY                string `json:"MIRRORPOLICY"`
	MIRRORTYPE                  string `json:"MIRRORTYPE"`
	NAME                        string `json:"NAME"`
	OWNINGCONTROLLER            string `json:"OWNINGCONTROLLER"`
	PARENTID                    string `json:"PARENTID"`
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

type ASSOCIATEMETADATA struct {
	HostLUNID int `json:"HostLUNID"`
}

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

const (
	ErrLunNotFound = "LUN is not found"
)

func EncodeLunName(u uuid.UUID) string {
	// MAX LUN name length is 31, but uuid
	// this function binding by huawei_utils.encode_name(id) in OpenStack cinder-driver.
	values := strings.Split(u.String(), "-")
	prefix := "w-" + values[0] + "-" // TODO(whywaita): delete w- later.

	hash := md5.Sum(u.Bytes())
	name := prefix + hex.EncodeToString(hash[:])[:MaxNameLength-len(prefix)]
	return name
}

func (d *Device) GetLUNs(ctx context.Context, query *SearchQuery) ([]LUN, error) {
	spath := "/lun"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	req = AddSearchQuery(req, query)

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	luns := []LUN{}
	if err = decodeBody(resp, &luns); err != nil {
		return nil, fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	if len(luns) == 0 {
		return nil, errors.New(ErrLunNotFound)
	}

	return luns, nil
}

func (d *Device) GetLUN(ctx context.Context, lunId string) (*LUN, error) {
	spath := fmt.Sprintf("/lun/%s", lunId)

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	lun := &LUN{}
	if err = decodeBody(resp, lun); err != nil {
		return nil, fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	return lun, nil
}

func (d *Device) CreateLUN(ctx context.Context, u uuid.UUID, capacityGB int, storagePoolName string) (*LUN, error) {
	// low level API
	storagePools, err := d.GetStoragePools(ctx, NewSearchQueryName(storagePoolName))
	if err != nil {
		return nil, fmt.Errorf("failed to get storagepool: %w", err)
	}

	if len(storagePools) != 1 {
		return nil, errors.New("found multiple storagepool in same name")
	}
	storagePoolId := storagePools[0].ID

	spath := "/lun"

	p := ParamCreateLUN{
		NAME:               EncodeLunName(u),
		PARENTID:           storagePoolId,
		DESCRIPTION:        "volume-" + u.String(),
		CAPACITY:           capacityGB * CapacityUnit,
		WRITEPOLICY:        "1",
		PREFETCHVALUE:      "0",
		ALLOCTYPE:          1,
		MIRRORPOLICY:       "1",
		DATATRANSFERPOLICY: "0",
		WORKLOADTYPEID:     "0",
		PREFETCHPOLICY:     "3",
	}
	jb, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf(ErrCreatePostValue+": %w", err)
	}
	req, err := d.newRequest(ctx, "POST", spath, bytes.NewBuffer(jb))
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	lun := &LUN{}
	if err := decodeBody(resp, lun); err != nil {
		return nil, fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	return lun, nil
}

func (d *Device) DeleteLUN(ctx context.Context, id string) error {
	spath := fmt.Sprintf("/lun/%s", id)
	req, err := d.newRequest(ctx, "DELETE", spath, nil)
	if err != nil {
		return fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	var i interface{} // this endpoint return N/A
	if err := decodeBody(resp, i); err != nil {
		return fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	return nil
}

func (d *Device) ExpandLUN(ctx context.Context, id string, newLunSizeGb int) error {
	spath := "/lun/expand"
	param := struct {
		ID       string `json:"ID"`
		TYPE     int    `json:"TYPE"`
		CAPACITY uint64 `json:"CAPACITY"`
	}{
		ID:       id,
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
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	var i interface{} // this endpoint return N/A
	if err = decodeBody(resp, i); err != nil {
		return fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	return nil
}

func (d *Device) GetAssociateLUNs(ctx context.Context, query *SearchQuery) ([]LUN, error) {
	spath := "/lun/associate"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req = AddSearchQuery(req, query)
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	luns := []LUN{}
	if err = decodeBody(resp, luns); err != nil {
		return nil, fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	if len(luns) == 0 {
		return nil, errors.New(ErrLunNotFound)
	}

	return luns, nil
}

// GetHostLUNID get LUN ID per host.
func (d *Device) GetHostLUNID(ctx context.Context, lunID, hostID int) (string, error) {
	query := &SearchQuery{
		AssociateObjType: strconv.Itoa(TypeHost),
		AssociateObjID:   strconv.Itoa(hostID),
	}

	luns, err := d.GetAssociateLUNs(ctx, query)
	if err != nil {
		return "", fmt.Errorf("failed to get associated LUNs: %w", err)
	}

	strLunID := strconv.Itoa(lunID)
	for _, lun := range luns {
		if lun.ID == strLunID {
			jsonStr := lun.ASSOCIATEMETADATA
			hostLunId := ASSOCIATEMETADATA{}
			err := json.Unmarshal([]byte(jsonStr), &hostLunId)
			if err != nil {
				return "", fmt.Errorf("failed to parse ASSOCIATEMETADATA: %w", err)
			}

			return strconv.Itoa(hostLunId.HostLUNID), nil
		}
	}

	return "", fmt.Errorf("LUN (ID: %d) is not associated host (ID: %d)", lunID, hostID)
}
