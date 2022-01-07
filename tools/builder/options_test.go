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
)

func TestListOptions(t *testing.T) {
	testCase := []struct {
		opt                 metaV1.ListOptions
		expectLabelSelector string
	}{
		{
			ListOptionsBuilder().Clusters("cluster-01").Options(),
			"search.clusterpedia.io/clusters in (cluster-01)",
		},
		{
			ListOptionsBuilder().Clusters("cluster-01", "cluster-02", "ABC").Options(),
			"search.clusterpedia.io/clusters in (ABC,cluster-01,cluster-02)",
		},
		{
			ListOptionsBuilder().Clusters("cluster-01", "cluster-02").
				Namespaces("default", "kube-system").Options(),
			"search.clusterpedia.io/clusters in (cluster-01,cluster-02),search.clusterpedia.io/namespaces in (default,kube-system)",
		},
		{
			ListOptionsBuilder().Clusters("cluster-01").
				Namespaces("kube-system").Offset(0).Size(5).Options(),
			"search.clusterpedia.io/clusters in (cluster-01),search.clusterpedia.io/namespaces in (kube-system),search.clusterpedia.io/offset=0,search.clusterpedia.io/size=5",
		},
		{
			ListOptionsBuilder().Clusters("cluster-01").
				Namespaces("kube-system").
				Offset(10).Size(5).
				OrderBy(Order{"dsad", false}).Options(),
			"search.clusterpedia.io/clusters in (cluster-01),search.clusterpedia.io/namespaces in (kube-system),search.clusterpedia.io/offset=10,search.clusterpedia.io/offset=dsad,search.clusterpedia.io/size=5",
		},
	}
	for _, test := range testCase {
		t.Run("", func(t *testing.T) {
			if test.opt.LabelSelector != test.expectLabelSelector {
				t.Errorf("Unexpect label selector: %s", test.opt.LabelSelector)
			}
		})
	}
}
