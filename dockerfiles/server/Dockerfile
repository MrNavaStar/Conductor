FROM gcr.io/distroless/base-debian11:debug as base
FROM busybox:stable-glibc as server

# Binaries for glibc
#COPY --from=base /lib/x86_64-linux-gnu/libdl.so.2 /lib/x86_64-linux-gnu/libdl.so.2
#COPY --from=base /lib/x86_64-linux-gnu/librt.so.1 /lib/x86_64-linux-gnu/librt.so.1
#COPY --from=base /lib/x86_64-linux-gnu/libresolv.so.2 /lib/x86_64-linux-gnu/libresolv.so.2
#COPY --from=base /lib/x86_64-linux-gnu/libm.so.6  /lib/x86_64-linux-gnu/libm.so.6
#COPY --from=base /lib/x86_64-linux-gnu/libpthread.so.0  /lib/x86_64-linux-gnu/libpthread.so.0
#COPY --from=base /lib/x86_64-linux-gnu/libc.so.6  /lib/x86_64-linux-gnu/libc.so.6
#COPY --from=base /lib64/ld-linux-x86-64.so.2  /lib64/ld-linux-x86-64.so.2

# External packages
RUN wget -O /bin/jq https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 && chmod +x /bin/jq

# User Setup
RUN mkdir /conductor
RUN adduser --no-create-home --disabled-password conductor
RUN chown -R conductor /conductor
RUN rm /bin/adduser /bin/chown

# Flatten Image
FROM scratch
COPY --from=server / /
CMD ["/bin/sh"]