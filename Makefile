all: build

.PHONY: build
build:
	go build -v

.PHONY: fmt
fmt:
	gofmt -s -w $$(git ls-files | grep -e '\.go$$' | grep -v -e vendor)

.PHONY: test
test:
	go test -v ./...

.PHONY: vet
vet:
	go tool vet --all *.go

.PHONY: test-all
test-all: vet test