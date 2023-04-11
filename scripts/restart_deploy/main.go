package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type RedeployBody struct {
	ConnectionType string
	Namespace      string
	Deployment     string
}

type gitResponse struct {
	Event      string
	Repository string
	Commit     string
	Ref        string
	Head       string
	Workflow   string
	RequestID  string
	Data       RedeployBody
}

func (r RedeployBody) getInfo() (string, string, string) {
	return r.ConnectionType, r.Namespace, r.Deployment
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
	fmt.Println("Hello World")
}

func payloadRedeploy(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	jsonBody := gitResponse{}
	err := json.Unmarshal([]byte(string(body)), &jsonBody)
	if err != nil {
		fmt.Printf("error getting json from body: %v\n", err)
	}
	connectionType, namespace, deploymentName := jsonBody.Data.getInfo()
	restartDeployHandler(connectionType, namespace, deploymentName)
	fmt.Fprintf(w, "Redeploy in process!")
}

// Endpoint format /redeploy/(inside/outside)/(namespace)/(deployName)
func redeployEndpoint(w http.ResponseWriter, r *http.Request) {
	partialUrl := r.URL.Path[len("/redeploy/"):]
	splitedUrl := strings.Split(partialUrl, "/")
	if len(splitedUrl) != 3 || splitedUrl[2] == "" || (splitedUrl[0] != "inside" && splitedUrl[0] != "outside") {
		fmt.Printf("Error getting the correct number of parameters or the value of connection. Was expecting 3 and got %v \nWas expecting inside or outside as values but got %v \n", len(splitedUrl), splitedUrl[0])
		os.Exit(1)
	}

	connectionType := splitedUrl[0]
	namespace := splitedUrl[1]
	deploymentName := splitedUrl[2]
	restartDeployHandler(connectionType, namespace, deploymentName)
	fmt.Fprintf(w, "Redeploy in process!")
}

func handleRequests() {
	http.HandleFunc("/redeploy/", redeployEndpoint)
	http.HandleFunc("/payload/", payloadRedeploy)
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":4567", nil))
}

func main() {
	handleRequests()
}
