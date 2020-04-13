package dorado

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	uuid "github.com/satori/go.uuid"
)

type LUN struct {
	ALLOCCAPACITY               string `json:"ALLOCCAPACITY"`
	ALLOCTYPE                   string `json:"ALLOCTYPE"`
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

func encodeLunName(u uuid.UUID) string {
	// MAX LUN name length is 31, but uuid
	// this function binding by huawei_utils.encode_name(id) in OpenStack cinder-driver.
	values := strings.Split(u.String(), "-")
	prefix := "w-" + values[0] + "-" // TODO(whywaita): delete w- later.

	hash := md5.Sum(u.Bytes())
	name := prefix + hex.EncodeToString(hash[:])[:MAX_NAME_LENGTH-len(prefix)]
	return name
}

func (d *Device) GetLUNs(ctx context.Context, query *SearchQuery) ([]LUN, error) {
	spath := "/lun"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}

	req = AddSearchQuery(req, query)

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	luns := []LUN{}
	if err = decodeBody(resp, &luns); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
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
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	lun := &LUN{}
	if err = decodeBody(resp, lun); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return lun, nil
}

func (d *Device) CreateLUN(ctx context.Context, u uuid.UUID, capacityGB int, storagePoolId string) (*LUN, error) {
	// low level API
	spath := "/lun"

	p := ParamCreateLUN{
		NAME:               encodeLunName(u),
		PARENTID:           storagePoolId,
		DESCRIPTION:        "volume-" + u.String(),
		CAPACITY:           capacityGB * CAPACITY_UNIT,
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
		return nil, errors.Wrap(err, ErrCreatePostValue)
	}
	req, err := d.newRequest(ctx, "POST", spath, bytes.NewBuffer(jb))
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	lun := &LUN{}
	if err := decodeBody(resp, lun); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	return lun, nil
}

func (d *Device) DeleteLUN(ctx context.Context, id string) error {
	spath := fmt.Sprintf("/lun/%s", id)
	req, err := d.newRequest(ctx, "DELETE", spath, nil)
	if err != nil {
		return errors.Wrap(err, ErrCreateRequest)
	}

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, ErrHTTPRequestDo)
	}

	var i interface{} // this endpoint return N/A
	if err := decodeBody(resp, i); err != nil {
		return errors.Wrap(err, ErrDecodeBody)
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
		CAPACITY: uint64(newLunSizeGb * CAPACITY_UNIT),
	}
	jb, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err, ErrCreatePostValue)
	}

	req, err := d.newRequest(ctx, "PUT", spath, bytes.NewBuffer(jb))
	if err != nil {
		return errors.Wrap(err, ErrCreateRequest)
	}
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, ErrHTTPRequestDo)
	}

	var i interface{} // this endpoint return N/A
	if err = decodeBody(resp, i); err != nil {
		return errors.Wrap(err, ErrDecodeBody)
	}

	return nil
}

func (d *Device) GetAssociateLUNs(ctx context.Context, query *SearchQuery) ([]LUN, error) {
	spath := "/lun/associate"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest)
	}
	req = AddSearchQuery(req, query)
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPRequestDo)
	}

	luns := []LUN{}
	if err = decodeBody(resp, luns); err != nil {
		return nil, errors.Wrap(err, ErrDecodeBody)
	}

	if len(luns) == 0 {
		return nil, errors.New(ErrLunNotFound)
	}

	return luns, nil
}
