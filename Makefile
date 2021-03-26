APP?=boolang
REGISTRY?='private-registry'
COMMIT_SHA=$(shell git rev-parse --short HEAD)

.PHONY: build
## build: build the application
build: clean
	@echo "Building..."
	@go build -o ${APP} main.go

.PHONY: run
## run: runs the go run build-binary
run: build
	./boolang

.PHONY: watch
## watch: watch the project for go file changes
watch:
	ulimit -n 1000 
	reflex -s -r '\.go$$' make run

.PHONY: clean
## clean: cleans the binary
clean:
	@echo "Cleaning..."
	@go clean

.PHONY: test
## test: runs go test with default values
test:
	go test -v -count=1 -race ./...


.PHONY: setup
## setup: setup go modules
setup:
	@go mod init \
		&& go mod tidy \
		&& go mod vendor

## helper rule for deployment
check-environment:
ifndef ENV
	${error ENV not set, allowed values - `staging` or `production`}
endif


.PHONY: docker-build
## docker-build: builds the boolang docker image to registry
docker-build: build
	docker build -t ${APP} .
	docker tag ${APP} ${APP}:${COMMIT_SHA}

.PHONY: docker-push
## docker-push: pushes the boolang docker image to registry
docker-push: docker-build
	docker push ${REGISTRY}/${ENV}/${APP}:${COMMIT_SHA}



.PHONY: help
## help: prints this help message
help:

	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

