package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	network "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v3"
)

type vnetSubscriptionsGroup struct {
	subscriptionID string
	info           []vnetInfo
}

type vnetInfo struct {
	rg   string
	name string
	ip   string
}

func addToVnetGroup(group []vnetSubscriptionsGroup, vnetSubscription string, vnetName string, vnetRG string) []vnetSubscriptionsGroup {
	found := false
	foundPosition := -1
	for i := 0; i < len(group); i++ {
		if vnetSubscription == group[i].subscriptionID {
			found = true
			foundPosition = i
			break
		}
	}

	if found {
		newVnetInfo := vnetInfo{rg: vnetRG, name: vnetName, ip: "NIL"}
		group[foundPosition].info = append(group[foundPosition].info, newVnetInfo)
	} else {
		newVnetInfo := vnetInfo{rg: vnetRG, name: vnetName, ip: "NIL"}
		vnetInfoArray := []vnetInfo{}
		vnetInfoArray = append(vnetInfoArray, newVnetInfo)
		newVnetSubscription := vnetSubscriptionsGroup{subscriptionID: vnetSubscription, info: vnetInfoArray}
		group = append(group, newVnetSubscription)
	}

	return group
}

func getVnetIps(cred *azidentity.ClientSecretCredential, vnetIDs []string) []vnetSubscriptionsGroup {
	ctx := context.Background()
	vnetSplit := []vnetSubscriptionsGroup{}

	for i := 0; i < len(vnetIDs); i++ {
		infoSplit := strings.Split(vnetIDs[i], "/")
		vnetNames := infoSplit[len(infoSplit)-1]
		vnetSubscriptions := infoSplit[2]
		vnetRG := infoSplit[4]
		vnetSplit = addToVnetGroup(vnetSplit, vnetSubscriptions, vnetNames, vnetRG)
	}

	// for i := 0; i < len(vnetSplit); i++ {
	// 	fmt.Printf("Subscriptions = %v\n", vnetSplit[i].subscriptionID)
	// 	fmt.Printf("Info = %v\n", vnetSplit[i].info)
	// }

	// Only retrieved the ip ranges but you can also access all info of the subnets
	for i := 0; i < len(vnetSplit); i++ {
		clientFactory, err := network.NewClientFactory(vnetSplit[i].subscriptionID, cred, nil)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		for j := 0; j < len(vnetSplit[i].info); j++ {
			res, err := clientFactory.NewVirtualNetworksClient().Get(ctx, vnetSplit[i].info[j].rg, vnetSplit[i].info[j].name, nil)
			if err != nil {
				// For testing leave this commnted
				//log.Fatalf("Failed to create client: %v", err)
				fmt.Printf("Failed to create client: %v", err)
			} else {
				updateSpitInfo := vnetSplit[i].info[j]
				updateSpitInfo.ip = *res.Properties.AddressSpace.AddressPrefixes[0]
				vnetSplit[i].info[j] = updateSpitInfo
			}

		}
	}
	return vnetSplit
}
