package watcher_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/playnet-public/gorcon/pkg/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/playnet-public/gorcon/pkg/watcher"
)

func TestWatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Watcher Suite")
}

var _ = Describe("Watcher", func() {
	var (
		w   *watcher.Watcher
		ctx context.Context
		p   *mocks.Process
	)

	BeforeEach(func() {
		ctx = context.Background()
		w = watcher.NewWatcher(ctx, "")
		p = &mocks.Process{}
		w.Process = p
	})

	Describe("Start", func() {
		BeforeEach(func() {
			w = watcher.NewWatcher(ctx, "")
			p = &mocks.Process{}
			w.Process = p
		})
		It("does not return error", func() {
			Expect(w.Start(ctx)).To(BeNil())
		})
		It("does set process outputs", func() {
			w.Start(ctx)
			time.Sleep(100 * time.Millisecond)
			Expect(p.SetOutCallCount()).To(BeEquivalentTo(1))
		})
		It("does start process", func() {
			w.Start(ctx)
			time.Sleep(100 * time.Millisecond)
			Expect(p.RunCallCount()).To(BeEquivalentTo(1))
		})
		It("does end OutputHandlers on tomb dying", func() {
			p.RunReturns(nil)
			w.Start(ctx)
			time.Sleep(100 * time.Millisecond)
			w.TLock.RLock()
			defer w.TLock.RUnlock()
			w.Tomb.Kill(nil)
			select {
			case <-time.After(1 * time.Second):
				Fail("did not exit")
			case <-w.Tomb.Dead():
				return
			}
		})
	})

	Describe("Stop", func() {
		BeforeEach(func() {
			w = watcher.NewWatcher(ctx, "")
			p = &mocks.Process{}
			w.Process = p
		})
		It("does return nil on successful stop", func() {
			w.TLock.RLock()
			defer w.TLock.RUnlock()
			w.Tomb.Go(func() error {
				w.TLock.RLock()
				defer w.TLock.RUnlock()
				<-w.Tomb.Dying()
				return nil
			})
			Expect(w.Stop(ctx)).To(BeNil())
		})
		It("does return error on kill timeout", func() {
			w := watcher.NewWatcher(ctx, "")
			p := &mocks.Process{}
			w.Process = p
			w.TLock.RLock()
			defer w.TLock.RUnlock()
			w.Tomb.Go(func() error {
				w.TLock.RLock()
				defer w.TLock.RUnlock()
				<-w.Tomb.Dying()
				time.Sleep(100 * time.Millisecond)
				return nil
			})
			w.StopTimeout = 1 * time.Millisecond
			Expect(w.Stop(ctx)).NotTo(BeNil())
		})
	})

	Describe("KeepAlive", func() {
		It("does revive the watcher once dead", func() {
			w.Start(ctx)
			w.KeepAlive(ctx)
			w.TLock.RLock()
			w.Tomb.Kill(errors.New("some test crash"))
			<-w.Tomb.Dead()
			w.TLock.RUnlock()
			time.Sleep(100 * time.Millisecond)
			w.TLock.RLock()
			defer w.TLock.RUnlock()
			Expect(w.Tomb.Alive()).To(BeTrue())
		})
		It("does revive the watcher on multiple deaths", func() {
			w.Start(ctx)
			w.KeepAlive(ctx)
			w.TLock.RLock()
			w.Tomb.Kill(errors.New("some test crash"))
			<-w.Tomb.Dead()
			w.TLock.RUnlock()
			time.Sleep(100 * time.Millisecond)
			w.TLock.RLock()
			w.Tomb.Kill(errors.New("some test crash"))
			w.TLock.RUnlock()
			time.Sleep(100 * time.Millisecond)
			w.TLock.RLock()
			defer w.TLock.RUnlock()
			Expect(w.Tomb.Alive()).To(BeTrue())
		})
		It("does not revive the wacher on ordered death (Stop)", func() {
			w.Start(ctx)
			w.KeepAlive(ctx)
			w.Stop(ctx)
			w.TLock.RLock()
			<-w.Tomb.Dead()
			w.TLock.RUnlock()
			time.Sleep(100 * time.Millisecond)
			w.TLock.RLock()
			defer w.TLock.RUnlock()
			Expect(w.Tomb.Alive()).To(BeFalse())
		})
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
			wr.Close()
			Expect(w.OutputHandler(re, "test")()).NotTo(BeNil())
		})
		It("does send event to subscriptions for each line", func() {
			w.TLock.RLock()
			defer w.TLock.RUnlock()
			w.Tomb.Go(w.OutputHandler(re, "test"))
			c := make(chan *watcher.Event)
			w.Subscribe(ctx, c)
			time.Sleep(500 * time.Millisecond)
			wr.Write([]byte(fmt.Sprintln("testLine1")))
			ev := <-c
			Expect(ev.Type).To(BeEquivalentTo("test"))
		})
		It("does send event to subscriptions for each line", func() {
			w.TLock.RLock()
			defer w.TLock.RUnlock()
			w.Tomb.Go(w.OutputHandler(re, "test"))
			c := make(chan *watcher.Event)
			w.Subscribe(ctx, c)
			time.Sleep(500 * time.Millisecond)
			wr.Write([]byte(fmt.Sprintln("testLine1")))
			ev := <-c
			Expect(ev.Type).To(BeEquivalentTo("test"))
		})
	})
})
