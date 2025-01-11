FROM docker.io/library/debian:bookworm

RUN apt-get update && \
    apt-get install --no-install-recommends -y curl ca-certificates gpg git xz-utils && \
    echo 'deb http://download.opensuse.org/repositories/network:/dns-oarc/Debian_12/ /' | tee /etc/apt/sources.list.d/network:dns-oarc.list && \
    curl -fsSL https://download.opensuse.org/repositories/network:dns-oarc/Debian_12/Release.key | gpg --dearmor | tee /etc/apt/trusted.gpg.d/network_dns-oarc.gpg > /dev/null && \
    mkdir -p /tmp/files && \
    git clone https://github.com/DNS-OARC/sample-query-data /tmp/files && \
    unxz /tmp/files/*.xz && \
    cat /tmp/files/queryfile* > /queryfile && \
    rm -rf /tmp/files && \
    apt-get update && \
    apt-get install -y dnsperf && \
    apt-get remove --purge -y curl ca-certificates gpg git xz-utils && \
    apt-get clean

ENTRYPOINT [ "/usr/bin/dnsperf", "-d", "/queryfile" ]
