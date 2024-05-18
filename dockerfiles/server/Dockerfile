FROM busybox:stable-uclibc as base
FROM gcr.io/distroless/base-debian12:latest as server

COPY --from=base /bin /bin

# External packages
RUN wget -O /tini https://github.com/krallin/tini/releases/download/v0.19.0/tini && chmod +x /tini
RUN wget -O /bin/jq https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 && chmod +x /bin/jq
RUN wget -O /bin/gdmp https://github.com/MrNavaStar/GDMP/releases/download/1.0.0/gdmp && chmod +x /bin/gdmp

# User Setup
RUN mkdir /conductor
RUN adduser --no-create-home --disabled-password conductor
RUN chown -R conductor /conductor

# Flatten Image
FROM scratch
COPY --from=server / /
ENTRYPOINT ["/tini", "--"]