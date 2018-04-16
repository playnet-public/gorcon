package gorcon

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoRcon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoRcon Suite")
}

var _ = Describe("GoRcon", func() {
	BeforeEach(func() {
	})

	Describe("Dummy Test", func() {
		It("runs", func() {
		})
	})
})
