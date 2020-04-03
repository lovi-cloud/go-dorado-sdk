package main

import (
	"context"
	"fmt"
	"log"

	"github.com/whywaita/go-dorado-sdk/example/lib"
)

func main() {
	client, err := lib.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	sps, err := client.LocalDevice.GetStoragePools(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range sps {
		fmt.Printf("%+v\n", v)
	}
}
