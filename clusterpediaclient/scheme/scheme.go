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

package scheme

import (
	clusterpediav1beta1 "github.com/clusterpedia-io/api/clusterpedia/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var Scheme = runtime.NewScheme()
var Codecs = serializer.NewCodecFactory(Scheme)
var ParameterCodec = runtime.NewParameterCodec(Scheme)
var localSchemeBuilder = runtime.NewSchemeBuilder()

func init() {
	localSchemeBuilder.Register(addKnownTypes)
	utilruntime.Must(localSchemeBuilder.AddToScheme(Scheme))
}

// scheme clusterpediav1beta1 miss metav1.ListOption.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(clusterpediav1beta1.SchemeGroupVersion,
		&clusterpediav1beta1.CollectionResource{},
		&clusterpediav1beta1.CollectionResourceList{},
	)
	metav1.AddToGroupVersion(scheme, clusterpediav1beta1.SchemeGroupVersion)
	return nil
}
