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
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/clusterpedia-io/client-go/constants"
	clusterv1alpha2 "github.com/clusterpedia-io/client-go/pkg/apis/cluster/v1alpha2"
)

const (
	DefaultQPS            float32 = 2000
	DefaultBurst          int     = 2000
	DefaultTimeoutSeconds         = 10
)

func Client() (client.Client, error) {
	restConfig, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}

	return newClient(restConfig)
}

func ClusterClient(cluster string) (client.Client, error) {
	restConfig, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}

	return newClient(restConfig, cluster)
}

func GetClient(restConfig *rest.Config, cluster ...string) (client.Client, error) {
	return newClient(restConfig, cluster...)
}

func newClient(restConfig *rest.Config, cluster ...string) (client.Client, error) {
	var err error

	if len(cluster) == 1 {
		restConfig, err = ClusterConfigFor(restConfig, cluster[0])
	} else {
		restConfig, err = ConfigFor(restConfig)
	}
	if err != nil {
		return nil, err
	}

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(clusterv1alpha2.AddToScheme(scheme))

	c, err := client.New(restConfig, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}

func ConfigFor(cfg *rest.Config) (*rest.Config, error) {
	configShallowCopy := *cfg

	// reset clusterpedia api path
	if err := setConfigDefaults(&configShallowCopy); err != nil {
		return nil, err
	}

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
	config.Host += constants.ClusterPediaAPIPath
	if config.Timeout == 0 {
		config.Timeout = DefaultTimeoutSeconds * time.Second
	}
	if config.Burst == 0 {
		config.Burst = DefaultBurst
	}
	if config.QPS == 0 {
		config.QPS = DefaultQPS
	}
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}
