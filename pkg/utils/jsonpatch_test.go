package utils

import (
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Json Patch Method Test", func() {
	var (
		pod corev1.Pod
	)

	BeforeEach(func() {
		By("create test pod")
		pod = corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-pod",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "container-1"},
				},
			},
		}
	})

	It("test add operator for image path merge the test pod", func() {
		op := NewJsonOperator()
		op.Add("/spec/containers/0/image", "image name")
		res, err := op.UpdateObject(pod)
		Expect(err).Should(BeNil())
		expectResult := `{"metadata":{"name":"test-pod","creationTimestamp":null},"spec":{"containers":[{"image":"image name","name":"container-1","resources":{}}]},"status":{}}`
		equal := cmp.Diff(res, expectResult)
		GinkgoWriter.Write([]byte(equal))
		Expect(equal).To(BeEmpty())
	})

	It("test add operator for container path  merge the test pod", func() {
		op := NewJsonOperator()
		container := corev1.Container{
			Name: "container-2",
		}
		op.Add("/spec/containers/1", container)
		res, err := op.UpdateObject(pod)
		Expect(err).Should(BeNil())
		expectResult := `{"metadata":{"name":"test-pod","creationTimestamp":null},"spec":{"containers":[{"name":"container-1","resources":{}},{"name":"container-2","resources":{}}]},"status":{}}`
		equal := cmp.Diff(res, expectResult)
		GinkgoWriter.Write([]byte(equal))
		Expect(equal).To(BeEmpty())
	})
})
