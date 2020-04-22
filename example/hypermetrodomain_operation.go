package main

import (
	"context"
	"fmt"
	"log"

	"github.com/whywaita/go-dorado-sdk/example/lib"
)

func main() {
	ctx := context.Background()
	client, err := lib.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("get hypermetro domains")
	hmds, err := client.LocalDevice.GetHyperMetroDomains(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", hmds)

	fmt.Println("operation is done!")
}
