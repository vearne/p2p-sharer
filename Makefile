VERSION :=0.0.1

RELEASE_DIR = dist
COMMAND_NAME = p2p-sharer
LDFLAGS = -ldflags "-s -w "

build: $(RELEASE_DIR)/$(COMMAND_NAME) ## Build release binaries

clean: ## Remove release binaries
	rm -rf $(RELEASE_DIR)

build-dirs:
	mkdir -p $(RELEASE_DIR)


$(RELEASE_DIR)/$(COMMAND_NAME): build-dirs $(wildcard *.go)
	env GOOS=linux GOARCH=amd64 go build -o $@ $(LDFLAGS)

