package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig *string

func makeHandler(fn func(http.ResponseWriter, *http.Request, *restclient.Config, *kubernetes.Clientset)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		fn(w, r, config, clientset)

	}

}

func handler(w http.ResponseWriter, r *http.Request, config *restclient.Config, clientset *kubernetes.Clientset) {

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	if len(pods.Items) > 0 {
		for _, pod := range pods.Items {
			fmt.Fprintf(w, "Pod Name:- %s\n", pod.Name)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "pod items = %d", len(pods.Items))
	}
}

func nsHandler(w http.ResponseWriter, r *http.Request, config *restclient.Config, clientset *kubernetes.Clientset) {

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

func main() {

	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// mux := http.NewServeMux()
	http.HandleFunc("/pods", makeHandler(handler))
	http.HandleFunc("/ns", makeHandler(nsHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))

}
