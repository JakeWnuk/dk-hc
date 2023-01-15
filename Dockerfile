# Build Go Binary Layer
FROM golang:1.18-buster as go-build-env
RUN mkdir src/app
WORKDIR /src/app
ADD ./main.go .
ADD ./go.mod .
RUN go build .

# Build Hashcat Layer
FROM nvidia/cuda:11.6.0-devel-ubuntu20.04

# Configuration
ENV HASHCAT_VERSION=v6.2.6
ENV FILE_DIR=/files/
ENV CONF_DIR=/etc/config.json

# Installation
RUN apt-get update && apt-get install -y --no-install-recommends \
        ocl-icd-libopencl1 opencl-headers \
        clinfo pkg-config make clinfo build-essential git libcurl4-openssl-dev \
    libssl-dev zlib1g-dev libcurl4-openssl-dev libssl-dev tini; \
    rm -rf /var/lib/apt/lists/*; mkdir -p /files/; \
    git clone https://github.com/hashcat/hashcat.git && cd hashcat && git checkout ${HASHCAT_VERSION} && make install -j4

COPY --from=go-build-env /src/app/dkhashcat /sbin/dkhashcat
COPY ./files/ ${FILE_DIR}
COPY ./config.json ${CONF_DIR}
ENTRYPOINT ["/usr/bin/tini", "--", "/sbin/dkhashcat"]

WORKDIR /data
