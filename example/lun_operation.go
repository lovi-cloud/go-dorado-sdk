// +build ignore

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

	fmt.Println("create lun")
	id, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	lun, err := client.LocalDevice.CreateLUN(ctx, id, 21, "0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", lun)

	fmt.Println("search lun")
	luns, err := client.LocalDevice.GetLUNs(ctx, dorado.NewSearchQueryName(lun.NAME))
	if err != nil {
		log.Fatal(err)
	}

	for _, l := range luns {
		fmt.Printf("%+v\n", l)
	}

	fmt.Println("create lungroup")
	lungroup, err := client.LocalDevice.CreateLunGroup(ctx, "w-cn0001")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", lungroup)

	fmt.Println("Associate lungroup")
	err = client.LocalDevice.AssociateLun(ctx, lungroup.ID, lun.ID)
	if err != nil {
		log.Fatal(err)
	}

	lungroup2, err := client.LocalDevice.GetLunGroup(ctx, lungroup.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", lungroup2)

	fmt.Println("DisAssociate lungroup")
	err = client.LocalDevice.DisAssociateLun(ctx, lungroup2.ID, lun.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("delete lun")
	err = client.LocalDevice.DeleteLUN(ctx, luns[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("delete lungroup")
	err = client.LocalDevice.DeleteLunGroup(ctx, lungroup2.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("operation is done!")
}
