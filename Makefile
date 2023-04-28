.PHONY: build
build:
	go build -o dist/terraform-j2md cmd/terraform-j2md/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: fmt
fmt:
	go fmt ./...