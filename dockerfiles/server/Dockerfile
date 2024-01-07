FROM alpine:latest
RUN mkdir /conductor

ARG GLIBC_VERSION=2.34-r0

RUN apk --no-cache add jq ca-certificates openssl

# Compat for games that use glibc (factorio)
RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub
RUN wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/${GLIBC_VERSION}/glibc-${GLIBC_VERSION}.apk
RUN wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/${GLIBC_VERSION}/glibc-bin-${GLIBC_VERSION}.apk
RUN apk --no-cache add glibc-${GLIBC_VERSION}.apk glibc-bin-${GLIBC_VERSION}.apk

# User setup
RUN adduser -H -D conductor conductor
RUN chown -R conductor:conductor /conductor

# Clean up
RUN rm -f glibc-${GLIBC_VERSION}.apk glibc-bin-${GLIBC_VERSION}.apk
RUN apk --no-cache del ca-certificates openssl