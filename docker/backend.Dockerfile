FROM debian:buster-slim

ARG GO_VERSION
ARG USER_ID
ARG GROUP_ID
ENV USER_ID ${USER_ID}
ENV GROUP_ID ${GROUP_ID}
ENV GO_VERSION ${GO_VERSION}

RUN set -eux; \
    apt-get update; \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        libssl-dev \
        inotify-tools \
        procps \
        wget \
        ca-certificates \
    ; \
    rm -rf /var/lib/apt/lists/*;

# install go
RUN wget "https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz"
RUN tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz
RUN rm go$GO_VERSION.linux-amd64.tar.gz

ENV PATH=${PATH}:/usr/local/go/bin
ENV GOPATH=/go

RUN mkdir -p $GOPATH
RUN chown $USER_ID:$GROUP_ID $GOPATH

COPY ./docker/certificates/server.pem /etc/ssl/certs

RUN mkdir /.cache
RUN chown ${USER_ID}:${GROUP_ID} /.cache

WORKDIR $GOPATH/src
