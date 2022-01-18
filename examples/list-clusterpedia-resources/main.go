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
	pedia "github.com/clusterpedia-io/client-go/client"
	"github.com/clusterpedia-io/client-go/tools/builder"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	c, err := ClusterpediaClientLoadFromFile()
	if err != nil {
		panic(err)
	}

	// build listOptions
	options := builder.ListOptionsBuilder().
		Clusters("cluster-01").
		Namespaces("kube-system").
		Offset(10).Limit(5).
		OrderBy("dsad", false).
		Options()

	pods, err := c.CoreV1().Pods("").List(context.TODO(), options)
	if err != nil {
		panic(err)
	}

	for _, item := range pods.Items {
		fmt.Printf("Pod info: %v", item)
	}
}

func ClusterpediaClientLoadFromFile() (kubernetes.Interface, error) {
	config, err := clientcmd.LoadFromFile("/path/to/config")
	if err != nil {
		return nil, err
	}
	overrides := clientcmd.ConfigOverrides{Timeout: "10s"}
	clientConfig, err := clientcmd.NewDefaultClientConfig(*config, &overrides).ClientConfig()
	if err != nil {
		return nil, err
	}

	return pedia.NewForConfig(clientConfig)
}
