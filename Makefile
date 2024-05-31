BINARY_NAME=minesweeper

install:
	@go mod tidy

dev:
	@go run cmd/minesweeper/main.go

build:
	@go build -o bin/${BINARY_NAME} cmd/minesweeper/main.go

run: build
	@./bin/${BINARY_NAME}

test:
	@go test -v ./... -count=1
