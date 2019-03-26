package watcher_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/playnet-public/gorcon/pkg/mocks"
	"github.com/playnet-public/gorcon/pkg/watcher"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
)

const debug = false

func TestWatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Watcher Suite")
}

var _ = Describe("Watcher", func() {
	setup := func(name string) (ctx context.Context, w *watcher.Watcher, p *mocks.Process) {
		ctx = context.Background()
		l := log.New("", debug)
		ctx = log.WithLogger(ctx, l.WithFields(zap.String("test", name)))
		w = watcher.NewWatcher(ctx, "")
		p = &mocks.Process{}
		w.Process = p

		return
	}

	Describe("Start", func() {
		It("does not return error", func() {
			ctx, w, _ := setup("Start.does not return error")

			Expect(w.Start(ctx)).To(BeNil())
		})
		It("does set process outputs", func() {
			ctx, w, p := setup("Start.does set process outputs")

			w.Start(ctx)
			time.Sleep(100 * time.Millisecond)
			Expect(p.SetOutCallCount()).To(BeEquivalentTo(1))
		})
		It("does start process", func() {
			ctx, w, p := setup("Start.does start process")

			w.Start(ctx)
			time.Sleep(100 * time.Millisecond)
			Expect(p.RunCallCount()).To(BeEquivalentTo(1))
		})
		It("does end OutputHandlers on closed context", func() {
			ctx, w, p := setup("Start.does end OutputHandlers on closed context")

			// TODO: check for running functions to verify there are no leaks
			p.RunReturns(nil)
			w.Stop(ctx)
			w.Start(ctx)
		})
		It("does exit all functions on ctx close", func() {
			ctx, w, p := setup("Start.does exit all functions on ctx close")
			ctx, close := context.WithCancel(ctx)
			w = watcher.NewWatcher(ctx, "")
			w.Process = p

			w.Start(ctx)
			close()
			<-time.After(10 * time.Millisecond)
			Expect(p.StopCallCount()).To(BeEquivalentTo(1))
		})
	})

	Describe("Stop", func() {
		It("does return nil", func() {
			ctx, w, p := setup("Stop.does return nil")

			Expect(w.Stop(ctx)).To(BeNil())
			Expect(p.StopCallCount()).To(BeEquivalentTo(1))
		})
	})

	Describe("KeepAlive", func() {
		It("does revive the watcher once dead", func() {
			ctx, w, p := setup("KeepAlive.does revive the watcher once dead")

			w.Start(ctx)
			w.KeepAlive(ctx)
			w.Stop(ctx)

			Expect(p.StopCallCount()).To(BeEquivalentTo(1))
		})
		It("does revive the watcher on multiple deaths", func() {
			ctx, w, p := setup("KeepAlive.does revive the watcher on multiple deaths")

			w.Start(ctx)
			w.KeepAlive(ctx)
			w.Stop(ctx)
			Expect(p.StopCallCount()).To(BeEquivalentTo(1))
			time.Sleep(10 * time.Millisecond)
			w.Stop(ctx)
			Expect(p.StopCallCount()).To(BeEquivalentTo(2))

		})
		It("does not revive the wacher on ordered death (Stop)", func() {
			ctx, w, p := setup("KeepAlive.does not revive the wacher on ordered death (Stop)")

			w.Start(ctx)
			w.KeepAlive(ctx)
			w.Stop(ctx)
			Expect(p.StopCallCount()).To(BeEquivalentTo(1))
		})
		It("does restart process on stop", func() {
			ctx, w, p := setup("KeepAlive.does restart process on stop")

			p.RunReturns(nil)
			// p.RunReturnsOnCall(0, nil)
			p.RunReturnsOnCall(0, nil)
			p.RunReturnsOnCall(1, errors.New("test err"))
			p.RunReturnsOnCall(2, nil)

			w.KeepAlive(ctx)
			Expect(w.Start(ctx)).To(BeNil())
			Expect(p.RunCallCount()).To(BeEquivalentTo(1))
		})
	})
})
