//go:build k8s

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

	fmt.Println("\n=>Listing Pods in namespace kube-system - with status running:\n")

	podClientII := clientset.CoreV1().Pods("kube-system")
	listOptions := metav1.ListOptions{
		FieldSelector: "status.phase=Running",
	}
	podListII, err := podClientII.List(listOptions)
	if err != nil {
		panic(err)
	}

	for _, pod := range podListII.Items {

		fmt.Printf(" * %s\n", pod.Name)
	}

	cmClient := clientset.CoreV1().ConfigMaps(apiv1.NamespaceAll)

	fmt.Printf("\n=> Listing all CMs and their NS: \n\n")
	cm, err := cmClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, c := range cm.Items {
		fmt.Printf(" * %s	-> %s\n", c.Namespace, c.Name)
	}
	//------//
	// list a namespace
	ks, err := clientset.CoreV1().Services("kube-system").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("\nkubectl cluster-info ")
	for _, something := range ks.Items {
		fmt.Printf(" %s -> %s\n", something.GetObjectMeta().GetName(), something.GetObjectMeta().GetSelfLink())
	}

	rs, err := clientset.Extensions().ReplicaSets(apiv1.NamespaceDefault).List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("\n kubectl get replicasets\n")
	fmt.Printf("%s\t\t%s\t\t%s\t\t%s\n", "NAME", "DESIRED", "CURRENT", "READY")
	for _, rst := range rs.Items {
		fmt.Printf("%s\t%d\t%d\t%d\n", rst.Name, *rst.Spec.Replicas, rst.Status.Replicas, rst.Status.ReadyReplicas)
	}

}
