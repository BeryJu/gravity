FROM docker.io/library/ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends iproute2 isc-dhcp-client tcpdump && \
    apt-get clean

COPY ./dhcp-client/entrypoint.sh /entrypoint.sh

CMD ["/entrypoint.sh"]
