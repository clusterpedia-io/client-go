package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/clusterpedia-io/client-go/client"
)

func main() {
	config, err := ctrl.GetConfig()
	if err != nil {
		panic(err)
	}
	client, err := client.NewClusterForConfig(config, "cluster1")
	if err != nil {
		panic(err)
	}

	pod, err := client.CoreV1().Pods("default").Get(context.TODO(), "pod1", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(pod)
}
