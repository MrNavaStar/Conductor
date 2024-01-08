#FROM busybox:stable-uclibc as base
#FROM gcr.io/distroless/base-debian12:latest as server

#COPY --from=base /bin /bin

# External packages
#RUN wget -O /bin/jq https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 && chmod +x /bin/jq

# User Setup
#RUN mkdir /conductor
#RUN adduser --no-create-home --disabled-password conductor
#RUN chown -R conductor /conductor

# Flatten Image
#FROM scratch
#COPY --from=server / /
#CMD ["/bin/sh", "-c", "tail", "-f", "/dev/null"]

FROM busybox:stable-uclibc as bin-provider

ARG GLIBC_VERSION=2.35
RUN wget -O glibc.tar.gz https://github.com/sgerrand/docker-glibc-builder/releases/download/2.35-0/glibc-bin-2.35-0-x86_64.tar.gz && tar -xzf glibc.tar.gz

FROM busybox:stable-uclibc as server

# Glibc Binaries
COPY --from=bin-provider /usr/glibc-compat/ld-linux-x86-64.so.2 /lib/ld-linux-x86-64.so.2
COPY --from=bin-provider /usr/glibc-compat/libdl.so.2 /lib/libdl.so.2
COPY --from=bin-provider /usr/glibc-compat/librt.so.1 /lib/librt.so.1
COPY --from=bin-provider /usr/glibc-compat/libresolv.so.2 /lib/libresolv.so.2
COPY --from=bin-provider /usr/glibc-compat/libm.so.6 /lib/libm.so.6
COPY --from=bin-provider /usr/glibc-compat/libpthread.so.0 /lib/libpthread.so.0
COPY --from=bin-provider /usr/glibc-compat/libc.so.6 /lib/libc.so.6

# External packages
RUN wget -O /bin/jq https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 && chmod +x /bin/jq

# User Setup
RUN mkdir /conductor
RUN adduser --no-create-home --disabled-password conductor
RUN chown -R conductor /conductor

# Flatten Image
FROM scratch
COPY --from=server / /
CMD ["/bin/sh", "-c", "tail", "-f", "/dev/null"]