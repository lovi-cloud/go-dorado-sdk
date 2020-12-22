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

	fmt.Println("search mapping view")
	mappingviews, err := client.LocalDevice.GetMappingViews(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", mappingviews)

	fmt.Println("get mapping view")
	mappingview, err := client.LocalDevice.GetMappingView(ctx, mappingviews[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", mappingview)

	fmt.Println("operation is done!")
}
