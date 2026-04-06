BIN_DIR=bin
CLI_BINARY=$(BIN_DIR)/shortener-cli

.PHONY: all
all: build

build: clean
	@mkdir -p $(BIN_DIR)
	go build -o $(CLI_BINARY) ./cmd/cli/main.go
	@echo "Binaries generated"

clean:
	@rm -rf $(BIN_DIR)

DOCKER_COMPOSE=deploy/docker-compose.yml
docker-up:
	podman compose -f $(DOCKER_COMPOSE) up -d

docker-down:
	podman compose -f $(DOCKER_COMPOSE) down
