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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	pedia "github.com/clusterpedia-io/client-go/client"
	"github.com/clusterpedia-io/client-go/tools/builder"
)

func main() {
	c, err := pedia.Client()
	if err != nil {
		panic(err)
	}

	// build listOptions
	options := builder.ListOptionsBuilder().
		Clusters("cluster-01").
		Namespaces("kube-system").
		Offset(10).Limit(5).
		OrderBy("dsad", false).
		Build()

	pods := &corev1.PodList{}
	err = c.List(context.TODO(), pods, options)
	if err != nil {
		panic(err)
	}

	for _, item := range pods.Items {
		fmt.Printf("Pod info: %v\n", item)
	}

	deploys := &appsv1.DeploymentList{}
	err = c.List(context.TODO(), deploys, options)
	if err != nil {
		panic(err)
	}

	for _, item := range deploys.Items {
		fmt.Printf("Deployment info: %v\n", item)
	}
}
