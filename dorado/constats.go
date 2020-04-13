package dorado

// this const value drip by https://github.com/Huawei/OpenStack_Driver/blob/master/Cinder/Queens/constants.py

const (
	CAPACITY_UNIT   = 1024 * 1024 * 2 // 2 is hypermetro capacity NOTE(whywaita): honnmani?
	MAX_NAME_LENGTH = 31
)

const (
	TypeHost           = 21
	TypeHostGroup      = 14
	TypeLUN            = 11
	TypeLUNGroup       = 256
	TypePortGroup      = 257
	TypeInitiator      = 222
	TypeMappingView    = 245
	TypeHyperMetroPair = 15361
)

type AssociateParam struct {
	ID               string `json:"ID,omitempty"`
	TYPE             string `json:"TYPE,omitempty"`
	ASSOCIATEOBJID   string `json:"ASSOCIATEOBJID,omitempty"`
	ASSOCIATEOBJTYPE int    `json:"ASSOCIATEOBJTYPE,omitempty"`
}
