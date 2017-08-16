#
# This Dockerfile builds the image of e2e 
# binary available in this folder
#
# This Dockerfile may be run with proper ENV pairs
# e.g. docker run --env MY_ENV=something
#

FROM ubuntu:16.04

RUN apt-get update && apt-get install -y \
    iproute2 \
    curl \
    net-tools \
    watch \
 && rm -rf /var/lib/apt/lists/*
    
COPY m-e2e /usr/local/bin/

COPY entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT entrypoint.sh \
