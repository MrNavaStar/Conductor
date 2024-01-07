FROM gcr.io/distroless/base-debian11:debug as base
FROM busybox:stable-glibc as server

# Binaries for glibc
RUN rm -f bin lib lib64
COPY --from=base /bin /bin
COPY --from=base /lib /lib
COPY --from=base /lib64 /lib64

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