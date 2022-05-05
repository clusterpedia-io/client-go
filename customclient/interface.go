package customclient

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Interface interface {
	Resource(resource schema.GroupVersionResource) NamespaceableResourceInterface
}

type ResourceInterface interface {
	List(ctx context.Context, opts metav1.ListOptions, params map[string]string, obj runtime.Object) error
}

type NamespaceableResourceInterface interface {
	Namespace(string) ResourceInterface
	ResourceInterface
}
