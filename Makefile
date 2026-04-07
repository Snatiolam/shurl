BIN_DIR=bin
CLI_BINARY=$(BIN_DIR)/shortener-cli
SERVER_BINARY=$(BIN_DIR)/shortener-server

.PHONY: all
all: build

build: clean
	@mkdir -p $(BIN_DIR)
	go build -v -o $(CLI_BINARY) ./cmd/cli/main.go
	go build -v -o $(SERVER_BINARY) ./cmd/server/main.go
	@echo "Binaries generated"

clean:
	@rm -rf $(BIN_DIR)

DOCKER_COMPOSE=deploy/docker-compose.yml
docker-up:
	podman compose -f $(DOCKER_COMPOSE) up -d

docker-down:
	podman compose -f $(DOCKER_COMPOSE) down
