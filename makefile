ensure-psql:
	@systemctl is-active --quiet postgresql || sudo systemctl start postgresql

run:
	cd ./frontend && pnpm run build
	go run ./cmd/server

frontend-dev:
	cd ./frontend/ && pnpm run dev

backend-dev:
	air -c air.toml

dev:
	make -j2 frontend-dev backend-dev


.PHONY: frontend-dev backend-dev dev run ensure-psql
