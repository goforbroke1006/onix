all: dep gen build test

.PHONY: dep gen build test/unit test/functional test/integration

clean:
	find ./ -name '*.generated.go' -type f -delete
	rm -rf ./onix
	rm -rf ./frontend/dashboard-admin/build/ || true
	rm -rf ./frontend/dashboard-admin/node_modules/ || true
	rm -rf ./frontend/dashboard-main/build/ || true
	rm -rf ./frontend/dashboard-main/node_modules/ || true

dep:
	go mod download
	npm --prefix ./frontend/dashboard-admin/ install
	npm --prefix ./frontend/dashboard-main/ install

gen:
	@echo "Generate backend boilerplate code"
	go generate ./...
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
	go test --tags=unit ./...
	npm --prefix ./frontend/dashboard-admin/ test -- --watchAll=false
	npm --prefix ./frontend/dashboard-main/ test -- --watchAll=false

test/functional:
	go test --tags=functional ./...

test/integration:
	go test --tags=integration ./...

setup:
	bash ./setup.sh