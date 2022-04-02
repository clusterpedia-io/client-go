package main

import (
	"context"
	"fmt"

	"github.com/clusterpedia-io/client-go/client"
	fake "github.com/clusterpedia-io/client-go/fake/apiserver"
	"github.com/clusterpedia-io/client-go/fake/storage/memory"
	"github.com/clusterpedia-io/client-go/tools/builder"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

func main() {
	// build fake memory storage
	storage := memory.NewFakeStorageFactory()
	server, err := fake.NewFakeApiserver(storage)
	if err != nil {
		panic(err)
	}
	defer server.Close()

	if err = storage.Create(context.TODO(), "cluster-01", NewObject("pod1")); err != nil {
		panic(err)
	}

	cs, err := client.NewForConfig(&rest.Config{Host: server.URL})
	if err != nil {
		panic(err)
	}
	options := builder.ListOptionsBuilder().
		Clusters("cluster-01").
		Namespaces("kube-system").
		Options()

	pods, err := cs.CoreV1().Pods("").List(context.TODO(), options)
	if err != nil {
		panic(err)
	}

	for _, item := range pods.Items {
		fmt.Printf("Pod info: %v", item)
	}
}

func NewObject(name string) runtime.Object {
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   "kube-system",
			ClusterName: "cluster-01",
		},
	}
}
