package namespacerules

import (
	"context"
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	fakeAdmission "github.com/shikanon/default-rules/pkg/utils/fake"
	"github.com/tidwall/gjson"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	ctrlAdmission "sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ = Describe("Admission Handler", func() {
	var (
		req     ctrlAdmission.Request
		handler PodAffinityMutatingHandler
		err     error
	)

	BeforeEach(func() {
		// init for env
		namespaceWithAnnotations := corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
				Annotations: map[string]string{
					"shikanon.com/node-affinity": `{"dev-groups": ["rcmd"]}`,
				},
			},
		}
		namespaceNoAnnotations := corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default-no-annotation",
			},
		}

		// init admission request
		podData := `{"metadata": {
			"name": "test",
			"namespace": "default",
			"uid": "e911857d-c318-11e8-bbad-025000000001",
			"creationTimestamp": "2022-04-08T12:20:39Z"
		  },
		  "spec": {
			  "containers": [
				  {
					  "image": "registry.cn-beijing.aliyuncs.com/tensorbytes/busybox",
					  "name": "test"
					}
				]
		  }
		  }`
		req, err = fakeAdmission.GetAdmissionRequest(podData, fakeAdmission.POD_TYPE, admissionv1.Create)
		Expect(err).Should(BeNil())

		decoder, err := ctrlAdmission.NewDecoder(runtime.NewScheme())
		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
			Expect(err).Should(BeNil())
		}

		// init for client and handler
		handler = PodAffinityMutatingHandler{
			Client:  fakeclient.NewClientBuilder().WithObjects(&namespaceWithAnnotations).WithObjects(&namespaceNoAnnotations).Build(),
			decoder: decoder,
		}

	})

	Describe("test namespace affinity", func() {
		It("test default namespace", func() {
			By("use podAffinity handler to handler affinity namespace")

			resp := handler.Handle(context.TODO(), req)
			Expect(len(resp.Patches)).To(Equal(1))
			Expect(resp.Patches[0].Operation).To(Equal("add"))
			Expect(resp.Patches[0].Path).To(Equal("/spec/affinity"))
			respValueByte, err := json.Marshal(resp.Patches[0].Value)
			Expect(err).To(BeNil())
			expectValue := `{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"dev-groups","operator":"In","values":["rcmd"]}]}]}}}`
			Expect(gjson.ParseBytes(respValueByte).Raw).To(Equal(expectValue))
		})

	})
})
