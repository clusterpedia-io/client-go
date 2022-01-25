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

package client

import (
	"github.com/clusterpedia-io/client-go/constants"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	DefaultQPS   float32 = 2000
	DefaultBurst int     = 2000
)

func ConfigFor(cfg *rest.Config) (*rest.Config, error) {
	configShallowCopy := *cfg

	// reset clusterpedia api path
	setConfigDefaults(&configShallowCopy)
	return &configShallowCopy, nil
}

func ClusterConfigFor(cfg *rest.Config, cluster string) (*rest.Config, error) {
	configShallowCopy, err := ConfigFor(cfg)
	if err != nil {
		return nil, err
	}
	configShallowCopy.Host += constants.ClusterAPIPath + cluster
	return configShallowCopy, nil
}

func NewForConfig(cfg *rest.Config) (kubernetes.Interface, error) {
	clientConfig, err := ConfigFor(cfg)
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	return kubeClient, nil
}

func NewClusterForConfig(cfg *rest.Config, cluster string) (kubernetes.Interface, error) {
	clientConfig, err := ClusterConfigFor(cfg, cluster)
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	return kubeClient, nil
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.Host += constants.ClusterPediaAPIPath
	config.Burst = DefaultBurst
	config.QPS = DefaultQPS

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}
