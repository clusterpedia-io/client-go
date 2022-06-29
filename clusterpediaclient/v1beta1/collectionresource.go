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

package v1beta1

import (
	"context"
	"net/http"
	"time"

	clusterpediav1beta1 "github.com/clusterpedia-io/api/clusterpedia/v1beta1"
	scheme "github.com/clusterpedia-io/client-go/clusterpediaclient/scheme"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type ClusterPediaV1beta1 interface {
	CollectionResource() CollectionResourceInterface
}

type ClusterPediaV1beta1Client struct {
	restClient rest.Interface
}

// New creates a new CoreV1Client for the given RESTClient.
func NewForConfig(c *rest.Config) (*ClusterPediaV1beta1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	httpClient, err := rest.HTTPClientFor(&config)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(&config, httpClient)
}

func NewForConfigAndClient(c *rest.Config, h *http.Client) (*ClusterPediaV1beta1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientForConfigAndClient(&config, h)
	if err != nil {
		return nil, err
	}
	return &ClusterPediaV1beta1Client{client}, nil
}

func setConfigDefaults(config *rest.Config) error {
	gv := clusterpediav1beta1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

func (c *ClusterPediaV1beta1Client) CollectionResource() CollectionResourceInterface {
	return &CollectionResource{c.restClient}
}

type CollectionResourceInterface interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*clusterpediav1beta1.CollectionResource, error)
	List(ctx context.Context, opts metav1.ListOptions) (*clusterpediav1beta1.CollectionResourceList, error)
	Fetch(ctx context.Context, name string, opts metav1.ListOptions, params map[string]string) (*clusterpediav1beta1.CollectionResource, error)
}

type CollectionResource struct {
	client rest.Interface
}

func (c *CollectionResource) Get(ctx context.Context, name string, opts metav1.GetOptions) (result *clusterpediav1beta1.CollectionResource, err error) {
	result = &clusterpediav1beta1.CollectionResource{}
	err = c.client.Get().
		Resource("collectionresources").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

func (c *CollectionResource) List(ctx context.Context, opts metav1.ListOptions) (result *clusterpediav1beta1.CollectionResourceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &clusterpediav1beta1.CollectionResourceList{}
	err = c.client.Get().
		Resource("collectionresources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

func (c *CollectionResource) Fetch(ctx context.Context, name string, opts metav1.ListOptions, params map[string]string) (result *clusterpediav1beta1.CollectionResource, err error) {
	request := c.client.Get().
		Resource("collectionresources").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec)

	for p, v := range params {
		request.Param(p, v)
	}

	result = &clusterpediav1beta1.CollectionResource{}
	request.Do(ctx).Into(result)
	return
}
