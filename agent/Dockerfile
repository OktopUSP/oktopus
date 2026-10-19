#
# file Dockerfile
#
# Base docker image with all required dependencies for OB-USP-A
#
# Based on Ubuntu 24.04 LTS (Noble Numbat) - long-term support until 2029
# This image includes libmosquitto, libwebsockets and basic compilation tools (automake, autoconf)
#
# One-liner execution line (straightforward build for OB-USP-A execution):
# > docker build -f Dockerfile -t obuspa:latest .
#
# Multi-stage builds execution lines (to tag build stages):
# 1) Create the build environment image:
# > docker build -f Dockerfile -t obuspa:build-env --target build-env .
# 2) Create the OB-USP-A image, then build the application
# > docker build -f Dockerfile -t obuspa:latest --target runner .
#
FROM ubuntu:24.04 AS build-env

# Install dependencies
RUN apt-get update && apt-get -y install \
        libssl-dev \
        libcurl4-openssl-dev\
        libsqlite3-dev \
        libz-dev \
        autoconf \
        automake \
        libtool \
        libmosquitto-dev \
        libwebsockets-dev \
        pkg-config \
        make \
    && apt-get clean

FROM build-env AS runner

ENV MAKE_JOBS=8
ENV OBUSPA_ARGS="-v4"

# Copy in all of the code
# Then compile, as root.
COPY obuspa /obuspa
RUN cd /obuspa/ && \
    autoreconf -fi && \
    ./configure && \
    make -j${MAKE_JOBS} && \
    make install

# Then delete the code
# that's no longer needed
RUN rm -rf /obuspa

# Create database directory and set permissions
RUN mkdir -p /usr/local/var/obuspa && \
    chmod 755 /usr/local/var/obuspa

# Run obuspa with args expanded
CMD obuspa ${OBUSPA_ARGS}
