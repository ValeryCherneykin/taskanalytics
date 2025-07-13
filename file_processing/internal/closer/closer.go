package closer

import (
	"os"
	"os/signal"
	"sync"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"go.uber.org/zap"
)

var globalCloser = New()

func Add(f ...func() error) {
	globalCloser.Add(f...)
}

func Wait() {
	globalCloser.Wait()
}

func CloseAll() {
	globalCloser.CloseAll()
}

type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			s := <-ch
			signal.Stop(ch)
			logger.Info("received shutdown signal", zap.String("signal", s.String()))
			c.CloseAll()
		}()
	}
	return c
}

func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

func (c *Closer) Wait() {
	<-c.done
}

func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		logger.Info("executing cleanup funcs")

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				logger.Error("error during cleanup", zap.Error(err))
			}
		}
		logger.Info("all cleanup funcs executed")
	})
}
