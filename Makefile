.PHONY: all
all: prepare build test lint

prepare:
	go mod download
	go generate ./...
	go mod tidy
	npm --prefix ./web/dashboard-main-app/ install
.PHONY: prepare

.PHONY: build
build: build/backend build/frontend

build/backend:
	CGO_ENABLED=0 go build -o ./application ./
.PHONY: build/backend

build/frontend:
	npm --prefix ./web/dashboard-main-app/ run build
.PHONY: build/frontend

.PHONY: test
test: test/backend test/frontend

test/backend:
	go test -short -race -coverprofile=./coverage.out ./...
.PHONY: test/backend

test/frontend:
	npm --prefix ./web/dashboard-main-app/ test -- --watchAll=false
.PHONY: test/frontend

lint:
	golangci-lint run
	ineffassign ./...
	cd ./web/dashboard-main-app/ && eslint src/**/*.js && cd ./../../
.PHONY: lint

coverage:
	go test -short -race -coverprofile=./coverage.out ./...
	go tool cover -html ./coverage.out
.PHONY: coverage

benchmark:
	go test -gcflags="-N" ./... -bench=.
.PHONY: benchmark

image:
	docker build --pull -f .docker/backend/Dockerfile -t docker.io/goforbroke1006/onix-backend:latest ./
	docker build --pull -f .docker/frontend/Dockerfile -t docker.io/goforbroke1006/onix-frontend:latest ./
.PHONY: image

clean:
	find ./ -name '*.generated.go' -type f -delete
	rm -rf ./application
	rm -rf ./web/dashboard-main-app/build/ || true
	rm -rf ./web/dashboard-main-app/node_modules/ || true
	rm -f ./coverage.out
.PHONY: clean

gen/frontend/snapshot:
	@echo "Generate jest test snapshots"
	npm --prefix ./frontend/dashboard-admin/ test -- -u --watchAll=false
	npm --prefix ./frontend/dashboard-main/ test -- -u --watchAll=false
.PHONY: gen/frontend/snapshot
