ensure-psql:
	@systemctl is-active --quiet postgresql || sudo systemctl start postgresql

run:
	./bin/rowsql

frontend-dev:
	cd ./frontend/ && pnpm run dev

backend-dev:
	air -c air.toml

dev:
	make -j2 frontend-dev backend-dev

build:
	cd ./frontend && pnpm run build
	go build -o bin/rowsql ./cmd/server
	
build-linux:
	cd ./frontend && pnpm run build
	GOOS=linux GOARCH=amd64 go build -o bin/rowsql-linux ./cmd/server

test:
	go test ./...

.PHONY: frontend-dev backend-dev dev run ensure-psql
