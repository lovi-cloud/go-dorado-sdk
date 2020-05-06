package dorado

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type EthernetPort struct {
	BONDID             string `json:"BONDID"`
	BONDNAME           string `json:"BONDNAME"`
	ERRORPACKETS       string `json:"ERRORPACKETS"`
	ETHDUPLEX          string `json:"ETHDUPLEX"`
	ETHNEGOTIATE       string `json:"ETHNEGOTIATE"`
	HEALTHSTATUS       string `json:"HEALTHSTATUS"`
	ID                 string `json:"ID"`
	INIORTGT           string `json:"INIORTGT"`
	IPV4ADDR           string `json:"IPV4ADDR"`
	IPV4GATEWAY        string `json:"IPV4GATEWAY"`
	IPV4MASK           string `json:"IPV4MASK"`
	IPV6ADDR           string `json:"IPV6ADDR"`
	IPV6GATEWAY        string `json:"IPV6GATEWAY"`
	IPV6MASK           string `json:"IPV6MASK"`
	ISCSINAME          string `json:"ISCSINAME"`
	ISCSITCPPORT       string `json:"ISCSITCPPORT"`
	LOCATION           string `json:"LOCATION"`
	LOGICTYPE          string `json:"LOGICTYPE"`
	LOSTPACKETS        string `json:"LOSTPACKETS"`
	MACADDRESS         string `json:"MACADDRESS"`
	MTU                string `json:"MTU"`
	NAME               string `json:"NAME"`
	OVERFLOWEDPACKETS  string `json:"OVERFLOWEDPACKETS"`
	OWNINGCONTROLLER   string `json:"OWNINGCONTROLLER"`
	PARENTID           string `json:"PARENTID"`
	PARENTTYPE         int    `json:"PARENTTYPE"`
	PORTSWITCH         string `json:"PORTSWITCH"`
	RUNNINGSTATUS      string `json:"RUNNINGSTATUS"`
	SHARETYPE          string `json:"SHARETYPE"`
	SPEED              string `json:"SPEED"`
	STARTTIME          string `json:"STARTTIME"`
	TYPE               int    `json:"TYPE"`
	CrcErrors          string `json:"crcErrors"`
	DswID              string `json:"dswId"`
	DswLinkRight       string `json:"dswLinkRight"`
	FrameErrors        string `json:"frameErrors"`
	FrameLengthErrors  string `json:"frameLengthErrors"`
	LightStatus        string `json:"lightStatus"`
	MaxSpeed           string `json:"maxSpeed"`
	NumberOfInitiators string `json:"numberOfInitiators"`
	SelectType         string `json:"selectType"`
	WorkModeList       string `json:"workModeList"`
	WorkModeType       string `json:"workModeType"`
	ZoneID             string `json:"zoneId"`
}

var (
	ErrEthernetPortNotFound = "ethernet port is not found"
)

// GetAssociatedEthernetPort get ethernet port associated ASSOCIATEOBJID (maybe port group).
// you must set ASSOCIATEOBJID and ASSOCIATEOBJTYPE. we recommend use dorado.GetPortalIPAddresses().
func (d *Device) GetAssociatedEthernetPort(ctx context.Context, query *SearchQuery) ([]EthernetPort, error) {
	spath := "/eth_port/associate"

	if query == nil || query.AssociateObjType == "" || query.AssociateObjID == "" {
		return nil, errors.New("you must set associated parameter")
	}

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req = AddSearchQuery(req, query)

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(ErrHTTPRequestDo+": %w", err)
	}

	etherports := []EthernetPort{}
	if err = decodeBody(resp, &etherports); err != nil {
		return nil, fmt.Errorf(ErrDecodeBody+": %w", err)
	}

	if len(etherports) == 0 {
		return nil, errors.New(ErrEthernetPortNotFound)
	}

	return etherports, nil
}

// GetPortalIPAddresses get iSCSI portal IP addresses that associated port group.
// return only IPv4 address.
func (d *Device) GetPortalIPAddresses(ctx context.Context, portgroupId int) ([]string, error) {
	query := &SearchQuery{
		AssociateObjID:   strconv.Itoa(portgroupId),
		AssociateObjType: strconv.Itoa(TypePortGroup),
	}

	ethernetports, err := d.GetAssociatedEthernetPort(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to associated ethernet port: %w", err)
	}

	var portalIPs []string
	for _, ethernetport := range ethernetports {
		portalIPs = append(portalIPs, ethernetport.IPV4ADDR)
	}

	return portalIPs, nil
}
