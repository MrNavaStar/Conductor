FROM gcr.io/distroless/base-debian11:debug as base
FROM busybox:stable-glibc as server

# Shell
COPY --from=base /busybox/sh /bin/sh
# Commands for templates
COPY --from=base /busybox/mkdir /bin/mkdir
COPY --from=base /busybox/ls /bin/ls
COPY --from=base /busybox/cat /bin/cat
COPY --from=base /busybox/rm /bin/rm
COPY --from=base /busybox/mv /bin/mv
COPY --from=base /busybox/cp /bin/cp
COPY --from=base /busybox/echo /bin/echo
COPY --from=base /busybox/tar /bin/tar
COPY --from=base /busybox/wget /bin/wget
COPY --from=base /busybox/chmod /bin/chmod
# Only used in build steps
COPY --from=base /busybox/adduser /bin/adduser
COPY --from=base /busybox/chown /bin/chown

# User Setup
RUN mkdir /conductor
RUN adduser --no-create-home --disabled-password conductor
RUN chown -R conductor /conductor
RUN rm /bin/adduser /bin/chown

# External packages
RUN wget -O /bin/jq https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 && chmod +x /bin/jq

# Binaries
COPY --from=base /lib/libdl.so.2 /lib/libdl.so.2
COPY --from=base /lib64/libdl.so.2 /lib/64libdl.so.2

# Flatten Image
FROM scratch
COPY --from=server / /
CMD ["/bin/sh"]