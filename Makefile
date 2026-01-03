SHELL := /bin/bash
.ONESHELL:
.SHELLFLAGS += -x -e -o pipefail
.PHONY:  web
PWD = $(shell pwd)
UID = $(shell id -u)
GID = $(shell id -g)
VERSION = "0.29.0"
LD_FLAGS = -X beryju.io/gravity/pkg/extconfig.Version=${VERSION}
GO_FLAGS = -ldflags "${LD_FLAGS}" -v
SCHEMA_FILE = schema.yml
TEST_COUNT = 1
TEST_FLAGS =
TEST_OUTPUT = ${PWD}/.test-output

GEN_API_TS = "gen-ts-api"
GEN_API_GO = "api"
GEN_ED_GO = "pkg/externaldns/generated/"

ci--env:
	echo "sha=${GITHUB_SHA}" >> ${GITHUB_OUTPUT}
	echo "build=${GITHUB_RUN_ID}" >> ${GITHUB_OUTPUT}
	echo "timestamp=$(shell date +%s)" >> ${GITHUB_OUTPUT}
	echo "version=${VERSION}" >> ${GITHUB_OUTPUT}

docker-build: internal/resources/macoui internal/resources/blocky internal/resources/tftp
	go build \
		-ldflags "${LD_FLAGS} -X beryju.io/gravity/pkg/extconfig.BuildHash=${GIT_BUILD_HASH}" \
		${GO_BUILD_FLAGS} \
		-v -a -o gravity ${PWD}/cmd/server/main

docker-build-external-dns:
	go build \
		-ldflags "${LD_FLAGS} -X beryju.io/gravity/pkg/extconfig.BuildHash=${GIT_BUILD_HASH}" \
		${GO_BUILD_FLAGS} \
		-v -a -o gravity-external-dns ${PWD}/cmd/external-dns/main

clean:
	rm -rf ${PWD}/data/
	rm -rf ${PWD}/bin/
	rm -rf ${PWD}/tests/coverage-*/

run: internal/resources/macoui internal/resources/blocky internal/resources/tftp
	export INSTANCE_LISTEN=0.0.0.0
	export DEBUG=true
	export LISTEN_ONLY=true
	$(eval LD_FLAGS := -X beryju.io/gravity/pkg/extconfig.Version=${VERSION} -X beryju.io/gravity/pkg/extconfig.BuildHash=dev-$(shell git rev-parse HEAD))
	go run \
		${GO_FLAGS} \
		${PWD}/cmd/server/main \
			server

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
	npm run tsc
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
	mkdir -p ${PWD}/internal/resources/macoui
	tshark -G manuf > ${PWD}/internal/resources/macoui/db.txt

internal/resources/blocky:
	mkdir -p ${PWD}/internal/resources/blocky
	curl -L https://adaway.org/hosts.txt -o ${PWD}/internal/resources/blocky/adaway.org.txt
	curl -L https://big.oisd.nl/domainswild -o ${PWD}/internal/resources/blocky/big.oisd.nl.txt
	curl -L https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts -o ${PWD}/internal/resources/blocky/StevenBlack.hosts.txt
	curl -L https://v.firebog.net/hosts/AdguardDNS.txt -o ${PWD}/internal/resources/blocky/AdguardDNS.txt
	curl -L https://v.firebog.net/hosts/Easylist.txt -o ${PWD}/internal/resources/blocky/Easylist.txt
	curl -L https://adguardteam.github.io/AdGuardSDNSFilter/Filters/filter.txt -o ${PWD}/internal/resources/blocky/AdGuardSDNSFilter.txt

