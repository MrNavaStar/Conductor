FROM busybox:stable-uclibc as busybox
FROM gcr.io/distroless/static-debian12:debug

COPY --from=busybox /bin/sh /bin/sh
# Commands for templates
COPY --from=busybox /bin/mkdir /bin/mkdir
COPY --from=busybox /bin/cat /bin/cat
COPY --from=busybox /bin/rm /bin/rm
COPY --from=busybox /bin/echo /bin/echo
COPY --from=busybox /bin/ls /bin/ls
COPY --from=busybox /bin/mv /bin/mv
COPY --from=busybox /bin/cp /bin/cp
COPY --from=busybox /bin/tar /bin/tar

# Packages
RUN apt-get update && apt-get install jq wget xz-utils -y && apt clean && rm -rf /var/lib/apt/lists/*

# User setup
RUN mkdir /conductor
RUN adduser --no-create-home --disabled-password conductor
RUN chown -R conductor:conductor /conductor