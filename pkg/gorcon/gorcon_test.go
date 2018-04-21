package gorcon_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoRcon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoRcon Suite")
}
