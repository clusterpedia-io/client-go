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
	"testing"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func TestListOptions(t *testing.T) {
	testCase := []struct {
		opt                 metaV1.ListOptions
		expectLabelSelector string
	}{
		{
			ListOptionsBuilder().Clusters("aaa").Clusters("bbb", "ccc").Options(),
			"search.clusterpedia.io/clusters in (aaa,bbb,ccc)",
		},
		{
			ListOptionsBuilder().Clusters("aaa", "bbb", "ccc").
				Namespaces("ddd").Options(),
			"search.clusterpedia.io/clusters in (aaa,bbb,ccc),search.clusterpedia.io/namespaces=ddd",
		},
		{
			ListOptionsBuilder().Clusters("aaa").
				Namespaces("bbb", "ccc").
				Namespaces("ddd").Options(),
			"search.clusterpedia.io/clusters=aaa,search.clusterpedia.io/namespaces in (bbb,ccc,ddd)",
		},
		{
			ListOptionsBuilder().Clusters("aaa").Clusters("bbbb").
				Namespaces("ccc").Options(),
			"search.clusterpedia.io/clusters in (aaa,bbbb),search.clusterpedia.io/namespaces=ccc",
		},
		{
			ListOptionsBuilder().Clusters("aaa").
				Namespaces("bbb").
				OrderBy("dsad", true).
				OrderBy("basd").Options(),
			"search.clusterpedia.io/clusters=aaa,search.clusterpedia.io/namespaces=bbb,search.clusterpedia.io/orderby in (basd,dsad_desc)",
		},
		{
			ListOptionsBuilder().Clusters("aaa").
				Namespaces("bbb").
				OrderBy("dsad", true).
				RemainingCount().Options(),
			"search.clusterpedia.io/clusters=aaa,search.clusterpedia.io/namespaces=bbb,search.clusterpedia.io/orderby=dsad_desc,search.clusterpedia.io/with-remaining-count=true",
		},
		{
			ListOptionsBuilder().Clusters("aaa").
				Selector(labels.SelectorFromSet(map[string]string{"a": "b"})).Options(),
			"a=b,search.clusterpedia.io/clusters=aaa",
		},
		{
			ListOptionsBuilder().Options(),
			"",
		},
	}

	for _, test := range testCase {
		t.Run("", func(t *testing.T) {
			if test.opt.LabelSelector != test.expectLabelSelector {
				t.Errorf("Unexpect label selector: %s, expect: %s", test.opt.LabelSelector, test.expectLabelSelector)
			}

		})
	}
}
