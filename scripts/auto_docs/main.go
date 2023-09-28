package main

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func connectToAzure() *azidentity.ClientSecretCredential {
	tenant_id := "dca04e92-b0c8-46b4-9d33-f893be36d972"
	client_id := "f4306dcb-71c4-493f-a6d0-500889530ec6"
	clientSecret := "mnw8Q~1ENZgu1a0oZ2sKpqWUAXSCbbtzE~Cjrcm-"

	cred, err := azidentity.NewClientSecretCredential(tenant_id, client_id, clientSecret, nil)
	//cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Printf("Error getting credentials: %v\n", err)
		os.Exit(1)
	}
	return cred
}

func main() {
	confluenceSpace := "~62553ca1ed2b3e0074fc2d7d"
	cred := connectToAzure()
	vnetIDs := checkVnet(cred)
	vnetGroups := getVnetIps(cred, vnetIDs)
	confluenceCred := confluenceClient()

	tbl := createWanPage(vnetGroups, confluenceSpace)
	createPage(confluenceCred, "Wan Tables", tbl, confluenceSpace)

}
