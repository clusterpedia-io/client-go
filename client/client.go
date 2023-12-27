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
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusterv1alpha2 "github.com/clusterpedia-io/api/cluster/v1alpha2"
	"github.com/clusterpedia-io/client-go/tools/transport"
)

const (
	DefaultQPS            float32 = 2000
	DefaultBurst          int     = 2000
	DefaultTimeoutSeconds         = 10
)

var Scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(Scheme))
	utilruntime.Must(clusterv1alpha2.AddToScheme(Scheme))
}

func Client() (client.Client, error) {
	config, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}

	return newClient(config)
}

func ClusterClient(cluster string) (client.Client, error) {
	config, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}

	return newClient(config, cluster)
}

func GetClient(config *rest.Config, cluster ...string) (client.Client, error) {
	return newClient(config, cluster...)
}

func newClient(config *rest.Config, cluster ...string) (client.Client, error) {
	var err error
	if len(cluster) == 1 {
		config, err = ClusterConfigFor(config, cluster[0])
	} else {
		config, err = ConfigFor(config)
	}
	if err != nil {
		return nil, err
	}

	c, err := client.New(config, client.Options{
		Scheme: Scheme,
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}

func NewForConfig(cfg *rest.Config) (kubernetes.Interface, error) {
	config, err := ConfigFor(cfg)
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return kubeClient, nil
}

func NewClusterForConfig(cfg *rest.Config, cluster string) (kubernetes.Interface, error) {
	config, err := ClusterConfigFor(cfg, cluster)
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return kubeClient, nil
}

func ConfigFor(cfg *rest.Config) (*rest.Config, error) {
	configShallowCopy := *cfg
	if err := SetConfigDefaults(&configShallowCopy); err != nil {
		return nil, err
	}

	// wrap a transport to rest client config
	configShallowCopy.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		return transport.NewTransport(configShallowCopy.Host, rt)
	})

	return &configShallowCopy, nil
}

func ClusterConfigFor(cfg *rest.Config, cluster string) (*rest.Config, error) {
	configShallowCopy := *cfg
	if err := SetConfigDefaults(&configShallowCopy); err != nil {
		return nil, err
	}

	// wrap a cluster transport to rest client config
	configShallowCopy.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		return transport.NewTransportForCluster(configShallowCopy.Host, cluster, rt)
	})

	return &configShallowCopy, nil
}

func SetConfigDefaults(config *rest.Config) error {
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
