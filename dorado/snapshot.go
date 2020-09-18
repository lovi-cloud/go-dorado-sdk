package dorado

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Snapshot is object of lun snapshot
type Snapshot struct {
	CASCADEDLEVEL         string `json:"CASCADEDLEVEL"`
	CASCADEDNUM           string `json:"CASCADEDNUM"`
	CONSUMEDCAPACITY      string `json:"CONSUMEDCAPACITY"`
	DESCRIPTION           string `json:"DESCRIPTION"`
	EXPOSEDTOINITIATOR    string `json:"EXPOSEDTOINITIATOR"`
	HEALTHSTATUS          string `json:"HEALTHSTATUS"`
	HYPERCOPYIDS          string `json:"HYPERCOPYIDS"`
	ID                    int    `json:"ID,string"`
	IOCLASSID             string `json:"IOCLASSID"`
	IOPRIORITY            string `json:"IOPRIORITY"`
	ISSCHEDULEDSNAP       string `json:"ISSCHEDULEDSNAP"`
	NAME                  string `json:"NAME"`
	PARENTID              int    `json:"PARENTID,string"`
	PARENTNAME            string `json:"PARENTNAME"`
	PARENTTYPE            int    `json:"PARENTTYPE"`
	ROLLBACKENDTIME       string `json:"ROLLBACKENDTIME"`
	ROLLBACKRATE          string `json:"ROLLBACKRATE"`
	ROLLBACKSPEED         string `json:"ROLLBACKSPEED"`
	ROLLBACKSTARTTIME     string `json:"ROLLBACKSTARTTIME"`
	ROLLBACKTARGETOBJID   string `json:"ROLLBACKTARGETOBJID"`
	ROLLBACKTARGETOBJNAME string `json:"ROLLBACKTARGETOBJNAME"`
	RUNNINGSTATUS         string `json:"RUNNINGSTATUS"`
	SOURCELUNCAPACITY     string `json:"SOURCELUNCAPACITY"`
	SOURCELUNID           string `json:"SOURCELUNID"`
	SOURCELUNNAME         string `json:"SOURCELUNNAME"`
	SUBTYPE               string `json:"SUBTYPE"`
	TIMESTAMP             string `json:"TIMESTAMP"`
	TYPE                  int    `json:"TYPE"`
	USERCAPACITY          string `json:"USERCAPACITY"`
	WORKINGCONTROLLER     string `json:"WORKINGCONTROLLER"`
	WORKLOADTYPEID        string `json:"WORKLOADTYPEID"`
	WORKLOADTYPENAME      string `json:"WORKLOADTYPENAME"`
	WWN                   string `json:"WWN"`
	ReplicationCapacity   string `json:"replicationCapacity"`
	SnapCgID              string `json:"snapCgId"`
}

// EncodeSnapshotName encode compatible name
func EncodeSnapshotName(u uuid.UUID) string {
	return EncodeLunName(u)
}

// GetSnapshots get snapshots by query
func (d *Device) GetSnapshots(ctx context.Context, query *SearchQuery) ([]Snapshot, error) {
	spath := "/snapshot"

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}
	req = AddSearchQuery(req, query)

	var snapshots []Snapshot
	if err = d.requestWithRetry(req, &snapshots, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	if len(snapshots) == 0 {
		return nil, ErrSnapshotNotFound
	}

	return snapshots, nil
}

