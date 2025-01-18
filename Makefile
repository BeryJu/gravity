SHELL := /bin/bash
.ONESHELL:
.SHELLFLAGS += -x -e -o pipefail
.PHONY:  web
PWD = $(shell pwd)
UID = $(shell id -u)
GID = $(shell id -g)
VERSION = "0.25.0"
LD_FLAGS = -X beryju.io/gravity/pkg/extconfig.Version=${VERSION}
GO_FLAGS = -ldflags "${LD_FLAGS}" -v
SCHEMA_FILE = schema.yml
TEST_COUNT = 1
TEST_FLAGS =

ci--env:
	echo "sha=${GITHUB_SHA}" >> ${GITHUB_OUTPUT}
	echo "build=${GITHUB_RUN_ID}" >> ${GITHUB_OUTPUT}
	echo "timestamp=$(shell date +%s)" >> ${GITHUB_OUTPUT}
	echo "version=${VERSION}" >> ${GITHUB_OUTPUT}

docker-build: internal/resources/macoui internal/resources/blocky internal/resources/tftp
	go build \
		-ldflags "${LD_FLAGS} -X beryju.io/gravity/pkg/extconfig.BuildHash=${GIT_BUILD_HASH}" \
		${GO_BUILD_FLAGS} \
		-v -a -o gravity ${PWD}

clean:
	rm -rf ${PWD}/data/
	rm -rf ${PWD}/bin/

run: internal/resources/macoui internal/resources/blocky internal/resources/tftp
	export INSTANCE_LISTEN=0.0.0.0
	export DEBUG=true
	export LISTEN_ONLY=true
	$(eval LD_FLAGS := -X beryju.io/gravity/pkg/extconfig.Version=${VERSION} -X beryju.io/gravity/pkg/extconfig.BuildHash=dev-$(shell git rev-parse HEAD))
	go run ${GO_FLAGS} ${PWD} server

# Web
web: web-lint web-build

web-install:
	cd ${PWD}/web
	npm ci
	npm version ${VERSION} || true

web-build:
	cd ${PWD}/web
	npm run build

web-watch:
	cd ${PWD}/web
	npm run watch

web-lint:
	cd ${PWD}/web
	npm run prettier
	npm run lint
	npm run lit-analyse

# CLI
bin/gravity-cli:
	$(eval LD_FLAGS := -X beryju.io/gravity/pkg/extconfig.Version=${VERSION} -X beryju.io/gravity/pkg/extconfig.BuildHash=dev-$(shell git rev-parse HEAD))
	mkdir -p ${PWD}/bin/
	go build ${GO_FLAGS} -o ${PWD}/bin/gravity-cli ${PWD}/cmd/cli/main/

# Website
website-watch:
	cd ${PWD}/docs
	open http://localhost:1313/ && hugo server --noBuildLock

internal/resources/macoui:
	mkdir -p internal/resources/macoui
	curl -L https://raw.githubusercontent.com/wireshark/wireshark/6885d787fda5f74a2d1f9eeea443fecf8dd58528/manuf -o ${PWD}/internal/resources/macoui/db.txt

internal/resources/blocky:
	mkdir -p internal/resources/blocky
	curl -L https://adaway.org/hosts.txt -o ${PWD}/internal/resources/blocky/adaway.org.txt
	curl -L https://big.oisd.nl/domainswild -o ${PWD}/internal/resources/blocky/big.oisd.nl.txt
	curl -L https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts -o ${PWD}/internal/resources/blocky/StevenBlack.hosts.txt
	curl -L https://v.firebog.net/hosts/AdguardDNS.txt -o ${PWD}/internal/resources/blocky/AdguardDNS.txt
	curl -L https://v.firebog.net/hosts/Easylist.txt -o ${PWD}/internal/resources/blocky/Easylist.txt
	curl -L https://adguardteam.github.io/AdGuardSDNSFilter/Filters/filter.txt -o ${PWD}/internal/resources/blocky/AdGuardSDNSFilter.txt

internal/resources/tftp:
	mkdir -p internal/resources/tftp
	curl -L http://boot.ipxe.org/undionly.kpxe -o ${PWD}/internal/resources/tftp/ipxe.undionly.kpxe
	curl -L https://boot.netboot.xyz/ipxe/netboot.xyz.kpxe -o ${PWD}/internal/resources/tftp/netboot.xyz.kpxe
	curl -L https://boot.netboot.xyz/ipxe/netboot.xyz-undionly.kpxe -o ${PWD}/internal/resources/tftp/netboot.xyz-undionly.kpxe
	curl -L https://boot.netboot.xyz/ipxe/netboot.xyz.efi -o ${PWD}/internal/resources/tftp/netboot.xyz.efi

gen-build:
	DEBUG=true go run ${GO_FLAGS} ${PWD} generateSchema ${SCHEMA_FILE}
	git add ${SCHEMA_FILE}

