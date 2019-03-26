package watcher

import (
	"bufio"
	"context"
	"errors"
	"io"
	"time"

	"github.com/playnet-public/gorcon/pkg/event"

	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
)

// Watcher is responsible for starting and keeping a process alive
type Watcher struct {
	Process   Process
	closeFunc func()
	close     chan error

	*event.Broker
	events chan event.Event

	StopTimeout time.Duration
}

// ErrStopEvent is being sent to the process when the watcher is being ordered to stop
var ErrStopEvent = errors.New("received external stop event")

// NewWatcher responsible for starting and keeping a process alive, restarting if necessary
func NewWatcher(ctx context.Context, path string, args ...string) *Watcher {
	w := &Watcher{
		Process:     nil,
		close:       make(chan error),
		events:      make(chan event.Event),
		StopTimeout: 5 * time.Second,
	}

	return w
}

// Start the underlying process and handle events
func (w *Watcher) Start(ctx context.Context) error {
	ctx, w.closeFunc = context.WithCancel(ctx)

	go func() {
		log.From(ctx).Debug("waiting for ctx to close", zap.String("span", "Watcher.Start"))
		<-ctx.Done()
		log.From(ctx).Debug("handling ctx close", zap.String("span", "Watcher.Start"))
		go func() { w.close <- ctx.Err() }()
	}()

	rerr, stderr := io.Pipe()
	rout, stdout := io.Pipe()

	go w.OutputHandler(ctx, rerr, "StdErr")
	go func() {
		log.From(ctx).Debug("waiting for ctx to close", zap.String("span", "OutputHandler.StdErr"))
		<-ctx.Done()
		log.From(ctx).Debug("handling ctx close", zap.String("span", "OutputHandler.StdErr"))
		stderr.CloseWithError(ctx.Err())
	}()
	go w.OutputHandler(ctx, rout, "StdOut")
	go func() {
		log.From(ctx).Debug("waiting for ctx to close", zap.String("span", "OutputHandler.StdOut"))
		<-ctx.Done()
		log.From(ctx).Debug("handling ctx close", zap.String("span", "OutputHandler.StdOut"))
		stdout.CloseWithError(ctx.Err())
	}()

	go func() {
		w.Broker = event.NewBroker(ctx, w.events)
		log.From(ctx).Debug("running broker")
		if err := w.Broker.Run(ctx); err != nil {
			log.From(ctx).Error("running broker", zap.Error(err))
			w.close <- err
		}
	}()

	w.Process.SetOut(stderr, stdout)
	go func() {
		log.From(ctx).Debug("waiting for ctx to close", zap.String("span", "Process.Stop"))
		<-ctx.Done()
		log.From(ctx).Debug("stopping process")
		if err := w.Process.Stop(); err != nil {
			log.From(ctx).Error("stopping process", zap.Error(err))
		}
	}()

	log.From(ctx).Debug("running process")
	if err := w.Process.Run(); err != nil {
		log.From(ctx).Error("running process", zap.Error(err))
		w.close <- err
		return err
	}

	return ctx.Err()
}

// Stop the underlying process and all event handling routines
func (w *Watcher) Stop(ctx context.Context) error {
	return w.Process.Stop()
}

// KeepAlive starts a go routine responsible for reviving the process once it dies
func (w *Watcher) KeepAlive(ctx context.Context) {
	go func() {
		if err := <-w.close; err != ErrStopEvent {
			log.From(ctx).Info("handling close event", zap.Error(err))
			w.KeepAlive(ctx)
			go func() {
				log.From(ctx).Debug("running process")
				if err := w.Process.Run(); err != nil {
					log.From(ctx).Error("running process", zap.Error(err))
					w.close <- err
				}
			}()
		}
	}()
}

// OutputHandler returns a function reading from io.Reader and creating events
func (w *Watcher) OutputHandler(ctx context.Context, r io.Reader, eventType string) func() error {
	return func() error {
		scn := bufio.NewScanner(r)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if scn.Scan() {
					w.events <- &Event{
						timestamp: time.Now(),
						kind:      eventType,
						payload:   scn.Text(),
					}
					continue
				}
				return errors.New("end of stream")
			}
		}
	}
}
