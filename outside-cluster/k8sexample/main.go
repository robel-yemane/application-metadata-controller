package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

func main() {

	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// mux := http.NewServeMux()
	http.HandleFunc("/pods", MakeHandler(pdHandler))
	http.HandleFunc("/namespaces", MakeHandler(nsHandler))
	http.HandleFunc("/deployments", MakeHandler(dpHandler))
	http.HandleFunc("/configmaps", MakeHandler(cmHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))

}
