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

	fmt.Println("get hypermetro domains")
	lhmds, err := client.LocalDevice.GetHyperMetroDomains(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", lhmds)

	rhmds, err := client.RemoteDevice.GetHyperMetroDomains(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", rhmds)

	fmt.Println("operation is done!")
}
