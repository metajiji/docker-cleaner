FROM alpine:3.18 AS base

RUN apk add --no-cache --upgrade \
        ca-certificates \
        tzdata

RUN mkdir -v /tmp/rootfs \
    && tar -cf- /usr/share/zoneinfo | tar -xf- -C /tmp/rootfs \
    && mkdir -v /tmp/rootfs/etc \
    && echo 'nobody:x:65534:65534:nobody:/:' > /tmp/rootfs/etc/passwd \
    && echo 'nobody:x:65534:' > /tmp/rootfs/etc/group \
    && install -vDm 644 /etc/ssl/certs/ca-certificates.crt /tmp/rootfs/etc/ssl/certs/ca-certificates.crt

FROM scratch

ARG APP_BIN_NAME=exe

COPY --from=base /tmp/rootfs /
COPY $APP_BIN_NAME /exe

WORKDIR /

USER nobody:nobody

ENTRYPOINT ["/exe"]
