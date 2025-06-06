//go:build k8s

package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func auth2K8s() *kubernetes.Clientset {

	// create the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}
