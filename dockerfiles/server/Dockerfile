#FROM nimmis/alpine-glibc as server

# External packages
#RUN wget -O /bin/jq https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 && chmod +x /bin/jq

# User Setup
#RUN mkdir /conductor
#RUN adduser --no-create-home --disabled-password conductor
#RUN chown -R conductor /conductor
#RUN rm /bin/adduser /bin/chown

# Flatten Image
#FROM scratch
#COPY --from=server / /
#CMD ["/bin/sh", "-c", "tail", "-f", "/dev/null"]

FROM busybox:stable-uclibc as base
FROM gcr.io/distroless/base-debian12:latest as server

COPY --from=base /bin/adduser /bin/cat /bin/chmod /bin/chown /bin/echo /bin/grep /bin/head /bin/ls /bin/mkdir /bin/mv /bin/rm /bin/rmdir /bin/sh /bin/tail /bin/tar /bin/touch /bin/wget /bin/

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
CMD ["/bin/sh", "-c", "tail", "-f", "/dev/null"]