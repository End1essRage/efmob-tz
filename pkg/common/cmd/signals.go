package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func Context() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-sigs:
			cancel()
		case <-ctx.Done():
		}
		signal.Stop(sigs)
		close(sigs)
	}()

	return ctx
}
