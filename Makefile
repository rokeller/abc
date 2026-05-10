VERSION_RAW = $(shell git describe --tags)

debug:
	@go build -ldflags \
		"-X github.com/rokeller/abc/cmd.version=local-${VERSION_RAW}"

release:
	@go build -ldflags "-s -w \
		-X github.com/rokeller/abc/cmd.version=local-${VERSION_RAW}"

test: debug
	@go test ./...

cover: debug
	@go test ./... -coverprofile=coverage.out
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@go-cover-treemap -coverprofile coverage.out > coverage.svg

blah:
	# @go tool cover -html=coverage.out
	# @go tool cover -html=coverage.out
