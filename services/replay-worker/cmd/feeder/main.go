package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/replay/platform/services/replay-worker/internal/feeder"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	f := feeder.NewFeeder()
	log.Println("replay feeder starting")
	if err := f.RunLoop(ctx); err != nil && err != context.Canceled {
		log.Fatal(err)
	}
}
