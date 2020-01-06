package katy

import (
	"log"
	"os"

	"github.com/VisImag/katy/controller/pods"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubernetesClient() *kubernetes.Clientset {

	// Build client configuration
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBE_CONFIG_PATH"))
	if err != nil {
		log.Fatalln(err)
	}

	// Get client object
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

func PushKubernetesClient(client *kubernetes.Clientset) {

	// Push kubernetes clientset to pod controller
	pods.PushClientSet(client)
}
