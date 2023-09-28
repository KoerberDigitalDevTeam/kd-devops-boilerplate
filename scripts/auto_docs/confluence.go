package main

import (
	"log"
	"time"

	confluence "github.com/virtomize/confluence-go-api"
)

func confluenceClient() *confluence.API {
	api, err := confluence.NewAPI("https://koerberdigital.atlassian.net/wiki/rest/api", "pedro.castro@koerber.com", "ATATT3xFfGF0v_1mTCBexcIlA518nxZRqI5TIXAtcxhAa_p2VYSAGlWEcHBj6Tbiv7Y8OTrMov8AWgGWmmB4YaO0yIkDCqam57fPb0IAk--ycgnFnD8ZpgtqG2v0-5jBbTSIuwC3n-3j8QOS21Nv-F6QGDNxPA9LNb7_7Klc-_OlYqIr62JaMGY=1A39EE75")
	if err != nil {
		log.Fatal("Problem with confluence client: %+v", err)
	}

	return api
}

func createPage(cred *confluence.API, title string, body string, confluenceSpace string) {
	timeStampTitle := time.Now().Format("2006-01-02 15:04:05") + " " + title

	data := &confluence.Content{
		Type:  "page",         // can also be page/blogpost
		Title: timeStampTitle, // page title (mandatory)
		Body: confluence.Body{
			Storage: confluence.Storage{
				Value:          body, // your page content here
				Representation: "storage",
			},
		},
		Version: &confluence.Version{ // mandatory
			Number: 1,
		},
		Space: &confluence.Space{
			Key: confluenceSpace, // Space
		},
		Metadata: &confluence.Metadata{
			Properties: &confluence.Properties{
				Editor: &confluence.Editor{
					Key:   "editor",
					Value: "",
				},
				ContentAppearancePublished: &confluence.ContentAppearancePublished{Value: "full-width"},
				ContentAppearanceDraft:     &confluence.ContentAppearanceDraft{Value: "full-width"},
			},
		},
	}

	_, err := cred.CreateContent(data)
	if err != nil {
		log.Fatalf("Error creating confluence page: %v\n", err)
	}
}

func createWanPage(wan []vnetSubscriptionsGroup, confluenceSpace string) string {

	var trName []string
	var trSubscription []string
	var trRG []string
	var trIp []string

	for i := 0; i < len(wan); i++ {
		for j := 0; j < len(wan[i].info); j++ {
			trName = append(trName, "<td>"+wan[i].info[j].name+"</td>")
			trSubscription = append(trSubscription, "<td>"+wan[i].subscriptionID+"</td>")
			trRG = append(trRG, "<td>"+wan[i].info[j].rg+"</td>")
			trIp = append(trIp, "<td>"+wan[i].info[j].ip+"</td>")
		}
	}

	body := `
	<html>
<head>
<style>
table, th, td {
  border: 1px solid black;
}
</style>
</head>
<body>

<h1>Wan Info</h1>

<table>
  <tr>
    <th>Name</th>
	<th>Subscription</th>
	<th>Resource Group</th>
	<th>Ip Range</th>
  </tr>
	`

	for i := 0; i < len(trName); i++ {
		body = body + "<tr>" + trName[i] + trSubscription[i] + trRG[i] + trIp[i] + "</tr>"
	}
	body = body + "</table></body></html>"
	return body
}
