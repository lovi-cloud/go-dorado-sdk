package dorado

import (
	"errors"
	"net/http"
	"strconv"
)

// this const value drip by https://github.com/Huawei/OpenStack_Driver/blob/master/Cinder/Queens/constants.py
const (
	CapacityUnit    = 1024 * 1024 * 2 // 2 is hypermetro capacity NOTE(whywaita): honnmani?
	MaxNameLength   = 31
	DefaultDeviceID = "xx"
)

// Object Type Numbers
const (
	TypeHost             = 21
	TypeHostGroup        = 14
	TypeLUN              = 11
	TypeLUNGroup         = 256
	TypeLUNCopy          = 219
	TypeSnapshot         = 27
	TypePortGroup        = 257
	TypeInitiator        = 222
	TypeMappingView      = 245
	TypeEthernetPort     = 213
	TypeHyperMetroPair   = 15361
	TypeHyperMetroDomain = 15362
)

// For HyperMetroPair RUNNINGSTATUS
const (
	StatusNormal           = 1
	StatusSynchronizing    = 23
	StatusInvalid          = 35
	StatusPause            = 41
	StatusForcedStart      = 93
	StatusToBeSynchronized = 100
)

// For HEALTHSTATUS status
const (
	StatusHealth = 1
)

// For a some RUNNNINGSTATUS
const (
	StatusVolumeReady      = 27
	StatusLunCopyReady     = 40
	StatusSnapshotActive   = 43
	StatusSnapshotInactive = 45
)

// Dorado return Error Codes
const (
	ErrorCodeUnAuthorized  = -401
	ErrorCodeUserIsOffline = 1077949069
)

// Error Values
var (
	ErrEthernetPortNotFound     = errors.New("ethernet port is not found")
	ErrHostNotFound             = errors.New("host is not found")
	ErrHostGroupNotFound        = errors.New("host group is not found")
	ErrHyperMetroDomainNotFound = errors.New("HyperMetroDomain ID is not found")
	ErrHyperMetroPairNotFound   = errors.New("HyperMetroPair is not found")
	ErrInitiatorNotFound        = errors.New("initiator is not found")
	ErrLunNotFound              = errors.New("LUN is not found")
	ErrLunGroupNotFound         = errors.New("LUN Group is not found")
	ErrLunCopyNotFound          = errors.New("LUN Copy is not found")
	ErrMappingViewNotFound      = errors.New("mapping view is not found")
	ErrPortGroupNotFound        = errors.New("port group is not found")
	ErrSnapshotNotFound         = errors.New("snapshot is not found")
	ErrStoragePoolNotFound      = errors.New("storage pool is not found")
	ErrTargetPortNotFound       = errors.New("target port is not found")

	ErrUnAuthorized = errors.New("failed to authorized token")
	ErrTimeoutWait  = errors.New("timeout to wait")

	// parent Error
	ErrCreateRequest    = "failed to create request"
	ErrHTTPRequestDo    = "failed to HTTP request"
	ErrDecodeBody       = "failed to decodeBody"
	ErrCreatePostValue  = "failed to create post value"
	ErrRequestWithRetry = "failed to request with retry"
)

// Default values
var (
	DefaultCopyTimeoutSecond = 180
	DefaultHTTPRetryCount    = 10
)

// AssociateParam is parameter of associate functions
type AssociateParam struct {
	ID               string `json:"ID,omitempty"`
	TYPE             string `json:"TYPE,omitempty"`
	ASSOCIATEOBJID   string `json:"ASSOCIATEOBJID,omitempty"`
	ASSOCIATEOBJTYPE int    `json:"ASSOCIATEOBJTYPE,omitempty"`
}

// AddAssociateParam add AssociateParam to http.Request
func AddAssociateParam(req *http.Request, param *AssociateParam) *http.Request {
	if param == nil {
		return req
	}

	q := req.URL.Query()

	if param.ID != "" {
		q.Add("ID", param.ID)
	}
	if param.TYPE != "" {
		q.Add("TYPE", param.TYPE)
	}
	if param.ASSOCIATEOBJID != "" {
		q.Add("ASSOCIATEOBJID", param.ASSOCIATEOBJID)
	}
	if param.ASSOCIATEOBJTYPE != 0 {
		q.Add("ASSOCIATEOBJTYPE", strconv.Itoa(param.ASSOCIATEOBJTYPE))
	}

	req.URL.RawQuery = q.Encode()

	return req
}
