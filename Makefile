.PHONY: release build run run_with_redis

current_version_number := $(shell git tag --list "v*" | sort -V | tail -n 1 | cut -c 2-)
next_version_number := $(shell echo $$(($(current_version_number)+1)))
latest_logstream_version = $(shell hub api -X GET repos/codecrafters-io/logstream/releases/latest \
							| jq ".tag_name")

update_logstream_version:
	sed -E -i "/^ARG logstream_version/ s/v[0-9]+/$(latest_logstream_version)/g" Dockerfile
	git add Dockerfile
	git commit -m "Update logstream version"

release_docker:
	git push origin master
	git tag v$(next_version_number)
	git push origin v$(next_version_number)

release:
	rm -rf dist
	goreleaser

build:
	go build -o dist/main.out

test: build
	dist/main.out --binary-path=./run_success.sh --config-path=./test_helpers/valid_config.yml

bump_version:
	bumpversion --verbose --tag patch

bump_release_upload: bump_version release upload_to_travis
