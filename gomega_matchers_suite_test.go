package matchers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGomegaMatchers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GomegaMatchers Suite")
}
