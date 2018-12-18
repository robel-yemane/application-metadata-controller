package main

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {

	clientset := auth2K8s()
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// list all deployments in the default namespace
	defaultNamespace := apiv1.NamespaceDefault
	deploymentClient := clientset.AppsV1().Deployments(defaultNamespace)
	fmt.Printf("\n=> Listing deployments in namespace %q:\n\n", defaultNamespace)
	list, err := deploymentClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}

	// list all namespaces
	nsClient := clientset.CoreV1().Namespaces()
	fmt.Printf("\n=> Listing all namespaces: \n\n")
	ns, err := nsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, n := range ns.Items {
		fmt.Printf(" * %s\n", n.Name)
	}

	// list pods in default namespace
	fmt.Printf("\n=>Listing Pods in namespace %q:\n\n", defaultNamespace)
	podClient := clientset.CoreV1().Pods(defaultNamespace)
	podList, err := podClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, pod := range podList.Items {
		fmt.Printf(" * %s\n", pod.Name)
	}

	fmt.Println("\n=>Listing Pods in namespace kube-system:\n")

	podClientII := clientset.CoreV1().Pods("kube-system")
	podListII, err := podClientII.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, pod := range podListII.Items {
		fmt.Printf(" * %s\n", pod.Name)
	}
}
