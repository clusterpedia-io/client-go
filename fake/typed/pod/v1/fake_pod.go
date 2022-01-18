package fake

import (
	"context"
	"github.com/clusterpedia-io/client-go/constants"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	testing "k8s.io/client-go/testing"
)

type PodsGetter interface {
	Pods(namespace string) PodInterface
}

type PodInterface interface {
	Create(ctx context.Context, pod *corev1.Pod, opts metav1.CreateOptions) (*corev1.Pod, error)
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*corev1.Pod, error)
	List(ctx context.Context, opts metav1.ListOptions) (*corev1.PodList, error)
}

// FakeClusters implements PodInterface
type FakePods struct {
	Fake *FakeCoreV1
	ns   string
}

// GroupName is the group name use in this package
const GroupName = ""

var podsResource = schema.GroupVersionResource{Group: GroupName, Version: "v1", Resource: "pods"}
var podsKind = schema.GroupVersionKind{Group: GroupName, Version: "v1", Kind: "Pod"}

// Get takes name of the pod, and returns the corresponding pod object, and an error if there is any.
func (c *FakePods) Get(ctx context.Context, name string, options metav1.GetOptions) (result *corev1.Pod, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(podsResource, c.ns, name), &corev1.Pod{})
	if obj == nil {
		return nil, err
	}
	return obj.(*corev1.Pod), err
}

// List takes label and field selectors, and returns the list of pod that match those selectors.
func (c *FakePods) List(ctx context.Context, opts metav1.ListOptions) (result *corev1.PodList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(podsResource, podsKind, c.ns, opts), &corev1.Pod{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}

	list := &corev1.PodList{ListMeta: obj.(*corev1.PodList).ListMeta}

	for _, item := range obj.(*corev1.PodList).Items {
		pl := make(map[string]string)
		// handle name and cluster name
		for k, v := range item.Labels {
			pl[k] = v
		}
		pl[constants.SearchLabelNames] = item.Name
		pl[constants.SearchLabelNamespaces] = item.Namespace

		if label.Matches(labels.Set(pl)) {
			list.Items = append(list.Items, item)
		}
	}

	// handle offSize and limit
	offset, _ := strconv.Atoi(opts.Continue)
	limt := int(opts.Limit)
	if offset <= len(list.Items) {
		if offset+limt > len(list.Items) {
			list.Items = list.Items[offset:]
		} else if limt > 0 {
			list.Items = list.Items[offset : offset+limt]
		}
	} else {
		list.Items = make([]corev1.Pod, 0)
	}

	return list, err
}

// Create takes the representation of a pod and creates it.  Returns the server's representation of the pod, and an error, if there is any.
func (c *FakePods) Create(ctx context.Context, pod *corev1.Pod, opts metav1.CreateOptions) (result *corev1.Pod, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(podsResource, c.ns, pod), &corev1.Pod{})
	if obj == nil {
		return nil, err
	}
	return obj.(*corev1.Pod), err
}
