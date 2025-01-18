FROM docker.io/library/ubuntu:24.04

RUN apt-get update && \
    apt-get install -y dnsutils

CMD ["/bin/sleep", "infinity"]
