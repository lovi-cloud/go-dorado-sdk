// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/lovi-cloud/go-dorado-sdk/dorado"
	"github.com/lovi-cloud/go-dorado-sdk/example/lib"
)

func main() {
	ctx := context.Background()
	client, err := lib.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	volumes, err := client.GetHyperMetroPairs(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range volumes {
		if isPrefixW(v) {
			// made by this library
			fmt.Println(v.ID)
			err = client.DeleteVolume(ctx, v.ID)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

func isPrefixW(volume dorado.HyperMetroPair) bool {
	return strings.HasPrefix(volume.LOCALOBJNAME, "w-")
}
