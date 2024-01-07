#FROM busybox:stable-glibc as server
FROM nimmis/alpine-glibc as server

# External packages
RUN wget -O /bin/jq https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 && chmod +x /bin/jq

# User Setup
RUN mkdir /conductor
RUN adduser --no-create-home --disabled-password conductor
RUN chown -R conductor /conductor
#RUN rm /bin/adduser /bin/chown

# Flatten Image
FROM scratch
COPY --from=server / /