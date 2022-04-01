package main

import (
	"context"
	"log"

	fake "github.com/clusterpedia-io/client-go/fake/storage"
	"github.com/clusterpedia-io/clusterpedia/pkg/storage"
	internal "github.com/clusterpedia-io/clusterpedia/pkg/apis/clusterpedia"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
)

func main() {
	config := &storage.ResourceStorageConfig{
		Codec: NewCodec(),
	}
	f := fake.NewFakeResourceStorage(config)
	ctx := context.TODO()
	f.Create(ctx, "test", NewObject("test0"))
	f.Create(ctx, "test", NewObject("test1"))
	f.Create(ctx, "test", NewObject("test2"))
	f.Create(ctx, "test", NewObject("test3"))
	f.Create(ctx, "test", NewObject("test4"))
	f.Create(ctx, "test", NewObject("test5"))
	f.Create(ctx, "test", NewObject("test6"))
	log.Println("create")
	pod := &corev1.Pod{}
	err := f.Get(ctx, "test", "test", "test0", pod)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(pod)
	pod.Finalizers = []string{"test111111111", "test2222222222222222"}
	err = f.Update(ctx, "test", pod)
	if err != nil {
		log.Println(err)
		return
	}
	err = f.Get(ctx, "test", "test", "test0", pod)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(pod)
	list := corev1.PodList{}
	opts := &internal.ListOptions{}
	opts.ClusterNames = []string{"test"}
	opts.Namespaces = []string{"test"}
	err = f.List(ctx, &list, opts)
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range list.Items {
		log.Println(v)
	}
}

func NewCodec() runtime.Codec {
	s := runtime.NewScheme()
	metav1.AddToGroupVersion(s, schema.GroupVersion{Group: "", Version: "v1"})
	localSchemeBuilder := runtime.SchemeBuilder{
		corev1.AddToScheme,
		batchv1.AddToScheme,
	}
	localSchemeBuilder.AddToScheme(s)
	return runtimeserializer.NewCodecFactory(s).LegacyCodec(schema.GroupVersion{Group: "", Version: "v1"})
}

func NewObject(name string) runtime.Object {
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   "test",
			ClusterName: "test",
		},
	}
}
