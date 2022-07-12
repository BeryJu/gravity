# Stage 1: Build website
FROM --platform=${BUILDPLATFORM} docker.io/node:18 as web-builder

COPY ./web /work/web/

ENV NODE_ENV=production
WORKDIR /work/web
RUN npm i && npm run build

# Stage 2: Build
FROM golang:1.18 as builder

ARG GIT_BUILD_HASH
ENV GIT_BUILD_HASH=$GIT_BUILD_HASH
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .
COPY --from=web-builder /work/web/dist/ /workspace/web/dist/

RUN go build \
    -ldflags "-X beryju.io/gravity/pkg/extconfig.BuildHash=$GIT_BUILD_HASH" \
    -v -a -o gravity .

# Stage 3: Run
FROM docker.io/library/debian:stable-slim

WORKDIR /

COPY --from=builder /workspace/gravity /app/gravity

RUN apt-get update && \
    apt-get install -y --no-install-recommends nmap bash && \
    rm -rf /tmp/* /var/lib/apt/lists/* /var/tmp/ && \
    mkdir /data && \
    chown 65532:65532 /data

# For debugging purposes
COPY --from=quay.io/coreos/etcd:v3.5.4 /usr/local/bin/etcdctl /usr/bin/etcdctl

USER 65532:65532

ENTRYPOINT ["/app/gravity"]
