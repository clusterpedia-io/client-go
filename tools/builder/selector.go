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
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

type simpleSelector []labels.Requirement

func NewSelector() labels.Selector {
	return simpleSelector(nil)
}

func (s simpleSelector) Matches(_ labels.Labels) bool                               { return false }
func (s simpleSelector) Empty() bool                                                { return false }
func (s simpleSelector) Requirements() (labels.Requirements, bool)                  { return nil, false }
func (s simpleSelector) DeepCopySelector() labels.Selector                          { return s }
func (s simpleSelector) RequiresExactMatch(label string) (value string, found bool) { return "", false }

func (s simpleSelector) Add(reqs ...labels.Requirement) labels.Selector {
	for ix, requirement := range reqs {
		if val, found := s.exists(requirement.Key()); found {
			s = s.remove(requirement.Key())

			values := requirement.Values().List()
			if requirement.Operator() == selection.In {
				values = append(values, val...)
				nr, _ := labels.NewRequirement(requirement.Key(), requirement.Operator(),
					append([]string(nil), values...))

				reqs[ix] = *nr
			}
		}
	}

	ret := make(simpleSelector, 0, len(s)+len(reqs))
	ret = append(ret, s...)
	ret = append(ret, reqs...)
	sort.Sort(labels.ByKey(ret))
	return ret
}

func (s simpleSelector) String() string {
	var reqs []string
	for ix := range s {
		reqs = append(reqs, s[ix].String())
	}
	return strings.Join(reqs, ",")
}

func (s simpleSelector) exists(label string) (value []string, found bool) {
	for _, v := range s {
		if v.Key() == label {
			return v.Values().List(), true
		}
	}
	return nil, false
}

func (s simpleSelector) remove(key string) simpleSelector {
	for i, v := range s {
		if v.Key() == key {
			s = append(s[:i], s[i+1:]...)
			return s
		}
	}
	return s
}
