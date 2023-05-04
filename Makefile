.PHONY: build test run-lint clean
build: test
	mkdir -p ./build/package/
	GOOS=linux go build -o ./build/package/ports-service ./cmd/service
	GOOS=linux go build -o ./build/package/ports-import ./cmd/import

clean:
	rm -rf build/*

test:
	go test ./...

run-lint:
	golangci-lint run
