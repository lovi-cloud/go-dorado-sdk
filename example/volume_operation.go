package main

import (
	"context"
	"fmt"
	"log"

	"github.com/whywaita/go-dorado-sdk/dorado"

	uuid "github.com/satori/go.uuid"
	"github.com/whywaita/go-dorado-sdk/example/lib"
)

func main() {
	ctx := context.Background()

	client, err := lib.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	err = getInitiators(client, ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("operation is done!")
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
	u, err := uuid.NewV4()
	if err != nil {
		return err
	}

	fmt.Println("get volume")
	hgs, err := client.LocalDevice.GetHyperMetroDomains(ctx, dorado.NewSearchQueryName(lib.HyperMetroDomainName))
	if err != nil {
		return err
	}

	fmt.Println("create volume")
	volume, err := client.CreateVolume(ctx, u, 21, lib.StoragePoolName, hgs[0].ID)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", volume)

	fmt.Println("attach volume")
	err = client.AttachVolume(ctx, volume.ID, "", "")
	if err != nil {
		return err
	}

	return nil
}

func volumeOperation(client *dorado.Client, ctx context.Context) error {
	fmt.Println("create volume")
	u, err := uuid.NewV4()
	if err != nil {
		return err
	}

	fmt.Println("get volume")
	hgs, err := client.LocalDevice.GetHyperMetroDomains(ctx, dorado.NewSearchQueryName(lib.HyperMetroDomainName))
	if err != nil {
		return err
	}

	fmt.Println("create volume")
	volume, err := client.CreateVolume(ctx, u, 21, lib.StoragePoolName, hgs[0].ID)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", volume)

	fmt.Println("attach volume")
	err = client.AttachVolume(ctx, volume.ID, "", "")
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
	u, err := uuid.NewV4()
	if err != nil {
		return err
	}

	hgs, err := client.LocalDevice.GetHyperMetroDomains(ctx, dorado.NewSearchQueryName(lib.HyperMetroDomainName))
	if err != nil {
		return err
	}

	hmp, err := client.CreateVolume(ctx, u, 21, lib.StoragePoolName, hgs[0].ID)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", hmp)

	fmt.Println("expand volume")
	err = client.ExtendVolume(ctx, hmp.ID, 30)
	if err != nil {
		return err
	}

	hmps, err := client.GetHyperMetroPairs(ctx, dorado.NewSearchQueryId(hmp.ID))
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
