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

	system, err := client.LocalDevice.GetSystem(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", *system)
}
