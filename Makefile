.ONESHELL:
.SHELLFLAGS += -x -e
PWD = $(shell pwd)
UID = $(shell id -u)
GID = $(shell id -g)
VERSION = "0.4.5"
LD_FLAGS = -X beryju.io/gravity/pkg/extconfig.Version=${VERSION}
GO_FLAGS = -ldflags "${LD_FLAGS}" -v
SCHEMA_FILE = schema.yml

ci--env:
	echo "sha=${GITHUB_SHA}" >> ${GITHUB_OUTPUT}
	echo "build=${GITHUB_RUN_ID}" >> ${GITHUB_OUTPUT}
	echo "timestamp=$(shell date +%s)" >> ${GITHUB_OUTPUT}
	echo "version=${VERSION}" >> ${GITHUB_OUTPUT}

docker-build:
	go build \
		-ldflags "${LD_FLAGS} -X beryju.io/gravity/pkg/extconfig.BuildHash=${GIT_BUILD_HASH}" \
		-v -a -o gravity .

run:
	export INSTANCE_LISTEN=0.0.0.0
	export DEBUG=true
	export LISTEN_ONLY=true
	go run ${GO_FLAGS} . server

web-build:
	cd web
	npm ci
	npm version ${VERSION} || true
	npm run build

gen-update-oui:
	curl -L https://gitlab.com/wireshark/wireshark/-/raw/master/manuf -o ./internal/macoui/db.txt || true

gen-build:
	DEBUG=true go run ${GO_FLAGS} . generateSchema ${SCHEMA_FILE}
	git add ${SCHEMA_FILE}

gen-clean:
	rm -rf gen-ts-api/

gen-tag:
	cd ${PWD}
	git commit -m "tag version v${VERSION}"
	git tag v${VERSION}

gen-client-go:
	cd ${PWD}/api
	rm *.go
	cd ${PWD}
	docker run \
		--rm -v ${PWD}:/local \
		--user ${UID}:${GID} \
		openapitools/openapi-generator-cli:v6.0.0 generate \
		-i /local/schema.yml \
		-g go \
		-o /local/api \
		-c /local/api/config.yaml
	cd ${PWD}/api/
	rm -f .travis.yml go.mod go.sum
	go get
	go fmt .
	go mod tidy
	git add .

gen-client-ts:
	docker run \
		--rm -v ${PWD}:/local \
		--user ${UID}:${GID} \
		openapitools/openapi-generator-cli:v6.0.0 generate \
		-i /local/${SCHEMA_FILE} \
		-g typescript-fetch \
		-o /local/gen-ts-api \
		--additional-properties=typescriptThreePlus=true,supportsES6=true,npmName=gravity-api,npmVersion=${VERSION} \
		--git-repo-id BeryJu \
		--git-user-id gravity
	cd gen-ts-api && npm i

gen-client-ts-update: gen-client-ts
	cd ${PWD}/gen-ts-api
	npm publish
	cd ${PWD}/web
	npm i gravity-api@${VERSION}
	npm version ${VERSION} || true
	git add package*.json

gen: gen-build gen-clean gen-client-go gen-client-ts-update gen-tag

test-env-start:
	cd hack/tests/
	docker compose --project-name gravity-test-env up -d

test-env-stop:
	cd hack/tests/
	docker compose --project-name gravity-test-env down -v

test:
	export BOOTSTRAP_ROLES="dns;dhcp;api;discovery;backup;debug;tsdb"
	export ETCD_ENDPOINT="localhost:2379"
	export DEBUG="true"
	etcdctl del --prefix / || true
	go test -p 1 -coverprofile=coverage.txt -covermode=atomic -v ./...
	go tool cover -html coverage.txt -o coverage.html
