# Stage 1: Build website
FROM --platform=${BUILDPLATFORM} docker.io/node:18 as web-builder

COPY ./web /work/web/

ENV NODE_ENV=production
WORKDIR /work/web
RUN npm ci && npm run build

# Stage 2: Build
FROM golang:1.18 as builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .
COPY --from=web-builder /work/web/dist/ /workspace/web/dist/

RUN CGO_ENABLED=0 GOOS=linux go build -v -a -o gravity .

# Stage 3: Run
FROM docker.io/library/debian:stable-slim

WORKDIR /

COPY --from=builder /workspace/gravity /app/gravity

RUN apt-get update && \
    apt-get install -y --no-install-recommends nmap bash && \
    rm -rf /tmp/* /var/lib/apt/lists/* /var/tmp/ && \
    mkdir /data && \
    chown 65532:65532 /data

USER 65532:65532

ENV INSTANCE_LISTEN=0.0.0.0

ENTRYPOINT ["/app/gravity"]
