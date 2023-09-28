package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	network "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

func checkVnet(cred *azidentity.ClientSecretCredential) []string {
	subscription := "f4cb84de-c926-4736-8743-74561684ce48"
	vnetIDs := []string{}
	ctx := context.Background()
	var wanIDS []string
	client, err := network.NewVirtualWansClient(subscription, cred, nil)
	if err != nil {
		fmt.Printf("Error creating wan cliend: %v\n", err)
		os.Exit(1)
	}

	//Collection of all the hub ids in the wan
	wanPage := client.NewListPager(nil)
	for wanPage.More() {
		nextResult, err := wanPage.NextPage(ctx)
		if err != nil {
			fmt.Printf("failed to advance page: %v", err)
		}
		for _, v := range nextResult.Value {
			hubs := v.Properties.VirtualHubs
			for i := 0; i < len(hubs); i++ {
				wanIDS = append(wanIDS, *hubs[i].ID)
			}
		}
	}
	//Filter the data for the name of the Hubs
	var hubNames []string
	var resourceGroupNames []string
	for i := range wanIDS {
		nameSplit := strings.Split(wanIDS[i], "/")
		name := nameSplit[len(nameSplit)-1]
		rg := nameSplit[4]
		hubNames = append(hubNames, name)
		resourceGroupNames = append(resourceGroupNames, rg)
	}

	hubClient, err := network.NewHubVirtualNetworkConnectionsClient(subscription, cred, nil)
	if err != nil {
		fmt.Printf("Faliled to get connection client: %v", err)
	}

	//Filter all the hubs for the connected vnet ids
	for i := 0; i < len(resourceGroupNames); i++ {
		hubPager := hubClient.NewListPager(resourceGroupNames[i], hubNames[i], nil)
		for hubPager.More() {
			nextResult, err := hubPager.NextPage(ctx)
			if err != nil {
				log.Fatalf("failed to advance page: %v", err)
			}
			for _, v := range nextResult.Value {
				vnetIDs = append(vnetIDs, *v.Properties.RemoteVirtualNetwork.ID)
			}
		}
	}

	return vnetIDs
}
