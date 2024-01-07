FROM busybox:stable-uclibc as busybox
FROM gcr.io/distroless/base-debian12:latest as server

# Shell
COPY --from=busybox /bin/sh /bin/sh
COPY --from=busybox /bin/kill /bin/kill
COPY --from=busybox /bin/pkill /bin/pkill
COPY --from=busybox /bin/killall /bin/killall
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
COPY --from=busybox /bin/chmod /bin/chmod
# Only used in build steps
COPY --from=busybox /bin/adduser /bin/adduser
COPY --from=busybox /bin/chown /bin/chown

# User Setup
RUN mkdir /conductor
RUN adduser --no-create-home --disabled-password conductor
RUN chown -R conductor /conductor
RUN rm /bin/adduser /bin/chown

# External packages
RUN wget -O /bin/jq https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 && chmod +x /bin/jq

# Flatten Image
FROM scratch
COPY --from=server / /
CMD ["/bin/sh", "-c"]