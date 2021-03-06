/*
Copyright 2021 clusterpedia Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
