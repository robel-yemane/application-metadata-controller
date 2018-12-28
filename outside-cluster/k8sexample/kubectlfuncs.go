package main

import (
	"fmt"
	"net/http"
	"os"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig *string
var defaultNamespace = apiv1.NamespaceDefault

func MakeHandler(fn func(http.ResponseWriter, *http.Request, *restclient.Config, *kubernetes.Clientset)) http.HandlerFunc {

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

//list pods in given namespace
func pdHandler(w http.ResponseWriter, r *http.Request, config *restclient.Config, clientset *kubernetes.Clientset) {

	namespace := r.URL.Query().Get("ns")
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	if len(pods.Items) <= 0 {

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "No resources found")
		return
	}
	fmt.Fprint(w, "NAME\tREADY\tSTATUS\n\n")
	for _, pod := range pods.Items {
		fmt.Fprintf(w, "%10s\t%10s\t%10s\n", pod.Name, pod.Status.Conditions[1].Status, pod.Status.Phase)
	}

}

//list namespaces
func nsHandler(w http.ResponseWriter, r *http.Request, config *restclient.Config, clientset *kubernetes.Clientset) {

	nsClient := clientset.CoreV1().Namespaces()
	ns, err := nsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, n := range ns.Items {
		fmt.Fprintf(w, " * %s\n", n.Name)
	}

}

// list deployments in given namespace
func dpHandler(w http.ResponseWriter, r *http.Request, config *restclient.Config, clientset *kubernetes.Clientset) {
	namespace := r.URL.Query().Get("ns")
	deploymentClient := clientset.AppsV1().Deployments(namespace)
	list, err := deploymentClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Fprintf(w, " * %s (%d/%d replicas)\n", d.Name, d.Status.Replicas, *d.Spec.Replicas)
	}
}

//list configmaps in given namespace
func cmHandler(w http.ResponseWriter, r *http.Request, config *restclient.Config, clientset *kubernetes.Clientset) {
	namespace := r.URL.Query().Get("ns")
	cmClient := clientset.CoreV1().ConfigMaps(namespace)
	cm, err := cmClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	if len(cm.Items) <= 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "No resources found")
		return
	}
	for _, c := range cm.Items {
		fmt.Fprintf(w, " * %s	-> %s\n", c.Namespace, c.Name)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// func main() {

// 	if home := homeDir(); home != "" {
// 		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
// 	} else {
// 		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
// 	}
// 	flag.Parse()

// 	// mux := http.NewServeMux()
// 	http.HandleFunc("/pods", makeHandler(pdHandler))
// 	http.HandleFunc("/namespaces", makeHandler(nsHandler))
// 	http.HandleFunc("/deployments", makeHandler(dpHandler))
// 	http.HandleFunc("/configmaps", makeHandler(cmHandler))
// 	log.Fatal(http.ListenAndServe(":8080", nil))

// }
