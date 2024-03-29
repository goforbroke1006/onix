.PHONY: all
all: dep gen build test lint

.PHONY: clean
clean:
	find ./ -name '*.generated.go' -type f -delete
	rm -rf ./onix
	rm -rf ./frontend/dashboard-admin/build/ || true
	rm -rf ./frontend/dashboard-admin/node_modules/ || true
	rm -rf ./frontend/dashboard-main/build/ || true
	rm -rf ./frontend/dashboard-main/node_modules/ || true
	rm -f ./coverage.out
	rm -rf internal/repository/mocks || true
	rm -rf internal/service/mocks || true

.PHONY: dep
dep:
	@echo "Install backend dependencies"
	go mod download
	go mod tidy
	@echo "Install frontend dependencies"
	npm --prefix ./frontend/dashboard-admin/ install
	npm --prefix ./frontend/dashboard-main/  install

.PHONY: gen
gen:
	@echo "Generate backend boilerplate code"
	go generate ./...

.PHONY: gen/frontend/snapshot
gen/frontend/snapshot:
	@echo "Generate jest test snapshots"
	npm --prefix ./frontend/dashboard-admin/ test -- -u --watchAll=false
	npm --prefix ./frontend/dashboard-main/ test -- -u --watchAll=false

.PHONY: build
build: build/backend build/frontend

.PHONY: build/backend
build/backend:
	@echo "Build backend"
	CGO_ENABLED=0 go build ./

.PHONY: build/frontend
build/frontend:
	@echo "Build frontend"
	npm --prefix ./frontend/dashboard-admin/ run build
	npm --prefix ./frontend/dashboard-main/ run build

.PHONY: test
test: test/backend test/frontend

.PHONY: test/backend
test/backend:
	go test -race ./...

.PHONY: test/frontend
test/frontend:
	npm --prefix ./frontend/dashboard-admin/ test -- --watchAll=false
	npm --prefix ./frontend/dashboard-main/ test -- --watchAll=false

.PHONY: lint
lint:
	golangci-lint run
	ineffassign ./...
	find . -type f -name '*.go' | xargs misspell -error
	cd ./frontend/dashboard-main/ && eslint src/**/*.js && cd ./../../

.PHONY: benchmark
benchmark:
	go test -gcflags="-N" ./... -bench=.

.PHONY: coverage
coverage:
	go test -race -coverprofile=./coverage.out ./...
	go tool cover -html ./coverage.out

.PHONY: image
image:
	echo "Build"
	docker build --pull --network=host -f ./.build/register/Dockerfile -t docker.io/goforbroke1006/onix-register:latest ./.build/register
	DOCKER_BUILDKIT=1 docker build --pull --network=host -f .build/backend/Dockerfile -t docker.io/goforbroke1006/onix-backend:latest ./
	docker build --pull --network=host -f .build/frontend/Dockerfile -t docker.io/goforbroke1006/onix-dashboard-admin:latest ./frontend/dashboard-admin

image/frontend/dashboard-main:
	docker build --pull --network=host -f .build/frontend/main/Dockerfile -t docker.io/goforbroke1006/onix-dashboard-main:latest ./
	docker push docker.io/goforbroke1006/onix-dashboard-main:latest

image/frontend/dashboard-admin:
	docker build --pull --network=host -f .build/frontend/admin/Dockerfile -t docker.io/goforbroke1006/onix-dashboard-admin:latest ./
	docker push docker.io/goforbroke1006/onix-dashboard-admin:latest


.PHONY: setup
setup:
	bash ./setup.sh
