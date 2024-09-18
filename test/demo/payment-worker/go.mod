module github.com/replay/platform/test/demo/payment-worker

go 1.25.0

require (
	github.com/IBM/sarama v1.46.3
	github.com/replay/platform/packages/agent-go v0.0.0
)

replace github.com/replay/platform/packages/agent-go => ../../../packages/agent-go

replace github.com/replay/platform/packages/shared-go => ../../../packages/shared-go
