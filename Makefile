.PHONY: build clean test release
VERSION=`egrep -o '[0-9]+\.[0-9a-z.\-]+' version.go`
GIT_SHA=`git rev-parse --short HEAD || echo`
GIT_TAG=`git tag --points-at=HEAD`

build:
	@echo "Building watchman..."
	@mkdir -p bin
	@go build -ldflags "-X main.GitSHA=${GIT_SHA}" -o bin/watchman .

clean:
	@rm -f bin/*

test:
	@echo "Running tests..."
	@go test ./...

release:
	@docker build -q -t ttousai/watchman:${GIT_TAG} -f Dockerfile.alpine .
	@docker push ttousai/watchman:${GIT_TAG}
