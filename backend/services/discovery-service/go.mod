module github.com/scouttalent/discovery-service

go 1.23

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/jackc/pgx/v5 v5.5.5
	github.com/scouttalent/pkg v0.0.0
	go.uber.org/zap v1.27.0
)

replace github.com/scouttalent/pkg => ../../pkg