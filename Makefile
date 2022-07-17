.SHELLFLAGS += -x -e
PWD = $(shell pwd)
UID = $(shell id -u)
GID = $(shell id -g)

run:
	INSTANCE_LISTEN=0.0.0.0 DEBUG=true LISTEN_ONLY=true go run -v . server

gen-build:
	DEBUG=true go run -v . cli generateSchema schema.yml

gen-clean:
	rm -rf gen-ts-api/

gen-client-web:
	docker run \
		--rm -v ${PWD}:/local \
		--user ${UID}:${GID} \
		openapitools/openapi-generator-cli:v6.0.0 generate \
		-i /local/schema.yml \
		-g typescript-fetch \
		-o /local/gen-ts-api \
		--additional-properties=typescriptThreePlus=true,supportsES6=true,npmName=gravity-api,npmVersion=1.0
	mkdir -p web/node_modules/gravity-api
	cd gen-ts-api && npm i
	\cp -rfv gen-ts-api/* web/node_modules/gravity-api

gen: gen-build gen-clean gen-client-web

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
