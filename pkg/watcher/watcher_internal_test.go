package watcher

import (
	"context"
	"fmt"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Watcher", func() {
	// var (
	// 	w   *Watcher
	// 	ctx context.Context
	// )

	// BeforeEach(func() {
	// 	ctx = context.Background()
	// 	w = &Watcher{}
	// })

	setup := func() (ctx context.Context, w *Watcher) {
		ctx = context.Background()
		w = NewWatcher(ctx, "")

		return
	}

	Describe("Subscribe", func() {
		// BeforeEach(func() {
		// 	w = &Watcher{
		// 		events: make(chan event.Event),
		// 	}
		// 	ctx = context.Background()
		// })
	})

	Describe("OutputHandler", func() {
		var (
			re *io.PipeReader
			wr *io.PipeWriter
		)
		BeforeEach(func() {
			re, wr = io.Pipe()
		})
		It("does return error on end of stream", func() {
			ctx, w := setup()

			wr.Close()
			Expect(w.OutputHandler(ctx, re, "test")()).NotTo(BeNil())
		})
		It("does send event to subscriptions for each line", func() {
			ctx, w := setup()

			go w.OutputHandler(ctx, re, "test")()
			wr.Write([]byte(fmt.Sprintln("testLine1")))
			ev := <-w.events
			Expect(ev.Kind()).To(BeEquivalentTo("test"))
		})
		It("does send event to subscriptions for each line", func() {
			ctx, w := setup()

			go w.OutputHandler(ctx, re, "test")()
			wr.Write([]byte(fmt.Sprintln("testLine1")))
			ev := <-w.events
			Expect(ev.Kind()).To(BeEquivalentTo("test"))
		})
	})
})
