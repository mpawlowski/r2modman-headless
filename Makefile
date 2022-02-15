default: all

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

ERROR_PREFIX:=[ \033[0;31mBUILD ERROR\033[0m ]

all: clean build bin

clean:
	@bazelisk clean || { echo "$(ERROR_PREFIX) Clean failed, check above for errors!"; exit 1; }
	@rm -rf $(ROOT_DIR)/work || { echo "$(ERROR_PREFIX) Clean failed, cannot remove work directory!"; exit 1; }
	@rm -rf $(ROOT_DIR)/r2modman-headless || { echo "$(ERROR_PREFIX) Clean failed, cannot remove binary!"; exit 1; }

build:
	@bazelisk run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies || { echo "$(ERROR_PREFIX) Unable to run dependency update, check above for errors!"; exit 1; }
	@bazelisk run //:gazelle -- update || { echo "$(ERROR_PREFIX) Unable to run gazelle, check above for errors!"; exit 1; }
	@bazelisk test //... || { echo "$(ERROR_PREFIX) Tests failed, check above for errors!"; exit 1; }

bin:
	@rm -rf $(ROOT_DIR)/r2modman-headless
	@cp $(shell bazelisk run --run_under "echo " //:r2modman-headless 2> /dev/null) r2modman-headless || ( echo "$(ERROR_PREFIX) Unable to copy bin."; exit 1; )

dev: build
	mkdir -p $(ROOT_DIR)/work/install
	mkdir -p $(ROOT_DIR)/work/cache
	bazelisk run //:r2modman-headless -- --install-dir=$(ROOT_DIR)/work/install --profile-zip=$(ROOT_DIR)/r2modman/testdata/Valheim_Creative_Mode.r2z --work-dir=$(ROOT_DIR)/work/cache