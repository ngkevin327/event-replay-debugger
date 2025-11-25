package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/replay/platform/apps/notifications/internal/config"
	"github.com/replay/platform/apps/notifications/internal/webhook"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sender := &webhook.Sender{
		Timeout: cfg.WebhookTimeout,
		Retries: cfg.MaxRetries,
	}
	log.Printf("notifications worker started (bus=%s)", cfg.EventBusURL)

	for {
		select {
		case <-ctx.Done():
			log.Println("notifications worker stopped")
			return
		default:
		}
		events, ok := pollEvents(ctx, cfg.EventBusURL)
		if !ok {
			continue
		}
		select {
		case <-ctx.Done():
			return
		case event, open := <-events:
			if !open {
				continue
			}
			if err := handleEvent(ctx, sender, event); err != nil {
				log.Printf("event %s: %v", event.Type, err)
			}
		}
	}
}

type domainEvent struct {
	Type       string
	WebhookURL string
	Payload    map[string]any
}

func handleEvent(ctx context.Context, sender *webhook.Sender, ev domainEvent) error {
	if ev.WebhookURL == "" {
		return nil
	}
	body, err := json.Marshal(ev.Payload)
	if err != nil {
		return err
	}
	return sender.Deliver(ctx, ev.WebhookURL, body)
}

func pollEvents(ctx context.Context, busURL string) (<-chan domainEvent, bool) {
	ch := make(chan domainEvent, 1)
	_ = busURL
	select {
	case <-ctx.Done():
		return nil, false
	default:
		return ch, true
	}
}
