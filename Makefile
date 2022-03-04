all: dep gen build test lint

.PHONY: dep gen build test/unit test/functional test/integration lint coverage

clean:
	find ./ -name '*.generated.go' -type f -delete
	rm -rf ./onix
	rm -rf ./frontend/dashboard-admin/build/ || true
	rm -rf ./frontend/dashboard-admin/node_modules/ || true
	rm -rf ./frontend/dashboard-main/build/ || true
	rm -rf ./frontend/dashboard-main/node_modules/ || true
	rm -f ./.coverage

dep:
	go mod download
	npm --prefix ./frontend/dashboard-admin/ install
	npm --prefix ./frontend/dashboard-main/ install

gen:
	@echo "Generate backend boilerplate code"
	go generate ./...

gen/frontend/snapshot:
	@echo "Generate jest test snapshots"
	npm --prefix ./frontend/dashboard-admin/ test -- -u --watchAll=false
	npm --prefix ./frontend/dashboard-main/ test -- -u --watchAll=false

build:
	go build ./
	npm --prefix ./frontend/dashboard-admin/ run build
	npm --prefix ./frontend/dashboard-main/ run build

test: test/unit
test-all: test/unit test/functional test/integration

test/unit:
	go test --tags=unit ./... -cover
	npm --prefix ./frontend/dashboard-admin/ test -- --watchAll=false
	npm --prefix ./frontend/dashboard-main/ test -- --watchAll=false

test/functional:
	go test --tags=functional ./...

test/integration:
	go test --tags=integration ./...

lint:
	golangci-lint run
	cd ./frontend/dashboard-main/ && eslint src/**/*.js && cd ./../../

benchmark:
	go test -gcflags="-N" ./... -bench=.

coverage:
	go test --coverprofile ./.coverage ./...
	go tool cover -html ./.coverage

setup:
	bash ./setup.sh