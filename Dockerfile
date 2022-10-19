# Stage 1: Build web
FROM --platform=${BUILDPLATFORM} docker.io/node:19 as web-builder

COPY . /work

ENV NODE_ENV=production
WORKDIR /work
RUN make web-build

# Stage 2: Build
FROM --platform=${BUILDPLATFORM} golang:1.19.2 as builder

ARG GIT_BUILD_HASH
ARG TARGETARCH

ENV GIT_BUILD_HASH=$GIT_BUILD_HASH
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=${TARGETARCH}

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .
COPY --from=web-builder /work/web/dist/ /workspace/web/dist/

RUN make docker-build

# Stage 3: Run
FROM docker.io/library/debian:stable-slim

WORKDIR /

COPY --from=builder /workspace/gravity /bin/gravity
# For debugging purposes
COPY --from=quay.io/coreos/etcd:v3.5.5 /usr/local/bin/etcdctl /usr/bin/etcdctl
COPY --from=quay.io/coreos/etcd:v3.5.5 /usr/local/bin/etcdutl /usr/bin/etcdutl

RUN apt-get update && \
    apt-get install -y --no-install-recommends nmap bash-completion bash ca-certificates && \
    mkdir -p /etc/bash_completion.d/ && \
    rm -rf /tmp/* /var/lib/apt/lists/* /var/tmp/ && \
    mkdir /data && \
    chown 65532:65532 /data && \
    /bin/gravity completion bash > /etc/bash_completion.d/gravity && \
    echo ". /etc/bash_completion" >> /etc/bash.bashrc

ENTRYPOINT ["/bin/gravity"]
CMD [ "server" ]
