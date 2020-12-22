// +build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lovi-cloud/go-dorado-sdk/dorado"

	uuid "github.com/satori/go.uuid"
	"github.com/lovi-cloud/go-dorado-sdk/examples/lib"
)

func main() {
	ctx := context.Background()
	var err error

	client, err := lib.GetClient()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(client)

	fmt.Println("operation is done!")
}

func truncateAllVolume(ctx context.Context, client *dorado.Client) error {
	volumes, err := client.GetHyperMetroPairs(ctx, nil)
	if err != nil && err != dorado.ErrHyperMetroPairNotFound {
		return err
	}

	for _, v := range volumes {
		fmt.Println(v.ID)
		if err := client.DeleteVolume(ctx, v.ID); err != nil {
			return err
		}
	}

	llc, err := client.LocalDevice.GetLUNCopys(ctx, nil)
	if err != nil && err != dorado.ErrLunCopyNotFound {
		return err
	}
	for _, lc := range llc {
		fmt.Println(lc.ID)
		if err := client.LocalDevice.DeleteLUNCopy(ctx, lc.ID); err != nil {
			return err
		}
	}

	ls, err := client.LocalDevice.GetSnapshots(ctx, nil)
	if err != nil && err != dorado.ErrSnapshotNotFound {
		return err
	}
	for _, s := range ls {
		fmt.Println(s.ID)
		if err := client.LocalDevice.DeleteSnapshot(ctx, s.ID); err != nil {
			return err
		}
	}

	ll, err := client.LocalDevice.GetLUNs(ctx, nil)
	if err != nil && err != dorado.ErrLunNotFound {
		return err
	}
	for _, l := range ll {
		fmt.Println(l.ID)
		if err := client.LocalDevice.DeleteLUN(ctx, l.ID); err != nil {
			return err
		}
	}

	rlc, err := client.RemoteDevice.GetLUNCopys(ctx, nil)
	if err != nil && err != dorado.ErrLunCopyNotFound {
		return err
	}
	for _, lc := range rlc {
		fmt.Println(lc.ID)
		if err := client.RemoteDevice.DeleteLUNCopy(ctx, lc.ID); err != nil {
			return err
		}
	}

	rs, err := client.RemoteDevice.GetSnapshots(ctx, nil)
	if err != nil && err != dorado.ErrSnapshotNotFound {
		return err
	}
	for _, s := range rs {
		fmt.Println(s.ID)
		if err := client.RemoteDevice.DeleteSnapshot(ctx, s.ID); err != nil {
			return err
		}
	}

	rl, err := client.RemoteDevice.GetLUNs(ctx, nil)
	if err != nil && err != dorado.ErrLunNotFound {
		return err
	}
	for _, l := range rl {
		fmt.Println(l.ID)
		if err := client.RemoteDevice.DeleteLUN(ctx, l.ID); err != nil {
			return err
		}
	}

	return nil
}

func truncateAttachedDevice(ctx context.Context, client *dorado.Client) error {
	// local
	_, host, err := client.LocalDevice.GetHostGroupForce(ctx, "isucn0001")
	if err != nil {
		return err
	}

	luns, err := client.LocalDevice.GetHostAssociatedLUNs(ctx, host.ID)
	fmt.Printf("Associated Number is %d\n", len(luns))

	hmps, err := client.GetHyperMetroPairs(ctx, nil)
	if err != nil {
		return err
	}

	var associated []dorado.HyperMetroPair
	for _, hmp := range hmps {
		for _, lun := range luns {
			if hmp.LOCALOBJID == lun.ID {
				associated = append(associated, hmp)
			}
		}
	}

	for _, hmp := range associated {
		fmt.Printf("deleted: %s\n", hmp.ID)
		if err := client.DeleteVolume(ctx, hmp.ID); err != nil {
			return err
		}
	}

	return nil
}

func getInitiators(client *dorado.Client, ctx context.Context) error {
	initiator, err := client.LocalDevice.GetInitiatorForce(ctx, "")
	if err != nil {
		return err
	}

	fmt.Println(initiator)
	return nil
}

func attachVolume(client *dorado.Client, ctx context.Context) error {
	fmt.Println("create volume")
	u := uuid.NewV4()

	fmt.Println("get volume")
	hgs, err := client.LocalDevice.GetHyperMetroDomains(ctx, dorado.NewSearchQueryName(lib.HyperMetroDomainName))
	if err != nil {
		return err
	}

	fmt.Println("create volume")
	volume, err := client.CreateVolumeRaw(ctx, u, 21, lib.StoragePoolName, hgs[0].ID)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", volume)

	fmt.Println("attach volume")
	err = client.AttachVolume(ctx, volume.ID, "w-cn0001", "dummy-iqn")
	if err != nil {
		return err
	}

	return nil
}

func volumeOperation(client *dorado.Client, ctx context.Context) error {
	fmt.Println("create volume")
	u := uuid.NewV4()

	fmt.Println("get volume")
	hgs, err := client.LocalDevice.GetHyperMetroDomains(ctx, dorado.NewSearchQueryName(lib.HyperMetroDomainName))
	if err != nil {
		return err
	}

	fmt.Println("create volume")
	volume, err := client.CreateVolumeRaw(ctx, u, 21, lib.StoragePoolName, hgs[0].ID)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", volume)

	fmt.Println("attach volume")
	err = client.AttachVolume(ctx, volume.ID, "w-cn0001", "dummy-iqn")
	if err != nil {
		return err
	}

	fmt.Println("detach volume")
	err = client.DetachVolume(ctx, volume.ID)
	if err != nil {
		return err
	}

	return nil
}

func singleLunOperation(client *dorado.Client, ctx context.Context) error {
	fmt.Println("create volume")
	u := uuid.NewV4()

	hgs, err := client.LocalDevice.GetHyperMetroDomains(ctx, dorado.NewSearchQueryName(lib.HyperMetroDomainName))
	if err != nil {
		return err
	}

	hmp, err := client.CreateVolumeRaw(ctx, u, 21, lib.StoragePoolName, hgs[0].ID)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", hmp)

	fmt.Println("expand volume")
	err = client.ExtendVolume(ctx, hmp.ID, 30)
	if err != nil {
		return err
	}

	hmps, err := client.GetHyperMetroPairs(ctx, dorado.NewSearchQueryID(hmp.ID))
	if err != nil {
		return err
	}

	for _, v := range hmps {
		fmt.Printf("%+v\n", v)
	}

	fmt.Println("delete volume")
	err = client.DeleteVolume(ctx, hmp.ID)
	if err != nil {
		return err
	}

	return nil
}
