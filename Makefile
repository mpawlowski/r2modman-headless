default: build

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

clean:
	@bazelisk clean || { echo "Clean failed, check above for errors!"; exit 1; }
	rm -rf $(ROOT_DIR)/work

build:
	@bazelisk run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies || { echo "Unable to run dependency update, check above for errors!"; exit 1; }
	@bazelisk run //:gazelle || { echo "Unable to run gazelle, check above for errors!"; exit 1; }
	@bazelisk build //... || { echo "Build failed, check above for errors!"; exit 1; }
	@bazelisk test //... || { echo "Tests failed, check above for errors!"; exit 1; }

dev: build
	mkdir -p $(ROOT_DIR)/work/install
	mkdir -p $(ROOT_DIR)/work/cache
	bazelisk run //:r2modman-headless -- --install-dir=$(ROOT_DIR)/work/install --profile-zip=$(ROOT_DIR)/r2modman/testdata/Valheim_Creative_Mode.r2z --work-dir=$(ROOT_DIR)/work/cache