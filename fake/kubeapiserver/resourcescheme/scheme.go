package resourcescheme

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/kubernetes/pkg/api/legacyscheme"

	unstructuredresourcescheme "github.com/clusterpedia-io/client-go/fake/kubeapiserver/resourcescheme/unstructured"
)

var (
	LegacyResourceScheme         = legacyscheme.Scheme
	LegacyResourceCodecs         = legacyscheme.Codecs
	LegacyResourceParameterCodec = legacyscheme.ParameterCodec

	CustomResourceScheme = unstructuredresourcescheme.NewScheme()
	CustomResourceCodecs = unstructured.UnstructuredJSONScheme
)
