all: dep gen build test/unit

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
	go generate ./...

build:
	go build ./
	npm --prefix ./frontend/dashboard-admin/ run build
	npm --prefix ./frontend/dashboard-main/ run build

test/unit:
	go test --tags=unit ./...

test/functional:
	go test --tags=functional ./...

test/integration:
	go test --tags=integration ./...
