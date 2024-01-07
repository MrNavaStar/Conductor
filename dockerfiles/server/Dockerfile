FROM busybox:stable-uclibc as busybox
FROM gcr.io/distroless/static-debian12:latest

# Shell
COPY --from=busybox /bin/sh /bin/sh
# Commands for templates
COPY --from=busybox /bin/mkdir /bin/mkdir
COPY --from=busybox /bin/ls /bin/ls
COPY --from=busybox /bin/cat /bin/cat
COPY --from=busybox /bin/rm /bin/rm
COPY --from=busybox /bin/mv /bin/mv
COPY --from=busybox /bin/cp /bin/cp
COPY --from=busybox /bin/echo /bin/echo
COPY --from=busybox /bin/tar /bin/tar
COPY --from=busybox /bin/wget /bin/wget

# External packages
RUN wget -O /bin/jq https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64

# User setup
RUN mkdir /conductor
RUN adduser --no-create-home --disabled-password conductor
RUN chown -R conductor:conductor /conductor