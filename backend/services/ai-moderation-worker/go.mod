module github.com/scouttalent/ai-moderation-worker

go 1.23

require (
	github.com/jackc/pgx/v5 v5.5.5
	github.com/nats-io/nats.go v1.34.0
	github.com/scouttalent/pkg v0.0.0
	go.uber.org/zap v1.27.0
)

replace github.com/scouttalent/pkg => ../../pkg