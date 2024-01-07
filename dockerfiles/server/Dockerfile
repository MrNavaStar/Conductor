FROM busybox:stable-uclibc as busybox
FROM gcr.io/distroless/static-debian12:latest

# Shell
COPY --from=busybox /bin/sh \
    /bin/mkdir \
    /bin/ls \
    /bin/cat \
    /bin/rm \
    /bin/mv \
    /bin/cp \
    bin/echo \
    /bin/tar \
    /bin/wget \
    /bin/adduser \
    /bin/chown \
    /bin/chmod \
    /bin/

# External packages
RUN wget -O /bin/jq https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 \
    # User Setup
    && mkdir /conductor \
    && adduser --no-create-home --disabled-password conductor \
    && chown -R conductor:conductor /conductor