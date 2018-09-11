package watcher_test

import (
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/playnet-public/gorcon/pkg/watcher"
)

var _ = Describe("Process", func() {
	var (
		p *watcher.OSProcess
	)

	BeforeEach(func() {
		p = watcher.NewOSProcess("")
	})

	Describe("SetOut", func() {
		It("does set process output", func() {
			_, w := io.Pipe()
			p.SetOut(w, w)
			Expect(p.Cmd.Stderr).NotTo(BeNil())
			Expect(p.Cmd.Stdout).NotTo(BeNil())
		})
	})
})
