include .env
export

install:
	go install github.com/volatiletech/sqlboiler/v4@latest
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
	go install github.com/golang/mock/mockgen@latest
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go get

run:
	go run main.go

run/feed-database:
	go run main.go feed_database

.PHONY: mock
mock:
	go generate ./...

test:
	go test ./...

test/bench:
	go test ./... -bench=.

test/cov:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

migrate:
	migrate -database ${MIGRATE_URL} -path db/migrations up

swag:
	swag init

models:
	sqlboiler mysql --wipe