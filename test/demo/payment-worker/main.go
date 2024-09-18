package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	agentgo "github.com/replay/platform/packages/agent-go"
	"github.com/replay/platform/packages/agent-go/kafka"
	"github.com/replay/platform/packages/shared-go/event"
)

func main() {
	brokers := envOr("KAFKA_BROKERS", "localhost:9092")
	topic := envOr("KAFKA_TOPIC", "payments.settlement")

	cfg := sarama.NewConfig()
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup([]string{brokers}, "payment-settlement", cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer group.Close()

	agCfg, err := agentgo.LoadFromEnv()
	if err != nil {
		log.Printf("agent disabled: %v", err)
	}
	ag := agentgo.New(agCfg)
	_ = ag.Start(context.Background())
	defer ag.Stop()

	handler := &paymentHandler{}
	wrapped := kafka.WrapConsumerHandler(handler, func(ev event.CapturedEvent) {
		log.Printf("captured topic=%s offset=%d", ev.Topic, ev.Offset)
	}, agCfg.CaptureTimeout)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		for {
			if err := group.Consume(ctx, []string{topic}, wrapped); err != nil {
				log.Printf("consume: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	log.Printf("payment-worker listening on %s topic=%s", brokers, topic)
	<-ctx.Done()
}

type paymentHandler struct{}

func (paymentHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (paymentHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *paymentHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("settle partition=%d offset=%d", msg.Partition, msg.Offset)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
