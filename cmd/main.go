package main

import (
	"context"
	"email-verification-service/internal/pkg/app"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	a := app.New()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	go func() {
		if err := a.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve returned err: %v", err)
		}
	}()

	<-ctx.Done()
	if err := a.Shutdown(context.TODO()); err != nil {
		log.Fatalf("app shutdown returned an err: %v\n", err)
	}
}
