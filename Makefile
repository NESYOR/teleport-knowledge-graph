.PHONY: build test fmt run serve clean

BINARY := twin

build:
	go build -o $(BINARY) ./cmd/twin

test:
	go test ./...

fmt:
	gofmt -w $$(rg --files -g '*.go')

run:
	go run ./cmd/twin collect

serve:
	go run ./cmd/twin serve

clean:
	rm -f $(BINARY)
