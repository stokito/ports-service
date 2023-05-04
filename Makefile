.PHONY: build_docker_image build test run-lint clean

build: test
	mkdir -p ./build/package/
	GOOS=linux go build -o ./build/package/ports-service ./cmd/service
	GOOS=linux go build -o ./build/package/ports-import ./cmd/import

build_docker_image: clean build
	docker build -t ports-service:latest --no-cache ./

clean:
	rm -rf build/*

test:
	go test ./...

run-lint:
	golangci-lint run
