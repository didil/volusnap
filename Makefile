test:
	go test ./pkg/...
test-cover:
	go test -race -coverprofile cover.out -covermode=atomic  ./pkg/...
	go tool cover -html=cover.out -o cover.html
	open cover.html
deps:
	go get -u ./...
	go get -u github.com/stretchr/testify/assert
deps-ci: deps
	go get golang.org/x/tools/cmd/cover
test-ci:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./pkg/...
build:
	go build cmd/volusnapd/volusnapd.go
	go build cmd/volusnapctl/volusnapctl.go