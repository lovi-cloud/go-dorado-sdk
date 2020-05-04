// +build ignore

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

	fmt.Println("search initiator")
	initiators, err := client.LocalDevice.GetInitiators(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", initiators)

	fmt.Println("get initiator")
	initiator, err := client.LocalDevice.GetInitiator(ctx, initiators[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", initiator)

	fmt.Println("operation is done!")
}
