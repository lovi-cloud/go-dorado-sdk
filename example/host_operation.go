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

	fmt.Println("create host")
	host, err := client.LocalDevice.CreateHost(ctx, "w-cn0001")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", host)

	fmt.Println("search host")
	hosts, err := client.LocalDevice.GetHosts(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range hosts {
		fmt.Printf("%+v\n", v)
	}

	fmt.Println("delete host")
	err = client.LocalDevice.DeleteHost(ctx, host.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("operation is done!")
}
