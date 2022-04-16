package utils

import (
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
	klogv2 "k8s.io/klog/v2"
)

// encapsulation of json patch
type JsonOperator struct {
	Operator []string
}

func NewJsonOperator() JsonOperator {
	return JsonOperator{
		Operator: make([]string, 0),
	}
}

// create add operator
func (o *JsonOperator) Add(path string, obj interface{}) (err error) {
	objJson, err := json.Marshal(obj)
	if err != nil {
		return
	}
	// jsonpatch doc: http://jsonpatch.com/
	patchOp := fmt.Sprintf(`{
		"op": "add",
		"path": "%s",
		"value": %s
	}`, path, string(objJson))
	o.Operator = append(o.Operator, patchOp)
	return
}

func (o *JsonOperator) UpdateObject(v interface{}) (result string, err error) {
	// marshal origin value
	orginValue, err := json.Marshal(v)
	if err != nil {
		klogv2.Error(err)
	}
	if len(o.Operator) == 0 {
		return string(orginValue), nil
	}
	// Assemble json patch operator string
	patchJsonString := "["
	for _, op := range o.Operator {
		patchJsonString = patchJsonString + op
	}
	patchJsonString = patchJsonString + "]"
	patchObject, err := jsonpatch.DecodePatch([]byte(patchJsonString))
	if err != nil {
		err = fmt.Errorf("error for json operator patch decode, %s, %w", patchJsonString, err)
		return
	}

	// apply to it
	patchObj, err := patchObject.Apply(orginValue)
	if err != nil {
		err = fmt.Errorf("error for json operator apply to value: %s, %w", string(orginValue), err)
		return
	}
	return string(patchObj), err
}
