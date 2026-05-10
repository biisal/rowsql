GREETING := Hello from RowSQL!

.PHONY: default backend-build frontend-build backend-dev frontend-dev dev build run doc test release lint clean install format

default:
	@echo "$(GREETING)"

generate-schema:
	go run ./cmd/schema/main.go

backend-build:
	go run ./cmd/schema/main.go
	go build -o bin/rowsql ./cmd/server

frontend-build:

	cd ./frontend && pnpm run build

backend-dev:
	air -c air.toml

frontend-dev: generate-schema
	cd ./frontend/ && pnpm run dev

dev: 
	$(MAKE) -j2 frontend-dev backend-dev

build: doc frontend-build backend-build
	echo "build was successful"

run: build
	./bin/rowsql

doc:
	go run ./cmd/schema/...

test:
	go test ./... -failfast

release:
	goreleaser release --clean --snapshot

lint:
	golangci-lint run
	cd ./frontend && pnpm lint

lint-staged:
	cd frontend && pnpx lint-staged

clean:
	rm -rf dist/
	rm -rf bin/
	rm -rf tmp/
	rm -rf frontend/dist/
	rm -rf frontend/node_modules/

install:
	go mod download
	cd ./frontend && pnpm install

format:
	gofmt -w .
	cd ./frontend && pnpm run format

format-check:
	gofmt -l .
	cd ./frontend && pnpm run format:check