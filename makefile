tools-local: tool-golangci-lint

tool-golangci-lint:
	scripts/goget.sh github.com/golangci/golangci-lint/cmd/golangci-lint

test:
	go test -race ./... --count=1

test-integration:
	go test -tags integration ./internal/infra/http/controllers_test/... --count=1

lint:
	golangci-lint run
	go mod tidy -v && git --no-pager diff --quiet go.mod go.sum
