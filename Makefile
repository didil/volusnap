test:
	go test ./pkg/...
test-cover:
	go test -race -coverprofile cover.out -covermode=atomic  ./pkg/...
	go tool cover -html=cover.out -o cover.html
	open cover.html
deps:
	go get -u ./...
	go get -u github.com/stretchr/testify/assert
	go get -u -t github.com/volatiletech/sqlboiler
	go get -u github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql
deps-ci: deps
	go get golang.org/x/tools/cmd/cover
	go get -v github.com/rubenv/sql-migrate/...
test-ci:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./pkg/...
build:
	go build -o volusnapd cmd/volusnapd/* 
sqlboiler:
	sqlboiler psql --wipe

migrate-test-up:
	sql-migrate up -config=sql-migrate.yml -env=test
migrate-test-down:
	sql-migrate down -config=sql-migrate.yml -env=test
migrate-test-status:
	sql-migrate status -config=sql-migrate.yml -env=test

migrate-dev-up:
	sql-migrate up -config=sql-migrate.yml -env=development
migrate-dev-down:
	sql-migrate down -config=sql-migrate.yml -env=development
migrate-dev-status:
	sql-migrate status -config=sql-migrate.yml -env=development

migrate-prod-up:
	sql-migrate up -config=sql-migrate.yml -env=production
migrate-prod-down:
	sql-migrate down -config=sql-migrate.yml -env=production
migrate-prod-status:
	sql-migrate status -config=sql-migrate.yml -env=production