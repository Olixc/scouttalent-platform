module github.com/scouttalent/media-service

go 1.23

require (
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.3.1
	github.com/gin-gonic/gin v1.10.0
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.5.5
	github.com/nats-io/nats.go v1.34.0
	github.com/scouttalent/pkg v0.0.0
	github.com/tus/tusd/v2 v2.4.0
	go.uber.org/zap v1.27.0
)

replace github.com/scouttalent/pkg => ../../pkg