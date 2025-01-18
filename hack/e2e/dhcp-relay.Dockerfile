FROM docker.io/library/ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends isc-dhcp-relay iproute2 tcpdump && \
    apt-get clean

STOPSIGNAL SIGINT

ENTRYPOINT [ "/usr/sbin/dhcrelay" ]
