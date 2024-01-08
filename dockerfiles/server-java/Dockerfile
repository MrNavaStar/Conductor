FROM debian:buster-slim as base
ARG JAVA_VERSION=21.0.1.12-1

# Install minimal amazon corretto - https://github.com/corretto/corretto-docker/blob/main/21/slim/debian/Dockerfile
RUN set -ux \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
        curl ca-certificates gnupg software-properties-common fontconfig \
    && curl -fL https://apt.corretto.aws/corretto.key | apt-key add - \
    && add-apt-repository 'deb https://apt.corretto.aws stable main' \
    && mkdir -p /usr/share/man/man1 \
    && apt-get update \
    && apt-get install -y java-21-amazon-corretto-jdk=1:$JAVA_VERSION binutils \
    && jlink --add-modules "$(java --list-modules | sed -e 's/@[0-9].*$/,/' | tr -d \\n)" --no-man-pages --no-header-files --strip-debug --output /opt/corretto-slim \
    && apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false \
            curl gnupg software-properties-common binutils java-21-amazon-corretto-jdk=1:$JAVA_VERSION

FROM mrnavastar/conductor:server as server-java
# Binaries
COPY --from=base /lib/x86_64-linux-gnu/libz.so.1.2.11 /lib/x86_64-linux-gnu/libz.so.1.2.11
COPY --from=base /lib/x86_64-linux-gnu/libudev.so.1.6.13 /lib/x86_64-linux-gnu/libudev.so.1.6.13
# Create symlinks
RUN ln -s /lib/x86_64-linux-gnu/libz.so.1.2.11 /lib/x86_64-linux-gnu/libz.so.1
RUN ln -s /lib/x86_64-linux-gnu/libudev.so.1.6.13 /lib/x86_64-linux-gnu/libudev.so.1

# Java
#RUN mkdir -p /usr/lib/jvm
#COPY --from=base /opt/corretto-slim /usr/lib/jvm/java-21-amazon-corretto
#RUN jdk_tools="java keytool rmid rmiregistry javac jaotc jlink jmod jhsdb jar jarsigner javadoc javap jcmd jconsole jdb jdeps jdeprscan jimage jinfo jmap jps jrunscript jshell jstack jstat jstatd serialver" \
#    && priority=$(echo "1${JAVA_VERSION}" | sed "s/\(\.\|-\)//g") \
#    && for i in ${jdk_tools}; do \
#          ln -s /usr/lib/jvm/java-21-amazon-corretto/bin/$i /usr/bin/$i; \
#       done

COPY --from=base /opt/corretto-slim/bin/java /usr/lib/jvm/java-21-amazon-corretto/bin/java
COPY --from=base /opt/corretto-slim/lib/libjli.so /usr/lib/jvm/java-21-amazon-corretto/lib/libjli.so
COPY --from=base /opt/corretto-slim/lib/libjava.so /usr/lib/jvm/java-21-amazon-corretto/lib/libjava.so
COPY --from=base /opt/corretto-slim/lib/server/libjvm.so /usr/lib/jvm/java-21-amazon-corretto/lib/server/libjvm.so
COPY --from=base /opt/corretto-slim/lib/jvm.cfg /usr/lib/jvm/java-21-amazon-corretto/lib/jvm.cfg

ENV LANG C.UTF-8
ENV JAVA_HOME=/usr/lib/jvm/java-21-amazon-corretto

# Flatten Image
FROM scratch
COPY --from=server-java / /
CMD ["/bin/sh", "-c", "tail", "-f", "/dev/null"]