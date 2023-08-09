# Stage 1: Build web
FROM --platform=${BUILDPLATFORM} docker.io/library/node:20 as web-builder

WORKDIR /work

COPY ./Makefile /work/Makefile
COPY ./web/package.json /work/web/package.json
COPY ./web/package-lock.json /work/web/package-lock.json

RUN cd web && npm ci

COPY ./web /work/web

ENV NODE_ENV=production
RUN make web-build

# Stage 2: Prepare external files
FROM --platform=${BUILDPLATFORM} docker.io/library/debian:stable-slim as downloader

WORKDIR /workspace
COPY Makefile .

RUN apt-get update && \
    apt-get install -y --no-install-recommends make curl ca-certificates && \
    mkdir -p ./internal/macoui ./internal/blocky && \
    make internal/resources/macoui internal/resources/blocky

# Stage 3: Build
FROM --platform=${BUILDPLATFORM} docker.io/library/golang:1.21.0 as builder

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
COPY --from=downloader /workspace/ /workspace/
COPY --from=web-builder /work/web/dist/ /workspace/web/dist/

RUN make docker-build

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
