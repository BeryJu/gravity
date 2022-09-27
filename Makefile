.SHELLFLAGS += -x -e
PWD = $(shell pwd)
UID = $(shell id -u)
GID = $(shell id -g)
VERSION = "0.1.18"
LD_FLAGS = -X beryju.io/gravity/pkg/extconfig.Version=${VERSION}
GO_FLAGS = -ldflags "${LD_FLAGS}" -v
SCHEMA_FILE = schema.yml

ci--env:
	echo "##[set-output name=sha]${GITHUB_SHA}"
	echo "##[set-output name=build]${GITHUB_RUN_ID}"
	echo "##[set-output name=timestamp]$(shell date +%s)"
	echo "##[set-output name=version]${VERSION}"

docker-build:
	go build \
		-ldflags "${LD_FLAGS} -X beryju.io/gravity/pkg/extconfig.BuildHash=${GIT_BUILD_HASH}" \
		-v -a -o gravity .

run:
	INSTANCE_LISTEN=0.0.0.0 DEBUG=true LISTEN_ONLY=true go run ${GO_FLAGS} . server

web-build:
	cd web && npm ci
	cd web && npm version ${VERSION} || true
	cd web && npm run build

gen-build:
	DEBUG=true go run ${GO_FLAGS} . cli generateSchema ${SCHEMA_FILE}
	git add ${SCHEMA_FILE}

gen-clean:
	rm -rf gen-ts-api/

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
	cd gen-ts-api && npm publish
	cd web && npm i gravity-api@${VERSION}
	cd web && git add package*.json
	cd web && npm version ${VERSION} || true

gen: gen-build gen-clean gen-client-ts-update

test-etcd-start:
	docker run \
		-d --rm \
		-p 2379:2379 \
		--name gravity-test-etcd \
		quay.io/coreos/etcd:v3.5.4 \
		/usr/local/bin/etcd \
		--listen-client-urls http://0.0.0.0:2379 \
		--advertise-client-urls http://127.0.0.1:2379

test-etcd-stop:
	docker stop gravity-test-etcd || true
	docker kill gravity-test-etcd || true
	docker rm gravity-test-etcd || true

test:
	export BOOTSTRAP_ROLES="dns;dhcp;api;discovery;backup"
	export ETCD_ENDPOINT="localhost:2379"
	go test -race -coverprofile=coverage.txt -covermode=atomic -v ./...
