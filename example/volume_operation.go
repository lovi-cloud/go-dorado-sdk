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

	fmt.Println("create volume")
	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	hgs, err := client.LocalDevice.GetHyperMetroDomain(ctx)
	if err != nil {
		log.Fatal(err)
	}

	hmp, err := client.CreateVolume(ctx, u, 21, "0", hgs[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", hmp)

	fmt.Println("delete volume")
	err = client.DeleteVolume(ctx, hmp.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("operation is done!")
}
