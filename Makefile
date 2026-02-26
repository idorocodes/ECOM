# Variables - change these to match your project name
BINARY_NAME=ecom-api
MAIN_PATH=./cmd/api/main.go

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## run: run the gateway application
run:
	go run ${MAIN_PATH}

## build: build the binary
build:
	go build -o ./bin/${BINARY_NAME} ${MAIN_PATH}

## test: run all tests
test:
	go test -v ./...

## clean: remove binary and temporary files
clean:
	go clean
	rm -f ./bin/${BINARY_NAME}