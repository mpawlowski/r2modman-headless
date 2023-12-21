default: all

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

ERROR_PREFIX:=[ \033[0;31mBUILD ERROR\033[0m ]

all: clean build

clean:
	@rm -rf $(ROOT_DIR)/work || { echo "$(ERROR_PREFIX) Clean failed, cannot remove work directory!"; exit 1; }
	@rm -rf $(ROOT_DIR)/r2modman-headless || { echo "$(ERROR_PREFIX) Clean failed, cannot remove binary!"; exit 1; }

build:
	@go build || { echo "$(ERROR_PREFIX) Tests failed, check above for errors!"; exit 1; }

dev: build
	mkdir -p $(ROOT_DIR)/work/install
	mkdir -p $(ROOT_DIR)/work/cache
	$(ROOT_DIR)/r2modman-headless --install-dir=$(ROOT_DIR)/work/install --profile-zip=$(ROOT_DIR)/r2modman/testdata/Megamodheim.r2z --work-dir=$(ROOT_DIR)/work/cache