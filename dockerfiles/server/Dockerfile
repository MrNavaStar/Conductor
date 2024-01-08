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

FROM gcr.io/distroless/base-debian12:latest as bin-provider
FROM busybox:stable-uclibc as server

# Glibc Binaries
COPY --from=bin-provider /lib/x86_64-linux-gnu/ld-2.28.so  /lib/x86_64-linux-gnu/ld-2.28.so
COPY --from=bin-provider /lib/x86_64-linux-gnu/libdl-2.28.so /lib/x86_64-linux-gnu/libdl-2.28.so
COPY --from=bin-provider /lib/x86_64-linux-gnu/librt-2.28.so /lib/x86_64-linux-gnu/librt-2.28.so
COPY --from=bin-provider /lib/x86_64-linux-gnu/libresolv-2.28.so /lib/x86_64-linux-gnu/libresolv-2.28.so
COPY --from=bin-provider /lib/x86_64-linux-gnu/libm-2.28.so /lib/x86_64-linux-gnu/libm-2.28.so
COPY --from=bin-provider /lib/x86_64-linux-gnu/libpthread-2.28.so /lib/x86_64-linux-gnu/libpthread-2.28.so
#COPY --from=bin-provider /lib/x86_64-linux-gnu/libc-2.28.so  /lib/x86_64-linux-gnu/libc-2.28.so

# Glibc symlinks
RUN ln -s /lib/x86_64-linux-gnu/ld-2.28.so /libx86_64-linux-gnu/ld-linux-x86-64.so.2
RUN ln -s /lib/x86_64-linux-gnu/libdl-2.28.so /lib/x86_64-linux-gnu/libdl.so.2
RUN ln -s /lib/x86_64-linux-gnu/librt-2.28.so /lib/x86_64-linux-gnu/librt.so.1
RUN ln -s /lib/x86_64-linux-gnu/libresolv-2.28.so /lib/x86_64-linux-gnu/libresolv.so.2
RUN ln -s /lib/x86_64-linux-gnu/libm-2.28.so /lib/x86_64-linux-gnu/libm.so.6
RUN ln -s /lib/x86_64-linux-gnu/libpthread-2.28.so /lib/x86_64-linux-gnu/libpthread.so.0
#RUN ln -s /lib/x86_64-linux-gnu/libc-2.28.so  /lib/x86_64-linux-gnu/libc.so.6

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