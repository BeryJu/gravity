FROM docker.io/library/ubuntu:24.04

RUN apt-get update && \
    apt-get install -y iproute2 isc-dhcp-client tcpdump

COPY ./dhcp-container/entrypoint.sh /entrypoint.sh

CMD ["/entrypoint.sh"]
