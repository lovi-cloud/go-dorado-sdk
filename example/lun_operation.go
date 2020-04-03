package main

import (
	"context"
	"fmt"
	"log"

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

	fmt.Println("search lun")
	luns, err := client.LocalDevice.GetLUNByName(ctx, lun.NAME)
	if err != nil {
		log.Fatal(err)
	}

	for _, lun := range luns {
		fmt.Printf("%+v\n", lun)
	}

	fmt.Println("delete lun")
	err = client.LocalDevice.DeleteLUN(ctx, luns[0].ID)
	if err != nil {
		log.Fatal(err)
	}
}
