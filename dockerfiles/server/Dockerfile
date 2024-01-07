FROM gcr.io/distroless/base-debian11:debug as base
FROM busybox:stable-glibc as server

# User Setup
RUN mkdir /conductor
RUN adduser --no-create-home --disabled-password conductor
RUN chown -R conductor /conductor
RUN rm /bin/adduser /bin/chown

# External packages
RUN wget -O /bin/jq https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 && chmod +x /bin/jq

# Binaries
COPY --from=base /lib/x86_64-linux-gnu /lib/x86_64-linux-gnu

# Flatten Image
FROM scratch
COPY --from=server / /
CMD ["/bin/sh"]