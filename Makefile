.PHONY: build
build:
	go build cmd/loy-point-api/main.go

.PHONY: test
test:
	go test -v ./...