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

	fmt.Println("search portgroup")
	portgroups, err := client.LocalDevice.GetPortGroups(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", portgroups)

	fmt.Println("get portgroup")
	portgroup, err := client.LocalDevice.GetPortGroup(ctx, portgroups[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", portgroup)

	fmt.Println("operation is done!")
}
