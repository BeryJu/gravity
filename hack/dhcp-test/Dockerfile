# Same build as the repo's top-level Dockerfile, except the "gravity-api" npm
# package used by the web UI is overlaid with a freshly generated client
# (see /gen-ts-api at the repo root, produced with:
#   docker run --rm -v "$(pwd)":/local openapitools/openapi-generator-cli:v7.15.0 generate \
#     -i /local/schema.yml -g typescript-fetch -o /local/gen-ts-api \
#     --additional-properties=typescriptThreePlus=true,supportsES6=true,npmName=gravity-api,npmVersion=0.33.0 \
#     --git-user-id BeryJu --git-repo-id gravity
# )
# instead of whatever version is currently published to npm. This is needed
# while testing backend fields that haven't been released in gravity-api yet
# (e.g. the DHCP lease `reservation` field) - without this, the web UI would
# silently drop those fields when parsing API responses.
#
# Build context must be the repo root, e.g.:
#   docker build -f hack/dhcp-test/Dockerfile .

# Stage 0: Build the locally generated gravity-api TS client
FROM --platform=${BUILDPLATFORM} docker.io/library/node:26 AS ts-client-builder

WORKDIR /client
COPY ./gen-ts-api/ /client/
RUN npm install && npm run build

# Stage 1: Build web
FROM --platform=${BUILDPLATFORM} docker.io/library/node:26 AS web-builder

ENV NODE_ENV=production

WORKDIR /work

COPY ./Makefile /work/Makefile
COPY ./web/package.json /work/web/package.json
COPY ./web/package-lock.json /work/web/package-lock.json

RUN make web-install

# Overlay the locally built gravity-api client, replacing the published one
COPY --from=ts-client-builder /client/ /work/web/node_modules/gravity-api/

COPY ./web /work/web

RUN make web-build

# Stage 2: Prepare external files
FROM --platform=${BUILDPLATFORM} docker.io/library/ubuntu:26.04 AS downloader

WORKDIR /workspace
COPY Makefile .

RUN apt-get update && \
    apt-get install -y --no-install-recommends make curl ca-certificates tshark && \
    make internal/resources/macoui internal/resources/blocky internal/resources/tftp

# Stage 3: Build
FROM --platform=${BUILDPLATFORM} docker.io/library/golang:1.26 AS builder

ARG GIT_BUILD_HASH
ARG TARGETARCH
ARG GRAVITY_BUILD_ARGS

ENV GIT_BUILD_HASH=$GIT_BUILD_HASH
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=${TARGETARCH}

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .
COPY --from=downloader /workspace/ /workspace/
COPY --from=web-builder /work/web/dist/ /workspace/web/dist/

RUN make ${GRAVITY_BUILD_ARGS} docker-build

# Stage 4: Run
FROM docker.io/library/debian:stable-slim

WORKDIR /

RUN apt-get update && \
    apt-get install -y --no-install-recommends nmap bash-completion bash ca-certificates && \
    mkdir -p /etc/bash_completion.d/ && \
    rm -rf /tmp/* /var/lib/apt/lists/* /var/tmp/ && \
    mkdir /data && \
    chown 65532:65532 /data && \
    echo ". /etc/bash_completion" >> /etc/bash.bashrc

COPY --from=builder /workspace/gravity /bin/gravity

RUN /bin/gravity completion bash > /etc/bash_completion.d/gravity

ENTRYPOINT ["/bin/gravity"]
CMD [ "server" ]
