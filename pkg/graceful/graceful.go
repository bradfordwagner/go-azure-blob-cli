package graceful

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

// New - creates a new graceful termination
func New() (ctx context.Context, cancel context.CancelFunc, errChan chan error) {
	// init
	errChan = make(chan error)
	ctx, cancel = signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// await a sigterm
	go func() {
		defer func() {
			cancel()
			close(errChan)
		}()

		<-ctx.Done()
	}()

	return
}

// AwaitLogError - reusable log error
func AwaitLogError(errChan chan error) {
	err, ok := <-errChan
	if ok {
		logrus.WithError(err).Error("failed execution")
	} else {
		logrus.Trace("canceled")
	}
}
