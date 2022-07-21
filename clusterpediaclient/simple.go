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

package clusterpediaclient

import (
	"fmt"

	"github.com/clusterpedia-io/client-go/clusterpediaclient/v1beta1"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/flowcontrol"
)

type clusterpediaClient struct {
	pediaClusterClient *v1beta1.ClusterPediaV1beta1Client
}

// AppsV1beta1 retrieves the  PediaClusterV1beta1
func (c *clusterpediaClient) PediaClusterV1beta1() v1beta1.ClusterPediaV1beta1 {
	return c.pediaClusterClient
}

func NewForConfig(cfg *rest.Config) (*clusterpediaClient, error) {
	configShallowCopy := cfg
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}

	var err error
	var cc clusterpediaClient
	cc.pediaClusterClient, err = v1beta1.NewForConfig(configShallowCopy)
	if err != nil {
		return nil, err
	}

	return &cc, nil
}
