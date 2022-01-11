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

package builder

import (
	"github.com/clusterpedia-io/client-go/constants"
	"strconv"
	"time"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

type ListOptionsInterface interface {
	Clusters(names ...string) ListOptionsInterface
	Names(names ...string) ListOptionsInterface
	Namespaces(namespaces ...string) ListOptionsInterface
	Size(size int) ListOptionsInterface
	Offset(offset int) ListOptionsInterface
	OrderBy(Order) ListOptionsInterface
	Timeout(timeout time.Duration) ListOptionsInterface
	Options() metaV1.ListOptions
}

type Order struct {
	Field string
	Desc  bool
}

type listOptions metaV1.ListOptions

func ListOptionsBuilder() ListOptionsInterface {
	return &listOptions{}
}

func (opts *listOptions) Clusters(names ...string) ListOptionsInterface {
	selector := parseToSelector(opts.LabelSelector)

	r, _ := labels.NewRequirement(constants.SearchLabelClusters, selection.In, append([]string(nil), names...))
	selector = selector.Add(*r)
	opts.LabelSelector = selector.String()
	return opts
}

func (opts *listOptions) Names(names ...string) ListOptionsInterface {
	selector := parseToSelector(opts.LabelSelector)

	r, _ := labels.NewRequirement(constants.SearchLabelNames, selection.In, append([]string(nil), names...))
	selector = selector.Add(*r)
	opts.LabelSelector = selector.String()
	return opts
}

func (opts *listOptions) Namespaces(names ...string) ListOptionsInterface {
	selector := parseToSelector(opts.LabelSelector)

	r, _ := labels.NewRequirement(constants.SearchLabelNamespaces, selection.In, append([]string(nil), names...))
	selector = selector.Add(*r)
	opts.LabelSelector = selector.String()
	return opts
}

func (opts *listOptions) Size(limit int) ListOptionsInterface {
	selector := parseToSelector(opts.LabelSelector)

	r, _ := labels.NewRequirement(constants.SearchLabelSize, selection.Equals, []string{strconv.Itoa(limit)})
	selector = selector.Add(*r)
	opts.LabelSelector = selector.String()
	return opts
}

func (opts *listOptions) Offset(offset int) ListOptionsInterface {
	selector := parseToSelector(opts.LabelSelector)

	r, _ := labels.NewRequirement(constants.SearchLabelOffset, selection.Equals, []string{strconv.Itoa(offset)})
	selector = selector.Add(*r)
	opts.LabelSelector = selector.String()
	return opts
}

func (opts *listOptions) OrderBy(order Order) ListOptionsInterface {
	selector := parseToSelector(opts.LabelSelector)
	var orderby string
	if order.Desc {
		orderby = order.Field + constants.OrderByDesc
	} else {
		orderby = order.Field
	}

	r, _ := labels.NewRequirement(constants.SearchLabelOffset, selection.Equals, []string{orderby})
	selector = selector.Add(*r)
	opts.LabelSelector = selector.String()
	return opts
}

func (opts *listOptions) Timeout(timeout time.Duration) ListOptionsInterface {
	timeoutSeconds := int64(timeout * time.Second)
	opts.TimeoutSeconds = &timeoutSeconds

	return opts
}

func (opts *listOptions) Options() metaV1.ListOptions {
	return metaV1.ListOptions(*opts)
}

func parseToSelector(labelSelector string) labels.Selector {
	selector := labels.NewSelector()
	if len(labelSelector) > 0 {
		_ = metaV1.Convert_string_To_labels_Selector(&labelSelector, &selector, nil)
	}

	return selector
}
