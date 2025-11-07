SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command

lint:
	golangci-lint run -v

local-env:
	docker-compose --env-file ./config/.env up -d
	$$env:ENV_PATH="./config/.env"; go run cmd/migrate/main.go --command=up
	go mod tidy
	$$env:ENV_PATH="./config/.env"; go run cmd/debezium/main.go
stop:
	docker stop postgres_container
	docker stop debezium_app
docker-build:
	docker-compose --env-file ./config/.env up --build -d