package dorado

import (
	"context"
	"fmt"
)

// System is system information
type System struct {
	CACHEWRITEQUOTA              string `json:"CACHEWRITEQUOTA"`
	CONFIGMODEL                  string `json:"CONFIGMODEL"`
	DESCRIPTION                  string `json:"DESCRIPTION"`
	DOMAINNAME                   string `json:"DOMAINNAME"`
	FREEDISKSCAPACITY            string `json:"FREEDISKSCAPACITY"`
	HEALTHSTATUS                 string `json:"HEALTHSTATUS"`
	HOTSPAREDISKSCAPACITY        string `json:"HOTSPAREDISKSCAPACITY"`
	ID                           string `json:"ID"`
	LOCATION                     string `json:"LOCATION"`
	MEMBERDISKSCAPACITY          string `json:"MEMBERDISKSCAPACITY"`
	NAME                         string `json:"NAME"`
	PRODUCTMODE                  string `json:"PRODUCTMODE"`
	PRODUCTVERSION               string `json:"PRODUCTVERSION"`
	RUNNINGSTATUS                string `json:"RUNNINGSTATUS"`
	SECTORSIZE                   string `json:"SECTORSIZE"`
	STORAGEPOOLCAPACITY          string `json:"STORAGEPOOLCAPACITY"`
	STORAGEPOOLFREECAPACITY      string `json:"STORAGEPOOLFREECAPACITY"`
	STORAGEPOOLHOSTSPARECAPACITY string `json:"STORAGEPOOLHOSTSPARECAPACITY"`
	STORAGEPOOLRAWCAPACITY       string `json:"STORAGEPOOLRAWCAPACITY"`
	STORAGEPOOLUSEDCAPACITY      string `json:"STORAGEPOOLUSEDCAPACITY"`
	THICKLUNSALLOCATECAPACITY    string `json:"THICKLUNSALLOCATECAPACITY"`
	THICKLUNSUSEDCAPACITY        string `json:"THICKLUNSUSEDCAPACITY"`
	THINLUNSALLOCATECAPACITY     string `json:"THINLUNSALLOCATECAPACITY"`
	THINLUNSMAXCAPACITY          string `json:"THINLUNSMAXCAPACITY"`
	THINLUNSUSEDCAPACITY         string `json:"THINLUNSUSEDCAPACITY"`
	TOTALCAPACITY                string `json:"TOTALCAPACITY"`
	TYPE                         int    `json:"TYPE"`
	UNAVAILABLEDISKSCAPACITY     string `json:"UNAVAILABLEDISKSCAPACITY"`
	USEDCAPACITY                 string `json:"USEDCAPACITY"`
	VASAALTERNATENAME            string `json:"VASA_ALTERNATE_NAME"`
	VASASUPPORTBLOCK             string `json:"VASA_SUPPORT_BLOCK"`
	VASASUPPORTFILESYSTEM        string `json:"VASA_SUPPORT_FILESYSTEM"`
	VASASUPPORTPROFILE           string `json:"VASA_SUPPORT_PROFILE"`
	WRITETHROUGHSW               string `json:"WRITETHROUGHSW"`
	WRITETHROUGHTIME             string `json:"WRITETHROUGHTIME"`
	MappedLunsCountCapacity      string `json:"mappedLunsCountCapacity"`
	PatchVersion                 string `json:"patchVersion"`
	UnMappedLunsCountCapacity    string `json:"unMappedLunsCountCapacity"`
	UserFreeCapacity             string `json:"userFreeCapacity"`
	Wwn                          string `json:"wwn"`
}

// GetSystem get system information
func (d *Device) GetSystem(ctx context.Context) (*System, error) {
	spath := "system"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req.URL.Path = req.URL.Path + "/" // path.Join trim last slash

	system := &System{}
	if err = d.requestWithRetry(req, system, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return system, nil
}
