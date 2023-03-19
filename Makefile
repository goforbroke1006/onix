SERVICE_NAME=onix
FRONTEND_DIR=./web/app

.PHONY: all
all: prepare build test lint

prepare:
	go mod download
	go generate ./...
	go mod tidy
	cd ${FRONTEND_DIR} && npm install && cd -
.PHONY: prepare

build: build/backend build/frontend
.PHONY: build

build/backend:
	CGO_ENABLED=0 go build -o ./application ./
.PHONY: build/backend

build/frontend:
	npm --prefix ${FRONTEND_DIR} run build
.PHONY: build/frontend

test: test/backend test/frontend
.PHONY: test

test/backend:
	go test -short -coverprofile=./coverage.out ./...
.PHONY: test/backend

test/frontend:
	npm --prefix ${FRONTEND_DIR} test -- --watchAll=false
.PHONY: test/frontend

lint:
	golangci-lint run
	cd ${FRONTEND_DIR} && npm run lint && cd -
.PHONY: lint

coverage:
	go test -short -coverprofile=./coverage.out ./...
	go tool cover -html ./coverage.out
.PHONY: coverage

benchmark:
	go test -gcflags="-N" ./... -bench=.
.PHONY: benchmark

image:
	docker build --pull -f .docker/backend/Dockerfile  -t docker.io/goforbroke1006/onix-backend:latest  ./
	docker build --pull -f .docker/frontend/Dockerfile -t docker.io/goforbroke1006/onix-frontend:latest ./
.PHONY: image

clean:
	find ./ -name '*.generated.go' -type f -delete
	rm -rf ./application
	rm -rf ${FRONTEND_DIR}/build/ || true
	rm -rf ${FRONTEND_DIR}/node_modules/ || true
	rm -f ./coverage.out
.PHONY: clean

gen/frontend/snapshot:
	npm --prefix ${FRONTEND_DIR} test -- -u --watchAll=false
.PHONY: gen/frontend/snapshot

dev:
	docker build -f ./.docker-compose/backend/Dockerfile  -t "local-env/${SERVICE_NAME}-backend:dev"  .
	docker build -f ./.docker-compose/frontend/Dockerfile -t "local-env/${SERVICE_NAME}-frontend:dev" .
