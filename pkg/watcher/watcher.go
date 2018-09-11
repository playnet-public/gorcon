package watcher

import (
	"bufio"
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"gopkg.in/tomb.v2"
)

// Watcher is responsible for starting and keeping a process alive
type Watcher struct {
	Process Process

	Tomb  *tomb.Tomb
	TLock sync.RWMutex

	subscriptions      []chan *Event
	subscriptionsMutex sync.RWMutex

	StopTimeout time.Duration
}

// ErrStopEvent is being sent to the process when the watcher is being ordered to stop
var ErrStopEvent = errors.New("received external stop event")

// NewWatcher responsible for starting and keeping a process alive, restarting if necessary
func NewWatcher(ctx context.Context, path string, args ...string) *Watcher {
	w := &Watcher{
		Process:       nil,
		subscriptions: []chan *Event{},
		StopTimeout:   5 * time.Second,
	}
	w.TLock.Lock()
	defer w.TLock.Unlock()
	w.Tomb, _ = tomb.WithContext(ctx)
	return w
}

// Start the underlying process and handle events
func (w *Watcher) Start(ctx context.Context) error {
	w.TLock.Lock()
	if !w.Tomb.Alive() {
		w.Tomb, _ = tomb.WithContext(ctx)
	}
	w.TLock.Unlock()
	rerr, stderr := io.Pipe()
	rout, stdout := io.Pipe()

	w.TLock.RLock()
	defer w.TLock.RUnlock()

	w.Tomb.Go(w.OutputHandler(rerr, "StdErr"))
	w.Tomb.Go(func() error {
		w.TLock.RLock()
		defer w.TLock.RUnlock()
		<-w.Tomb.Dying()
		stderr.CloseWithError(w.Tomb.Err())
		return nil
	})
	w.Tomb.Go(w.OutputHandler(rout, "StdOut"))
	w.Tomb.Go(func() error {
		w.TLock.RLock()
		defer w.TLock.RUnlock()
		<-w.Tomb.Dying()
		stdout.CloseWithError(w.Tomb.Err())
		return nil
	})

	w.Process.SetOut(stderr, stdout)
	w.Tomb.Go(w.Process.Run)
	w.Tomb.Go(func() error {
		w.TLock.RLock()
		defer w.TLock.RUnlock()
		<-w.Tomb.Dying()
		return w.Process.Stop()
	})

	return nil
}

// Stop the underlying process and all event handling routines
func (w *Watcher) Stop(ctx context.Context) error {
	w.TLock.RLock()
	defer w.TLock.RUnlock()
	w.Tomb.Kill(ErrStopEvent)
	select {
	case <-w.Tomb.Dead():
		return nil
	case <-time.After(w.StopTimeout):
		return errors.New("timeout stopping watcher")
	}
}

// KeepAlive starts a go routine responsible for reviving the process once it dies
func (w *Watcher) KeepAlive(ctx context.Context) {
	go func() {
		w.TLock.RLock()
		defer w.TLock.RUnlock()
		<-w.Tomb.Dead()
		if w.Tomb.Err() != ErrStopEvent {
			go w.Start(ctx)
			w.KeepAlive(ctx)
		}
	}()
}

// Subscribe for events on the process until the context ends
func (w *Watcher) Subscribe(ctx context.Context, to chan *Event) {
	w.subscriptionsMutex.Lock()
	defer w.subscriptionsMutex.Unlock()
	w.subscriptions = append(w.subscriptions, to)
	go func() {
		<-ctx.Done()
		w.subscriptionsMutex.Lock()
		defer w.subscriptionsMutex.Unlock()
		for i, s := range w.subscriptions {
			if s == to {
				w.subscriptions[i] = w.subscriptions[len(w.subscriptions)-1]
				w.subscriptions[len(w.subscriptions)-1] = nil
				w.subscriptions = w.subscriptions[:len(w.subscriptions)-1]
			}
		}
	}()
}

// OutputHandler returns a function reading from io.Reader and creating events
func (w *Watcher) OutputHandler(r io.Reader, eventType string) func() error {
	return func() error {
		scn := bufio.NewScanner(r)
		for scn.Scan() {
			w.publishEvent(&Event{
				Timestamp: time.Now(),
				Type:      eventType,
				Payload:   scn.Text(),
			})
			continue
		}
		return errors.New("end of stream")
	}
}

func (w *Watcher) publishEvent(e *Event) {
	w.subscriptionsMutex.RLock()
	defer w.subscriptionsMutex.RUnlock()
	for _, l := range w.subscriptions {
		go func(l chan *Event) { l <- e }(l)
	}
}
