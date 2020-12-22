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

	dummyIqn := "dummyiqn"
	fmt.Println("get initiator force")
	_, err = client.LocalDevice.GetInitiatorForce(ctx, dummyIqn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("delete initiator")
	err = client.LocalDevice.DeleteInitiator(ctx, dummyIqn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("get initiator force (maybe create)")
	_, err = client.LocalDevice.GetInitiatorForce(ctx, dummyIqn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("get initiator force (not create)")
	_, err = client.LocalDevice.GetInitiatorForce(ctx, dummyIqn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("delete initiator")
	err = client.LocalDevice.DeleteInitiator(ctx, dummyIqn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("operation is done!")
}
