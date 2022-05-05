package main

import (
	"context"
	"fmt"
	"log"

	"github.com/clusterpedia-io/client-go/customclient"
	"github.com/clusterpedia-io/client-go/tools/builder"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	config, err := ctrl.GetConfig()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}
	customClient, err := customclient.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to init customClient: %v", err)
	}

	deploys := &appsv1.DeploymentList{}
	options := builder.ListOptionsBuilder().
		Offset(0).Limit(10).
		RemainingCount().
		Options()
	customClient.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).
		Namespace("default").
		List(context.TODO(), options, map[string]string{"clusters": "kpanda-global-cluster"}, deploys)

	for _, item := range deploys.Items {
		fmt.Printf("namespace: %s, name: %s\n", item.Namespace, item.Name)
	}
}
