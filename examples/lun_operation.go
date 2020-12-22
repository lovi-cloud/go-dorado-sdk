// +build ignore

package main

import (
	"context"
	"fmt"
	"log"

	uuid "github.com/satori/go.uuid"
	"github.com/lovi-cloud/go-dorado-sdk/dorado"

	"github.com/lovi-cloud/go-dorado-sdk/examples/lib"
)

func main() {
	ctx := context.Background()

	client, err := lib.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	err = lunOperation(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("operation is done!")
}

func lunOperation(ctx context.Context, client *dorado.Client) error {
	fmt.Println("create lun")
	id := uuid.NewV4()
	lun, err := client.LocalDevice.CreateLUN(ctx, id, 21, lib.StoragePoolName)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", lun)

	fmt.Println("search lun")
	luns, err := client.LocalDevice.GetLUNs(ctx, dorado.NewSearchQueryName(lun.NAME))
	if err != nil {
		return err
	}

	for _, l := range luns {
		fmt.Printf("%+v\n", l)
	}

	fmt.Println("create lungroup")
	lungroup, err := client.LocalDevice.CreateLunGroup(ctx, "w-cn0001")
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", lungroup)

	fmt.Println("Associate lungroup")
	err = client.LocalDevice.AssociateLun(ctx, lungroup.ID, lun.ID)
	if err != nil {
		return err
	}

	lungroup2, err := client.LocalDevice.GetLunGroup(ctx, lungroup.ID)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", lungroup2)

	fmt.Println("DisAssociate lungroup")
	err = client.LocalDevice.DisAssociateLun(ctx, lungroup2.ID, lun.ID)
	if err != nil {
		return err
	}

	fmt.Println("delete lun")
	err = client.LocalDevice.DeleteLUN(ctx, luns[0].ID)
	if err != nil {
		return err
	}

	fmt.Println("delete lungroup")
	err = client.LocalDevice.DeleteLunGroup(ctx, lungroup2.ID)
	if err != nil {
		return err
	}

	return nil
}
