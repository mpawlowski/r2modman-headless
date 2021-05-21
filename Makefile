default: build

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

clean:
	@bazelisk clean || { echo "Clean failed, check above for errors!"; exit 1; }

build:
	@bazelisk run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies || { echo "Unable to run dependency update, check above for errors!"; exit 1; }
	@bazelisk run //:gazelle || { echo "Unable to run gazelle, check above for errors!"; exit 1; }
	@bazelisk build //... || { echo "Build failed, check above for errors!"; exit 1; }
	@bazelisk test //... || { echo "Tests failed, check above for errors!"; exit 1; }
