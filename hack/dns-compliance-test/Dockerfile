FROM docker.io/library/debian:buster

RUN apt-get update && \
    apt-get install --no-install-recommends -y build-essential autoconf pkg-config libssl-dev autotools-dev automake git ca-certificates && \
    cd /tmp && \
    git clone https://gitlab.isc.org/isc-projects/DNS-Compliance-Testing.git && \
    cd DNS-Compliance-Testing && \
    autoreconf -fi && \
    ./configure && \
    make && \
    make install && \
    apt-get remove build-essential autoconf autotools-dev --purge -y && \
    apt-get autoremove --purge -y && \
    apt-get clean && \
    rm -rf tmp/*

CMD [ "/usr/local/bin/genreport" ]