internal/resources/tftp:
	mkdir -p ${PWD}/internal/resources/tftp
	curl -L https://boot.ipxe.org/ipxe.efi -o ${PWD}/internal/resources/tftp/ipxe.efi
	curl -L https://boot.ipxe.org/undionly.kpxe -o ${PWD}/internal/resources/tftp/ipxe.undionly.kpxe
	curl -L https://boot.netboot.xyz/ipxe/netboot.xyz.kpxe -o ${PWD}/internal/resources/tftp/netboot.xyz.kpxe
	curl -L https://boot.netboot.xyz/ipxe/netboot.xyz-undionly.kpxe -o ${PWD}/internal/resources/tftp/netboot.xyz-undionly.kpxe
	curl -L https://boot.netboot.xyz/ipxe/netboot.xyz.efi -o ${PWD}/internal/resources/tftp/netboot.xyz.efi

gen-build:
	export DEBUG=true
	go run \
		${GO_FLAGS} \
		${PWD}/cmd/server/main \
			generateSchema ${SCHEMA_FILE}
	git add ${SCHEMA_FILE}

gen-proto:
	protoc \
		--proto_path ${PWD} \
		--go_out ${PWD} \
		${PWD}/protobuf/**

gen-clean:
	rm -rf ${PWD}/${GEN_API_TS}/
	rm -rf ${PWD}/${GEN_ED_GO}/
	rm -rf ${PWD}/${GEN_API_GO}/api/
	rm -rf ${PWD}/${GEN_API_GO}/docs/
	rm -rf ${PWD}/${GEN_API_GO}/test/
	rm -rf ${PWD}/${GEN_API_GO}/*.go

gen-tag:
	git add Makefile
	cd ${PWD}
	git commit -m "release version v${VERSION}"
	git tag v${VERSION}
	git push
	git push --tags

gen-client-go:
	docker run \
		--rm -v ${PWD}:/local \
		--user ${UID}:${GID} \
		openapitools/openapi-generator-cli:v7.15.0 generate \
		--git-host beryju.io \
		--git-user-id gravity \
		--git-repo-id api \
		--additional-properties=packageName=api \
		-i /local/schema.yml \
		-g go \
		-o /local/${GEN_API_GO} \
		-c /local/${GEN_API_GO}/config.yaml
	cd ${PWD}/${GEN_API_GO}/
	rm -rf .travis.yml go.mod go.sum test/
	go get
	go fmt ${PWD}/${GEN_API_GO}/
	go tool mvdan.cc/gofumpt -l -w ${PWD}/${GEN_API_GO}/
	git add ${PWD}/${GEN_API_GO}/

gen-client-ts:
	docker run \
		--rm -v ${PWD}:/local \
		--user ${UID}:${GID} \
		openapitools/openapi-generator-cli:v7.15.0 generate \
		-i /local/${SCHEMA_FILE} \
		-g typescript-fetch \
		-o /local/${GEN_API_TS} \
		--additional-properties=typescriptThreePlus=true,supportsES6=true,npmName=gravity-api,npmVersion=${VERSION} \
		--git-user-id BeryJu \
		--git-repo-id gravity
	cd ${PWD}/${GEN_API_TS} && npm i
	\cp -rf ${PWD}/${GEN_API_TS}/* ${PWD}/web/node_modules/gravity-api

gen-client-ts-publish: gen-client-ts
	cd ${PWD}/${GEN_API_TS}
	npm publish
	cd ${PWD}/web
	npm i gravity-api@${VERSION}
	npm version ${VERSION} || true
	git add package*.json

gen-external-dns:
	wget https://kubernetes-sigs.github.io/external-dns/v0.18.0/api/webhook.yaml -O pkg/externaldns/schema.yaml
	docker run \
		--rm -v ${PWD}:/local \
		--user ${UID}:${GID} \
		openapitools/openapi-generator-cli:v7.15.0 generate \
		--git-host beryju.io \
		--git-user-id gravity \
		--git-repo-id api \
		--additional-properties=packageName=externaldnsapi \
		--additional-properties=sourceFolder=externaldnsapi \
		-i /local/pkg/externaldns/schema.yaml \
		-g go-server \
		-o /local/${GEN_ED_GO} \
		-c /local/${GEN_API_GO}/config.yaml
	cd ${PWD}/${GEN_ED_GO}/
	sed -i 's|application/json; charset=UTF-8|application/external.dns.webhook+json;version=1|g' externaldnsapi/routers.go
	rm -f .travis.yml go.mod go.sum main.go Dockerfile
	go fmt ${PWD}/${GEN_ED_GO}/externaldnsapi
	go tool mvdan.cc/gofumpt -l -w ${PWD}/${GEN_ED_GO}/
	git add ${PWD}/${GEN_ED_GO}/

gen-go-tidy:
	go mod tidy
	git add go.mod go.sum

release: gen-build gen-clean gen-client-go gen-external-dns gen-go-tidy gen-client-ts-publish gen-tag

lint: web-lint
	golangci-lint run -v --timeout 5000s

test-env-start:
	cd ${PWD}/hack/tests/
	docker compose --project-name gravity-test-env up -d

test-env-stop:
	cd ${PWD}/hack/tests/
	docker compose --project-name gravity-test-env down -v

install-deps:
	sudo apt-get update
	sudo apt-get install -y nmap libpcap-dev tshark
	sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"
	$(MAKE) internal/resources/macoui internal/resources/blocky internal/resources/tftp

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
		-json \
		${TEST_FLAGS} \
		$(shell go list ./... | grep -v beryju.io/gravity/api | grep -v beryju.io/gravity/cmd | grep -v beryju.io/gravity/pkg/externaldns/generated) \
			2>&1 | tee ${TEST_OUTPUT}
	go tool cover \
		-html ${PWD}/coverage.txt \
		-o ${PWD}/coverage.html

test-e2e-container-build:
	docker build \
		--build-arg=GRAVITY_BUILD_ARGS=GO_BUILD_FLAGS=-cover \
		-t gravity:e2e-test \
		${PWD}

test-e2e: test-e2e-container-build
	export CI="true"
	go test \
		-p 1 \
		-v \
		-coverprofile=${PWD}/coverage.txt \
		-covermode=atomic \
		-count=${TEST_COUNT} \
		-timeout=300s \
		-json \
		-tags=e2e \
		${TEST_FLAGS} \
		beryju.io/gravity/tests \
			2>&1 | tee ${TEST_OUTPUT}
	go tool covdata textfmt \
		-i ${PWD}/tests/coverage/ \
		--pkg $(shell go list ./... | grep -v beryju.io/gravity/api | grep -v beryju.io/gravity/cmd | grep -v beryju.io/gravity/pkg/externaldns/generated | xargs | sed 's/ /,/g') \
		-o ${PWD}/coverage_in_container.txt
	go tool cover \
		-html ${PWD}/coverage_in_container.txt \
		-o ${PWD}/coverage_in_container.html

bench: internal/resources/macoui internal/resources/blocky internal/resources/tftp
	export BOOTSTRAP_ROLES="dns;dhcp;api;discovery;backup;debug;tsdb;tftp"
	export ETCD_ENDPOINT="localhost:2385"
	export LISTEN_ONLY="true"
	export LOG_LEVEL="FATAL"
	export CI="true"
	go test \
		-run=^$$ \
		-bench=^Benchmark \
		-coverprofile=${PWD}/coverage.txt \
		-covermode=atomic \
		-benchmem \
		$(shell go list ./... | grep -v beryju.io/gravity/api | grep -v beryju.io/gravity/cmd | grep -v beryju.io/gravity/pkg/externaldns/generated) \
			2>&1 | tee ${TEST_OUTPUT}
	go tool cover \
		-html ${PWD}/coverage.txt \
		-o ${PWD}/coverage.html

test-convert:
	go tool github.com/jstemmer/go-junit-report/v2 -parser gojson -in ${TEST_OUTPUT} -out ${PWD}/junit.xml
