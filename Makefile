all: go-fmt lint

fmt: go-fmt

lint:
	golangci-lint run

go-fmt:
	gofumpt -l -w .

