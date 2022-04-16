package namespacerules

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/shikanon/default-rules/pkg/utils"
	"github.com/tidwall/gjson"

	corev1 "k8s.io/api/core/v1"
	klogv2 "k8s.io/klog/v2"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
	ctrlAdmission "sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutating-pod-affinity,mutating=true,failurePolicy=Never,groups="",resources=pods,verbs=create;update,sideEffects=none,versions=v1,name=namespace.rule.dev.shikanon.com,admissionReviewVersions={v1,v1beta1}

var (
	NSAnnotation = "shikanon.com/node-affinity"
)

// MutatingHandler handles Component
type PodAffinityMutatingHandler struct {
	Client ctrlClient.Client

	// Decoder decodes objects
	decoder *ctrlAdmission.Decoder
}

// pod affinity mutating rule function, can read the namespace labels to set pod affinity
func (h *PodAffinityMutatingHandler) Handle(ctx context.Context, req ctrlAdmission.Request) ctrlAdmission.Response {
	var pod corev1.Pod
	err := h.decoder.Decode(req, &pod)
	if err != nil || pod.Namespace == "" {
		return ctrlAdmission.Errored(http.StatusBadRequest, err)
	}
	var ns corev1.Namespace
	if err := h.Client.Get(ctx, ctrlClient.ObjectKey{Name: pod.Namespace}, &ns); err != nil {
		klogv2.Errorln(err)
		return ctrlAdmission.Allowed(err.Error())
	}
	klogv2.Info(ns.Labels, ns.Name)
	var updatePod corev1.Pod
	for k, vannotations := range ns.Annotations {
		if k == NSAnnotation {
			updatePod, err = updatePodAffinity(pod, vannotations)
			if err != nil {
				klogv2.Errorln(err)
			}
			break
		}
	}
	podRaw, err := json.Marshal(pod)
	if err != nil {
		klogv2.Errorln(err)
		return ctrlAdmission.Allowed(err.Error())
	}
	updatePodRaw, err := json.Marshal(updatePod)
	if err != nil {
		klogv2.Errorln(err)
		return ctrlAdmission.Allowed(err.Error())
	}
	resp := ctrlAdmission.PatchResponseFromRaw(podRaw, updatePodRaw)
	return resp
}

// InjectDecoder injects the decoder.
func (h *PodAffinityMutatingHandler) InjectDecoder(d *ctrlAdmission.Decoder) error {
	h.decoder = d
	return nil
}

// update pod NodeAffinity by affinity Expression config
func updatePodAffinity(pod corev1.Pod, affinityExpr string) (updatePod corev1.Pod, err error) {
	parseAffinityResult := gjson.Parse(affinityExpr)
	if !parseAffinityResult.IsObject() {
		err = fmt.Errorf("%s is not expect", affinityExpr)
		return
	}

	jsonPatchOprator := utils.NewJsonOperator()

	// parse expression and merge patch the  pod
	parseAffinityResult.ForEach(func(key, value gjson.Result) bool {
		// value type must array
		if value.IsArray() {
			var values = make([]string, 0)
			for _, v := range value.Array() {
				values = append(values, v.String())
			}
			var nodeAffinityTerm = corev1.NodeSelectorTerm{
				MatchExpressions: []corev1.NodeSelectorRequirement{
					{
						Key:      key.String(),
						Operator: corev1.NodeSelectorOpIn,
						Values:   values,
					},
				},
			}

			nodeAffinity := corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: []corev1.NodeSelectorTerm{
							nodeAffinityTerm,
						},
					},
				},
			}

			// jsonpatch doc: http://jsonpatch.com/
			// add to path and value to jsonpatch object
			jsonPatchOprator.Add("/spec/affinity", nodeAffinity)

			return true
		}
		return false
	})

	// Assemble json patch operator string
	patchString, err := jsonPatchOprator.UpdateObject(pod)
	if err != nil {
		return
	}
	klogv2.Info("patchString: ", patchString)

	err = json.Unmarshal([]byte(patchString), &updatePod)
	if err != nil {
		klogv2.Error(err)
		return
	}

	klogv2.Info("updatePod: ", updatePod.String())
	return
}
