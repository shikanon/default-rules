package namespacerules_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNamespacerules(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Namespacerules Suite")
}
