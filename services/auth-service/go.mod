module github.com/scouttalent/auth-service

go 1.22

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/jackc/pgx/v5 v5.5.5
	github.com/redis/go-redis/v9 v9.5.1
	github.com/scouttalent/pkg v0.0.0
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.20.0
)

replace github.com/scouttalent/pkg => ../../pkg