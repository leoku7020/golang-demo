package bootkit

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"demo/pkg/logger"
)

type handlerFunc func() error

// ShutdownFunc executes the handlerFunc on app shutdown.
type ShutdownFunc func(handlerFunc)

type registerFunc func(shutdownHandler ShutdownFunc) error

const (
	topPriorityShutdownLevel = -1
)

var (
	// ErrNoRegister indicates no register before running Serve()
	ErrNoRegister = errors.New("no register")
)

var (
	// shutdownLevel is used to prioritize the order of shutdownHandlers.
	// The same level of callback function will be triggered at the same time.
	// Default level is zero.
	shutdownHandlers = map[int][]handlerFunc{}
	shutdownHlrMux   = sync.Mutex{}

	serveHandlers = []registerFunc{}
	servHlrMux    = sync.Mutex{}
	servHlrWg     = sync.WaitGroup{}

	wait = make(chan struct{})
)

type orderedHandler struct {
	level    int
	handlers []handlerFunc
}

// AddShutdownHandler hooks callback function that want to be called on app shutdown.
func AddShutdownHandler(handler handlerFunc, options ...BootOptions) {
	o := applyServOptions(options...)

	shutdownHlrMux.Lock()
	defer shutdownHlrMux.Unlock()

	shutdownHandlers[o.level] = append(shutdownHandlers[o.level], handler)
}

func addShutdownHandler(handler handlerFunc) {
	AddShutdownHandler(handler, withTopPriorityShutdownLevel())
}

func fireShutdownHandlers(ctx context.Context) {
	termSignals := make(chan os.Signal)
	allSignals := make(chan os.Signal)

	// syscall.SIGKILL cannot be trapped by linter. remove it
	signal.Notify(termSignals, syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(allSignals)

	go func() {
		for sig := range allSignals {
			switch sig {
			case syscall.SIGCHLD:
				fallthrough
			case syscall.SIGURG:
				continue
			default:
			}
			logger.Info(fmt.Sprintf("got signal: %v", sig.String()))
		}
	}()

	go func() {
		select {
		case sig := <-termSignals:
			logger.Error(fmt.Sprintf("got terminational signal: %v", sig.String()))
		case <-ctx.Done():
			logger.Error(fmt.Sprintf("terminated system manually: %v", ctx.Err()))
		}

		shutdownHlrMux.Lock()
		defer shutdownHlrMux.Unlock()

		now := time.Now()

		ordered := []orderedHandler{}
		for l, hlrs := range shutdownHandlers {
			ordered = append(ordered, orderedHandler{
				level:    l,
				handlers: hlrs,
			})
		}
		// sort in ascending order
		sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].level < ordered[j].level })

		var wg sync.WaitGroup
		for _, o := range ordered {
			// Registered callbacks at the same level run concurrently.
			// Deliver the last-in handler at the same level as possible.
			for i := len(o.handlers) - 1; i >= 0; i-- {
				wg.Add(1)
				go func(idx int, cb handlerFunc) {
					defer wg.Done()
					if err := cb(); err != nil {
						logger.Error(fmt.Sprintf("shutdown callback function failed: %v", err))
					}
					logger.Debug("finish shutdown callback", logger.WithFields(logger.Fields{"level": o.level, "idx": idx}))
				}(i, o.handlers[i])
			}

			wg.Wait()
		}
		logger.Info(fmt.Sprintf("shutdown callbacks finished within %fs", float64(time.Since(now)/time.Second)))
		close(wait)
	}()
}

// Register registers the callback function launching a server running in a infinite loop. ex: http.ListenAndServe()
func Register(handler registerFunc) {
	servHlrMux.Lock()
	defer servHlrMux.Unlock()

	serveHandlers = append(serveHandlers, handler)
}

// ListenAndServe launches all servers registered above at the same time, blockes until all of them down.
func ListenAndServe(ctx context.Context) error {
	servHlrMux.Lock()
	defer servHlrMux.Unlock()

	if len(serveHandlers) == 0 {
		return ErrNoRegister
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fireShutdownHandlers(ctx)

	errCh := make(chan error, len(serveHandlers))
	for i := range serveHandlers {
		serv := serveHandlers[i]
		servHlrWg.Add(1)
		go func() {
			defer servHlrWg.Done()
			if err := serv(addShutdownHandler); err != nil {
				errCh <- err
			}
		}()
	}

	select {
	case err := <-errCh:
		logger.Ctx(ctx).Error("Serve failed", logger.WithError(err))
		cancel()
		return err
	case <-time.After(5 * time.Second):
		logger.Info("start serving without error within 5 seconds")
	}

	// wait until all servers are down.
	servHlrWg.Wait()

	return nil
}

// Wait until all shutdown handlers added by AddShutdownHandler() finished.
func Wait() {
	<-wait
}
