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
	"strconv"
	"strings"
	"time"

	"github.com/clusterpedia-io/client-go/constants"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
)

type ListOptionsInterface interface {
	Clusters(clusters ...string) ListOptionsInterface
	Names(names ...string) ListOptionsInterface
	FuzzyNames(names ...string) ListOptionsInterface
	Namespaces(namespaces ...string) ListOptionsInterface
	Limit(limit int) ListOptionsInterface
	Offset(offset int) ListOptionsInterface
	OrderBy(field string, desc ...bool) ListOptionsInterface
	Timeout(timeout time.Duration) ListOptionsInterface
	TimeoutSeconds(timeout int64) ListOptionsInterface
	RemainingCount() ListOptionsInterface
	OwnerUID(uid string) ListOptionsInterface
	OwnerName(name string) ListOptionsInterface
	OwnerSeniority(ownerSeniority int) ListOptionsInterface
	OwnerGroupResource(groupResource schema.GroupResource) ListOptionsInterface
	LabelSelector(field string, values []string) ListOptionsInterface
	Selector(ls labels.Selector) ListOptionsInterface
	FieldSelector(field string, values []string) ListOptionsInterface
	Options() metav1.ListOptions
	Build() *client.ListOptions
}

type listOptions struct {
	options       metav1.ListOptions
	labels        map[string][]string
	labelSelector labels.Selector
	fieldSelector map[string][]string
}

func ListOptionsBuilder() ListOptionsInterface {
	return &listOptions{
		options:       metav1.ListOptions{},
		labels:        make(map[string][]string),
		fieldSelector: make(map[string][]string),
	}
}

func (opts *listOptions) Clusters(clusters ...string) ListOptionsInterface {
	if len(clusters) > 0 {
		opts.labels[constants.SearchLabelClusters] =
			append(opts.labels[constants.SearchLabelClusters], clusters...)
	}
	return opts
}

func (opts *listOptions) Names(names ...string) ListOptionsInterface {
	if len(names) > 0 {
		opts.labels[constants.SearchLabelNames] =
			append(opts.labels[constants.SearchLabelNames], names...)
	}
	return opts
}

func (opts *listOptions) FuzzyNames(names ...string) ListOptionsInterface {
	if len(names) > 0 {
		opts.labels[constants.SearchLabelFuzzyName] =
			append(opts.labels[constants.SearchLabelFuzzyName], names...)
	}
	return opts
}

func (opts *listOptions) OwnerUID(uid string) ListOptionsInterface {
	uid = strings.TrimSpace(uid)
	if len(uid) > 0 {
		opts.labels[constants.SearchLabelOwnerUID] = []string{uid}
	}
	return opts
}

func (opts *listOptions) OwnerName(name string) ListOptionsInterface {
	name = strings.TrimSpace(name)
	if len(name) > 0 {
		opts.labels[constants.SearchLabelOwnerName] = []string{name}
	}
	return opts
}

func (opts *listOptions) OwnerSeniority(ownerSeniority int) ListOptionsInterface {
	if ownerSeniority > 0 {
		opts.labels[constants.SearchLabelOwnerSeniority] = []string{strconv.Itoa(ownerSeniority)}
	}
	return opts
}

func (opts *listOptions) OwnerGroupResource(groupResource schema.GroupResource) ListOptionsInterface {
	gr := strings.TrimSpace(groupResource.String())
	if len(gr) > 0 {
		opts.labels[constants.SearchLabelOwnerGroupResource] = []string{gr}
	}
	return opts
}

func (opts *listOptions) Namespaces(namespaces ...string) ListOptionsInterface {
	if len(namespaces) > 0 {
		opts.labels[constants.SearchLabelNamespaces] =
			append(opts.labels[constants.SearchLabelNamespaces], namespaces...)
	}
	return opts
}

func (opts *listOptions) Limit(limit int) ListOptionsInterface {
	if limit > 0 {
		opts.options.Limit = int64(limit)
	}
	return opts
}

func (opts *listOptions) Offset(offset int) ListOptionsInterface {
	if offset >= 0 {
		opts.options.Continue = strconv.Itoa(offset)
	}
	return opts
}

func (opts *listOptions) OrderBy(field string, desc ...bool) ListOptionsInterface {
	var orderby string
	if len(field) > 0 {
		orderby = field

		if len(desc) > 0 && desc[len(desc)-1] {
			orderby += constants.OrderByDesc
		}

		opts.labels[constants.SearchLabelOrderBy] =
			append(opts.labels[constants.SearchLabelOrderBy], orderby)
	}
	return opts
}

func (opts *listOptions) Timeout(timeout time.Duration) ListOptionsInterface {
	if timeout > 0 {
		timeoutSeconds := int64(timeout * time.Second)
		opts.options.TimeoutSeconds = &timeoutSeconds
	}
	return opts
}

func (opts *listOptions) TimeoutSeconds(timeout int64) ListOptionsInterface {
	if timeout > 0 {
		opts.options.TimeoutSeconds = &timeout
	}
	return opts
}

func (opts *listOptions) RemainingCount() ListOptionsInterface {
	opts.labels[constants.SearchLabelWithRemainingCount] =
		append(opts.labels[constants.SearchLabelWithRemainingCount], strconv.FormatBool(true))

	return opts
}

func (opts *listOptions) LabelSelector(field string, values []string) ListOptionsInterface {
	opts.labels[field] =
		append(opts.labels[field], values...)
	return opts
}

func (opts *listOptions) Selector(ls labels.Selector) ListOptionsInterface {
	opts.labelSelector = ls
	return opts
}

func (opts *listOptions) FieldSelector(field string, values []string) ListOptionsInterface {
	opts.fieldSelector[field] =
		append(opts.fieldSelector[field], values...)
	return opts
}

func (opts *listOptions) Options() metav1.ListOptions {
	ls := labels.Everything()
	if opts.labelSelector != nil {
		ls = opts.labelSelector
	}
	for label, values := range opts.labels {
		var op selection.Operator
		if len(values) > 1 {
			op = selection.In
		} else {
			op = selection.Equals
		}
		r, _ := labels.NewRequirement(label, op, append([]string(nil), values...))
		ls = ls.Add(*r)
	}
	opts.options.LabelSelector = ls.String()

	if len(opts.fieldSelector) == 0 {
		opts.options.FieldSelector = fields.Everything().String()
	} else {
		requirements := make([]labels.Requirement, 0, len(opts.fieldSelector))
		for label, values := range opts.fieldSelector {
			var op selection.Operator
			if len(values) > 1 {
				op = selection.In
			} else {
				op = selection.Equals
			}

			r, _ := labels.NewRequirement(label, op, append([]string(nil), values...))
			requirements = append(requirements, *r)
		}
		selector := labels.NewSelector()
		selector = selector.Add(requirements...)
		opts.options.FieldSelector = selector.String()
	}
	return opts.options
}

func (opts *listOptions) Build() *client.ListOptions {
	opt := opts.Options()

	return &client.ListOptions{Raw: &opt, Limit: opt.Limit, Continue: opt.Continue}
}
