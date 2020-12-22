// +build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lovi-cloud/go-dorado-sdk/examples/lib"
)

func main() {
	ctx := context.Background()

	client, err := lib.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	targetports, err := client.LocalDevice.GetTargetPort(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range targetports {
		fmt.Printf("%+v\n", t)
	}

	targetIqns, err := client.LocalDevice.GetTargetIQNs(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range targetIqns {
		fmt.Printf("%+v\n", t)
	}

	fmt.Println("operation is done!")
}
