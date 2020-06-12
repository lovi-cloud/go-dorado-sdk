package dorado

import (
	"net/http"
	"strconv"
)

// this const value drip by https://github.com/Huawei/OpenStack_Driver/blob/master/Cinder/Queens/constants.py
const (
	CapacityUnit  = 1024 * 1024 * 2 // 2 is hypermetro capacity NOTE(whywaita): honnmani?
	MaxNameLength = 31
)

// Type definition
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

// For HEALTHSTATUS
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
