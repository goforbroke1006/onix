all: dep gen build test/unit

.PHONY: dep gen build test/unit test/functional test/integration

clean:
	find ./ -name '*.generated.go' | xargs rm -f '{}'
	rm -rf ./onix
	rm -rf ./frontend/build || true
	rm -rf ./frontend/node_modules || true

dep:
	go mod download
	npm --prefix ./frontend install

gen:
	go generate ./...

build:
	go build ./
	npm --prefix ./frontend run build

test/unit:
	go test --tags=unit ./...

test/functional:
	go test --tags=functional ./...

test/integration:
	go test --tags=integration ./...
