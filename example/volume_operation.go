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

	fmt.Println("create volume")
	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	hgs, err := client.LocalDevice.GetHyperMetroDomains(ctx)
	if err != nil {
		log.Fatal(err)
	}

	hmp, err := client.CreateVolume(ctx, u, 21, "0", hgs[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", hmp)

	fmt.Println("expand volume")
	err = client.ExtendVolume(ctx, hmp.ID, 30)
	if err != nil {
		log.Fatal(err)
	}

	hmps, err := client.GetHyperMetroPairs(ctx, dorado.CreateSearchId(hmp.ID))
	if err != nil {
		log.Fatal()
	}

	for _, v := range hmps {
		fmt.Printf("%+v\n", v)
	}

	fmt.Println("delete volume")
	err = client.DeleteVolume(ctx, hmp.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("operation is done!")
}
