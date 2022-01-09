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
	"client-go/constants"
	"strconv"
	"time"

	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

type ListOptionsInterface interface {
	Clusters(clusters ...string) ListOptionsInterface
	Names(names ...string) ListOptionsInterface
	Namespaces(namespaces ...string) ListOptionsInterface
	Size(size int) ListOptionsInterface
	Offset(offset int) ListOptionsInterface
	OrderBy(field string, desc ...bool) ListOptionsInterface
	Timeout(timeout time.Duration) ListOptionsInterface
	Options() metav1.ListOptions
}

type listOptions metainternal.ListOptions

func ListOptionsBuilder() ListOptionsInterface {
	return &listOptions{
		LabelSelector: NewSelector(),
	}
}

func (opts *listOptions) Clusters(clusters ...string) ListOptionsInterface {
	var r *labels.Requirement
	r, _ = labels.NewRequirement(constants.SearchLabelClusters, selection.In, append([]string(nil), clusters...))
	opts.LabelSelector = opts.LabelSelector.Add(*r)
	return opts
}

func (opts *listOptions) Names(names ...string) ListOptionsInterface {
	r, _ := labels.NewRequirement(constants.SearchLabelNames, selection.In, append([]string(nil), names...))
	opts.LabelSelector = opts.LabelSelector.Add(*r)
	return opts
}

func (opts *listOptions) Namespaces(namespaces ...string) ListOptionsInterface {
	r, _ := labels.NewRequirement(constants.SearchLabelNamespaces, selection.In, append([]string(nil), namespaces...))
	opts.LabelSelector = opts.LabelSelector.Add(*r)
	return opts
}

func (opts *listOptions) Size(limit int) ListOptionsInterface {
	r, _ := labels.NewRequirement(constants.SearchLabelLimit, selection.Equals, []string{strconv.Itoa(limit)})
	opts.LabelSelector = opts.LabelSelector.Add(*r)
	return opts
}

func (opts *listOptions) Offset(offset int) ListOptionsInterface {
	r, _ := labels.NewRequirement(constants.SearchLabelOffset, selection.Equals, []string{strconv.Itoa(offset)})
	opts.LabelSelector = opts.LabelSelector.Add(*r)
	return opts
}

func (opts *listOptions) OrderBy(field string, desc ...bool) ListOptionsInterface {
	if len(desc) > 0 && desc[0] {
		field += constants.OrderByDesc
	}

	r, _ := labels.NewRequirement(constants.SearchLabelOrderBy, selection.In, []string{field})
	opts.LabelSelector = opts.LabelSelector.Add(*r)
	return opts
}

func (opts *listOptions) Timeout(timeout time.Duration) ListOptionsInterface {
	timeoutSeconds := int64(timeout * time.Second)
	opts.TimeoutSeconds = &timeoutSeconds

	return opts
}

func (opts *listOptions) Options() metav1.ListOptions {
	var v1 metav1.ListOptions
	metainternal.Convert_internalversion_ListOptions_To_v1_ListOptions((*metainternal.ListOptions)(opts), &v1, nil)
	return v1
}
