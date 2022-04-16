package namespacerules

import (
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func RegisterMutatingHandler(mgr manager.Manager) {
	server := mgr.GetWebhookServer()
	server.Register("/mutating-pod-affinity", &webhook.Admission{Handler: &PodAffinityMutatingHandler{Client: mgr.GetClient()}})
}
