# Stage 1: Build
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

RUN make ${GRAVITY_BUILD_ARGS} docker-build-external-dns

# Stage 2: Run
FROM docker.io/library/debian:stable-slim

WORKDIR /

COPY --from=builder /workspace/gravity-external-dns /bin/gravity-external-dns

USER 1000:1000

ENTRYPOINT ["/bin/gravity-external-dns"]
