package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func restartDeployHandler(connectionType string, namespace string, deploymentName string) {
	//This is for debugging. Do not uncomment
	// connectionType := os.Args[1]
	// namespace := os.Args[2]
	// deployment_name := os.Args[3]

	var clientset *kubernetes.Clientset

	if connectionType == "inside" {
		clientset = insideCluster()
	} else if connectionType == "outside" {
		clientset = outsideCluster()
	} else {
		fmt.Println("Please use either outside or inside as the first argument. This is meant to tell the program if it should connect using ~/.kube/config or the internal k8s rest client")
		os.Exit(1)
	}

	// This is to get the current docker image of the pod, it is not necessary now but it might be usefull in the future
	//pods, _ := clientset.CoreV1().Pods(namespace).List(context.Background(), v1.ListOptions{})
	// for _, pod := range pods.Items {
	// 	if strings.Contains(pod.Name, deploymentName) {
	// 		imageIdWithSha := pod.Status.ContainerStatuses[0].ImageID
	// 		fmt.Printf(imageIdWithSha)
	// 	}
	// }

	// deploymentsClient := clientset.AppsV1().Deployments(namespace)
	// deployment, _ := deploymentsClient.Get(context.Background(),deployment_name, v1.GetOptions{} )

	restartDeploy(clientset, namespace, deploymentName)
}

func restartDeploy(clientset *kubernetes.Clientset, namespace string, deploymentName string) {
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	data := fmt.Sprintf(`{"spec": {"template": {"metadata": {"annotations": {"kubectl.kubernetes.io/restartedAt": "%s"}}}}}`, time.Now().Format("20060102150405"))
	_, err := deploymentsClient.Patch(context.Background(), deploymentName, types.MergePatchType, []byte(data), v1.PatchOptions{})

	if err != nil {
		fmt.Printf("error getting deployments patch: %v\n", err)
		os.Exit(1)
	}
}

func outsideCluster() *kubernetes.Clientset {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		os.Exit(1)
	}
	// Same as ~/.kube/config
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("error getting Kubernetes config: %v\n", err)
		os.Exit(1)
	}
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Printf("Error getting clientSet: %v\n", err)
		os.Exit(1)
	}

	return clientset
}
func insideCluster() *kubernetes.Clientset {

	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("error getting Kubernetes config: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Printf("Error getting clientSet: %v\n", err)
		os.Exit(1)
	}

	return clientset
}
