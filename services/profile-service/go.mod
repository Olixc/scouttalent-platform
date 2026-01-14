module github.com/scouttalent/profile-service

go 1.22

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.5.5
	github.com/nats-io/nats.go v1.31.0
	github.com/scouttalent/pkg v0.0.0
	go.uber.org/zap v1.27.0
)

replace github.com/scouttalent/pkg => ../../pkg