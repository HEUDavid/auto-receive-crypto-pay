FROM alpine:3.18

RUN apk add --no-cache bash curl tar ca-certificates gettext

WORKDIR /app

ARG BIN=receivepay-linux-amd64
ARG VERSION=latest

RUN curl -L -o ${TAR} "https://github.com/HEUDavid/auto-receive-crypto-pay/releases/${VERSION}/download/${BIN}.tar.gz" && \
    tar -xzf ${BIN}.tar.gz && \
    rm ${BIN}.tar.gz && \
    mv ${BIN} receivepay


ENTRYPOINT ["/bin/sh", "-c", "envsubst < conf/conf.toml.tpl > conf/conf.toml && ./receivepay"]
