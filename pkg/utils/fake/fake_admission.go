package fake

import (
	"errors"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	ctrlAdmission "sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	POD_TYPE        ObjectType = "Pod"
	DEPLOYMENT_TYPE ObjectType = "Deployment"
)

type ObjectType string

// get admission request for test
func GetAdmissionRequest(jsonObj string, objType ObjectType, op admissionv1.Operation) (req ctrlAdmission.Request, err error) {
	var version string
	switch objType {
	case POD_TYPE, DEPLOYMENT_TYPE:
		version = "v1"
	default:
		err := errors.New("the object type is not support")
		return req, err
	}
	admReqPod := admissionv1.AdmissionRequest{
		Kind: metav1.GroupVersionKind{
			Kind:    string(objType),
			Version: version,
		},
		Operation: op,
		Object: k8sRuntime.RawExtension{
			Raw: []byte(jsonObj),
		},
	}
	req = ctrlAdmission.Request{
		AdmissionRequest: admReqPod,
	}
	return req, err
}
