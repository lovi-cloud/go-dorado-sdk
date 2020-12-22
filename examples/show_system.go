// +build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lovi-cloud/go-dorado-sdk/examples/lib"
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
