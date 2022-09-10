package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/clusterpedia-io/client-go/tools/transport"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	config, err := ctrl.GetConfig()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}
	config.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		return transport.NewTransportForCluster(config.Host, "cluster", rt)
	})

	client, err := clientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to init clientset: %v", err)
	}

	pod, err := client.CoreV1().Pods("demo-system").Get(context.TODO(), "pod1", metav1.GetOptions{})
	if err != nil {
		log.Fatalf("failed to list pods: %v", err)
	}
	fmt.Println(pod)
}
