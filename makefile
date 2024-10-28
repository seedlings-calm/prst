GITHUB_TOKEN =  $(shell cat ghp_sercet)
# Determine operating system
ifeq ($(OS),Windows_NT)
    OS_TYPE := Windows
else
    OS_TYPE := $(shell uname -s)
endif

# Define the default target
.PHONY: swag
swag: check_make check_swag generate_swag

# Check if make command exists, if not install it based on the OS
check_make:
    @echo "Checking for make command..."
    @command -v make >/dev/null 2>&1 || { \
        echo "make command not found. Installing make..."; \
        if [ "$(OS_TYPE)" = "Windows" ]; then \
            choco install make; \
        elif [ "$(OS_TYPE)" = "Darwin" ]; then \
            brew install make; \
        else \
            echo "Unsupported OS. Please install make manually."; exit 1; \
        fi; \
    }

# Check if swag command exists, if not install it based on the OS
check_swag:
    @echo "Checking for swag command..."
    @command -v swag >/dev/null 2>&1 || { \
        echo "swag command not found. Installing swag..."; \
        go install github.com/swaggo/swag/cmd/swag@latest; \
    }

# Generate Swag documentation
generate_swag:
    @echo "Generating Swag documentation..."
    swag init --parseDependency --parseDepth=6 --instanceName prst -o ./docs/   

# Define cleanup target (optional)
.PHONY: clean
clean:
    @echo "Cleaning up generated files..."
    rm -rf docs



window:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o prst.exe ./
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o prst ./

.PHONY: window linux


dev:
	go  run ./ -c config.dev.yml
release: 
	go run  ./ -c config.yml

.PHONY: dev release


build:
	@docker login ghcr.io -u seedlings-calm -p $(GITHUB_TOKEN)
	@docker build -t ghcr.io/seedlings-calm/prst:latest . 
	docker push ghcr.io/seedlings-calm/prst:latest

.PHONY: build



ssl:
	@chmod +x ./ssl/setup.sh
	./ssl/setup.sh
.PHONY: ssl


SHELL := /bin/bash
BINARY_NAME := YfApi
SERVICE_NAME := yf-api
SSH_USER := root
DEV_SERVER := 192.168.77.112
REMOTE_PATH := /data/$(SERVICE_NAME)
SYSTEMD_CONF_NAME := $(SERVICE_NAME).service
SYSTEMD_CONF_PATH := /etc/systemd/system
#max编译
mac-scp: build_mac_to_linux
	@echo "Deploying to server: $(DEV_SERVER)"
	@echo "Remove old binary file"
	ssh $(SSH_USER)@$(DEV_SERVER) "rm -f $(REMOTE_PATH)/$(BINARY_NAME)"
	@echo "Sending new binary file to dev server"
	scp $(BINARY_NAME) $(SSH_USER)@$(DEV_SERVER):$(REMOTE_PATH)
	@echo "Restart service"
	ssh $(SSH_USER)@$(DEV_SERVER) "systemctl restart $(SERVICE_NAME)"
	@echo "Done!"
build_mac_to_linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)
#windows编译
deploy_prod_windows_to_linux: build_windows_to_linux
	@echo "Deploying to server: $(DEV_SERVER)"
	@echo "Remove old binary file"
	ssh $(SSH_USER)@$(DEV_SERVER) "rm -f $(REMOTE_PATH)/$(BINARY_NAME)"
	@echo "Sending new binary file to dev server"
	scp $(BINARY_NAME) $(SSH_USER)@$(DEV_SERVER):$(REMOTE_PATH)
	ssh $(SSH_USER)@$(DEV_SERVER) "chmod +x $(REMOTE_PATH)/$(BINARY_NAME)"
	@echo "Restart service"
	ssh $(SSH_USER)@$(DEV_SERVER) "systemctl restart $(SERVICE_NAME)"
	@echo "Done!"
build_windows_to_linux:
	SET CGO_ENABLED=0
	SET GOOS=linux
	SET GOARCH=amd64
	go build -o $(BINARY_NAME)

.PHONY: mac-scp build_mac_to_linux deploy_prod_windows_to_linux


swag:
	swag init --parseDependency --parseDepth=6 --exclude ./core,./internal/,./util/  -o ./docs

dev:
	go run main.go  -f ./config/config.yaml

.PHONY: swag dev
