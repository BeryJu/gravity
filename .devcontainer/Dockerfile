FROM mcr.microsoft.com/devcontainers/go:2-1.25-bookworm

USER root

RUN apt-get update && \
    apt-get install -y --no-install-recommends nmap bash-completion bash ca-certificates && \
    mkdir -p /etc/bash_completion.d/ && \
    rm -rf /tmp/* /var/lib/apt/lists/* /var/tmp/

USER 1000
