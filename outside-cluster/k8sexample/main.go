package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	http.HandleFunc("/", handler)
	http.HandleFunc("/ns", nsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func handler(w http.ResponseWriter, r *http.Request) {

	var kubeconfig *string

	temp := "/Users/ryemane/.kube/config"

	kubeconfig = &temp
	//kubeconfig := &

	// const kubeconfig *string
	// if home := homeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// }
	// flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range pods.Items {
		fmt.Fprintf(w, "Pod Name:- %s\n", pod.Name)
	}
}

func nsHandler(w http.ResponseWriter, r *http.Request) {

	var kubeconfig *string

	temp := "/Users/ryemane/.kube/config"

	kubeconfig = &temp
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	nsClient := clientset.CoreV1().Namespaces()
	fmt.Printf("\n=> Listing all namespaces: \n\n")
	ns, err := nsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, n := range ns.Items {
		fmt.Fprintf(w, " * %s\n", n.Name)
	}

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