gen-proto:
	protoc \
		--proto_path . \
		--go_out . \
		protobuf/**

gen-clean:
	rm -rf ${PWD}/gen-ts-api/
	rm -rf ${PWD}/api/api/
	rm -rf ${PWD}/api/docs/
	rm -rf ${PWD}/api/test/
	rm -rf ${PWD}/api/*.go

gen-tag:
	git add Makefile
	cd ${PWD}
	git commit -m "release version v${VERSION}"
	git tag v${VERSION}

gen-client-go:
	docker run \
		--rm -v ${PWD}:/local \
		--user ${UID}:${GID} \
		openapitools/openapi-generator-cli:v6.6.0 generate \
		--git-host beryju.io \
		--git-user-id gravity \
		--git-repo-id api \
		--additional-properties=packageName=api \
		-i /local/schema.yml \
		-g go \
		-o /local/api \
		-c /local/api/config.yaml
	cd ${PWD}/api/
	rm -f .travis.yml go.mod go.sum
	go get
	go fmt .
	go mod tidy
	gofumpt -l -w . || true
	git add .

gen-client-ts:
	docker run \
		--rm -v ${PWD}:/local \
		--user ${UID}:${GID} \
		openapitools/openapi-generator-cli:v6.6.0 generate \
		-i /local/${SCHEMA_FILE} \
		-g typescript-fetch \
		-o /local/gen-ts-api \
		--additional-properties=typescriptThreePlus=true,supportsES6=true,npmName=gravity-api,npmVersion=${VERSION} \
		--git-repo-id BeryJu \
		--git-user-id gravity
	cd ${PWD}/gen-ts-api && npm i
	\cp -rf ${PWD}/gen-ts-api/* ${PWD}/web/node_modules/gravity-api

gen-client-ts-publish: gen-client-ts
	cd ${PWD}/gen-ts-api
	npm publish
	cd ${PWD}/web
	npm i gravity-api@${VERSION}
	npm version ${VERSION} || true
	git add package*.json

release: gen-build gen-clean gen-client-go gen-client-ts-publish gen-tag

lint: web-lint
	golangci-lint run -v --timeout 5000s

test-env-start:
	cd ${PWD}/hack/tests/
	docker compose --project-name gravity-test-env up -d

test-env-stop:
	cd ${PWD}/hack/tests/
	docker compose --project-name gravity-test-env down -v

install-deps: internal/resources/macoui internal/resources/blocky internal/resources/tftp
	sudo apt-get install -y nmap libpcap-dev
	sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"

test-local:
	$(eval TEST_COUNT := 100)
	$(eval TEST_FLAGS := -shuffle=on -failfast)

test: internal/resources/macoui internal/resources/blocky internal/resources/tftp
	export BOOTSTRAP_ROLES="dns;dhcp;api;discovery;backup;debug;tsdb;tftp"
	export ETCD_ENDPOINT="localhost:2385"
	export DEBUG="true"
	export LISTEN_ONLY="true"
	export CI="true"
	go test \
		-p 1 \
		-v \
		-coverprofile=${PWD}/coverage.txt \
		-covermode=atomic \
		-count=${TEST_COUNT} \
		${TEST_FLAGS} \
		$(shell go list ./... | grep -v beryju.io/gravity/api) \
			2>&1 | tee test-output
	go tool cover \
		-html ${PWD}/coverage.txt \
		-o ${PWD}/coverage.html

test-e2e-container-build:
	docker build \
		--build-arg=GRAVITY_BUILD_ARGS=GO_BUILD_FLAGS=-cover \
		-t gravity:e2e-test \
		.

test-e2e: test-e2e-container-build
	cd ${PWD}/tests
	go get .
	go test \
		-p 1 \
		-v \
		-coverprofile=${PWD}/coverage.txt \
		-covermode=atomic \
		-count=${TEST_COUNT} \
		${TEST_FLAGS} \
		beryju.io/gravity/tests \
			2>&1 | tee test-output
	cd ${PWD}
	mkdir -p ${PWD}/tests/coverage-node-1/ \
		${PWD}/tests/coverage-node-2/ \
		${PWD}/tests/coverage-node-3/
	go tool covdata textfmt \
		-i ${PWD}/tests/coverage-node-1/ \
		-i ${PWD}/tests/coverage-node-2/ \
		-i ${PWD}/tests/coverage-node-3/ \
		--pkg $(shell go list ./... | grep -v beryju.io/gravity/api | xargs | sed 's/ /,/g') \
		-o ${PWD}/coverage_in_container.txt
	go tool cover \
		-html ${PWD}/coverage_in_container.txt \
		-o ${PWD}/coverage.html

bench: internal/resources/macoui internal/resources/blocky internal/resources/tftp
	export BOOTSTRAP_ROLES="dns;dhcp;api;discovery;backup;debug;tsdb;tftp"
	export ETCD_ENDPOINT="localhost:2385"
	export LISTEN_ONLY="true"
	export LOG_LEVEL="FATAL"
	export CI="true"
	go test \
		-run=^$$ \
		-bench=^Benchmark \
		-benchmem \
		$(shell go list ./... | grep -v beryju.io/gravity/api) \
			| tee test-output
