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

	fmt.Println("create host")
	host, err := client.LocalDevice.CreateHost(ctx, "w-cn0001")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", host)

	fmt.Println("search host")
	hosts, err := client.LocalDevice.GetHost(ctx, host.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", hosts)

	fmt.Println("create hostgroup")
	hostgroup, err := client.LocalDevice.CreateHostGroup(ctx, "w-cn0001")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", hostgroup)

	fmt.Println("search hostgroup")
	hostgroups, err := client.LocalDevice.GetHostGroup(ctx, hostgroup.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", hostgroups)

	fmt.Println("associate host")
	err = client.LocalDevice.AssociateHost(ctx, hostgroup.ID, host.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("disassociate host")
	err = client.LocalDevice.DisAssociateHost(ctx, hostgroup.ID, host.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("delete host")
	err = client.LocalDevice.DeleteHost(ctx, host.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("delete hostgroup")
	err = client.LocalDevice.DeleteHostGroup(ctx, hostgroup.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("operation is done!")
}
