.SHELLFLAGS += -x -e
PWD = $(shell pwd)
UID = $(shell id -u)
GID = $(shell id -g)

run-debug:
	DEBUG=true go run -v .

gen-build:
	DEBUG=true go run -v . generateSchema schema.yml

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

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