// GetSnapshot get snapshot by id
func (d *Device) GetSnapshot(ctx context.Context, snapshotID int) (*Snapshot, error) {
	spath := fmt.Sprintf("/snapshot/%d", snapshotID)

	req, err := d.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	snapshots := &Snapshot{}
	if err = d.requestWithRetry(req, snapshots, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return snapshots, nil
}

// CreateSnapshot create object of snapshot
func (d *Device) CreateSnapshot(ctx context.Context, lunID int, name uuid.UUID, description string) (*Snapshot, error) {
	spath := "/snapshot"
	param := struct {
		TYPE        string `json:"TYPE"`
		NAME        string `json:"NAME"`
		PARENTTYPE  string `json:"PARENTTYPE"`
		PARENTID    string `json:"PARENTID"`
		DESCRIPTION string `json:"DESCRIPTION"`
	}{
		TYPE:        strconv.Itoa(TypeSnapshot),
		NAME:        EncodeSnapshotName(name),
		PARENTTYPE:  strconv.Itoa(TypeLUN),
		PARENTID:    strconv.Itoa(lunID),
		DESCRIPTION: description,
	}
	jb, err := json.Marshal(param)
	if err != nil {
		return nil, fmt.Errorf(ErrCreatePostValue+": %w", err)
	}

	req, err := d.newRequest(ctx, "POST", spath, bytes.NewBuffer(jb))
	if err != nil {
		return nil, fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	snapshot := &Snapshot{}
	if err = d.requestWithRetry(req, snapshot, DefaultHTTPRetryCount); err != nil {
		return nil, fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return snapshot, nil
}

// CreateSnapshotWithWait create snapshot and waiting ready
func (d *Device) CreateSnapshotWithWait(ctx context.Context, lunID int, name uuid.UUID, description string) (*Snapshot, error) {
	snapshot, err := d.CreateSnapshot(ctx, lunID, name, description)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}

	// wait 10 seconds
	for i := 0; i < 10; i++ {
		isReady, err := d.snapshotIsReady(ctx, snapshot.ID)
		if err != nil {
			if err := d.DeleteSnapshot(ctx, snapshot.ID); err != nil {
				d.Logger.Printf("failed to delete snapshot: %v\n", err)
			}
			return nil, fmt.Errorf("failed to wait that snapshot is ready: %w", err)
		}

		if isReady == true {
			return d.GetSnapshot(ctx, snapshot.ID)
		}

		time.Sleep(1 * time.Second)
	}

	return nil, ErrTimeoutWait
}

func (d *Device) snapshotIsReady(ctx context.Context, snapshotID int) (bool, error) {
	snapshot, err := d.GetSnapshot(ctx, snapshotID)
	if err != nil {
		return false, fmt.Errorf("failed to get snapshot (ID: %d): %w", snapshotID, err)
	}

	if snapshot.HEALTHSTATUS != strconv.Itoa(StatusHealth) {
		return false, fmt.Errorf("snapshot health status is bad (HEALTHSTATUS: %s)", snapshot.HEALTHSTATUS)
	}

	if snapshot.RUNNINGSTATUS == strconv.Itoa(StatusSnapshotActive) || snapshot.RUNNINGSTATUS == strconv.Itoa(StatusSnapshotInactive) {
		return true, nil
	}

	return false, nil
}

// DeleteSnapshot delete snapshot
func (d *Device) DeleteSnapshot(ctx context.Context, snapshotID int) error {
	spath := fmt.Sprintf("/snapshot/%d", snapshotID)
	param := struct {
		TYPE string `json:"TYPE"`
		ID   string `json:"ID"`
	}{
		TYPE: strconv.Itoa(TypeSnapshot),
		ID:   strconv.Itoa(snapshotID),
	}
	jb, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf(ErrCreatePostValue+": %w", err)
	}

	req, err := d.newRequest(ctx, "DELETE", spath, bytes.NewBuffer(jb))
	if err != nil {
		return fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	var i interface{} // this endpoint return N/A
	if err = d.requestWithRetry(req, i, DefaultHTTPRetryCount); err != nil {
		return fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return nil
}

// ActivateSnapshot activate snapshot
func (d *Device) ActivateSnapshot(ctx context.Context, snapshotID int) error {
	spath := "/snapshot/activate"
	param := struct {
		SNAPSHOTLIST []string `json:"SNAPSHOTLIST"`
	}{
		SNAPSHOTLIST: []string{strconv.Itoa(snapshotID)},
	}
	jb, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf(ErrCreatePostValue+": %w", err)
	}

	req, err := d.newRequest(ctx, "POST", spath, bytes.NewBuffer(jb))
	if err != nil {
		return fmt.Errorf(ErrCreateRequest+": %w", err)
	}

	var i interface{} // this endpoint return N/A
	if err = d.requestWithRetry(req, i, DefaultHTTPRetryCount); err != nil {
		return fmt.Errorf(ErrRequestWithRetry+": %w", err)
	}

	return nil
}

// StopSnapshot stop snapshot
func (d *Device) StopSnapshot(ctx context.Context, snapshotID int) error {
	spath := "/snapshot/stop"
	param := struct {
		ID string `json:"ID"`
	}{
		ID: strconv.Itoa(snapshotID),
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
