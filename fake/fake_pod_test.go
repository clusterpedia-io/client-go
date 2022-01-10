package fake

import (
	"client-go/constants"
	"client-go/tools/builder"
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFakeClient(t *testing.T) {
	cs := NewSimpleClientset()

	// generic pod data
	podDates := []struct {
		name        string
		clusterNmae string
		namespace   string
	}{
		{"pod01", "cluster01", "default"},
		{"pod02", "cluster01", "default"},
		{"pod03", "cluster01", "default"},
		{"pod04", "cluster01", "default"},
		{"pod05", "cluster01", "kube-system"},
		{"pod06", "cluster01", "kube-system"},
		{"pod07", "cluster01", "kube-system"},
		{"pod08", "cluster02", "default"},
		{"pod09", "cluster02", "default"},
		{"pod10", "cluster02", "kube-system"},
		{"pod11", "cluster03", "kube-system"},
		{"pod12", "cluster03", "kube-system"},
	}
	for _, p := range podDates {
		pod := newPod(p.name, p.clusterNmae)
		cs.CoreV1().Pods(p.namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	}

	testCase := []struct {
		opt         metav1.ListOptions
		expectCount int
	}{
		{
			builder.ListOptionsBuilder().Clusters("cluster01").Options(), 7,
		},
		{
			builder.ListOptionsBuilder().Clusters("cluster01").
				Namespaces("kube-system").Options(), 3,
		},
		{
			builder.ListOptionsBuilder().Clusters("cluster01", "cluster02").
				Namespaces("kube-system").Options(), 4,
		},
		{
			builder.ListOptionsBuilder().Clusters("cluster01", "cluster02").
				Namespaces("kube-system", "default").Options(), 10,
		},
		{
			builder.ListOptionsBuilder().Clusters("cluster01", "cluster02").
				Offset(2).Size(4).Options(), 4,
		},
	}

	for _, test := range testCase {
		t.Run("", func(t *testing.T) {
			pods, err := cs.CoreV1().Pods("").List(context.TODO(), test.opt)
			if err != nil {
				t.Error(err)
			}
			if len(pods.Items) != test.expectCount {
				t.Errorf("Unexpect label selector: %s", test.opt.LabelSelector)
			}
		})
	}
}

func newPod(name, clusterName string) *corev1.Pod {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{Kind: "pods", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	if len(clusterName) > 0 {
		pod.Labels = map[string]string{
			constants.SearchLabelClusters: clusterName,
		}
	}
	return pod
}
