FROM mrnavastar/conductor:server
ARG JAVA_VERSION=21

RUN wget --output-document=java.deb https://corretto.aws/downloads/latest/amazon-corretto-${JAVA_VERSION}-x64-linux-jdk.deb

# installing eudev-dev stops minecraft from throwing a warning
RUN apt-get update && apt-get install java-common -y && apt clean && rm -rf /var/lib/apt/lists/*
RUN dpkg --install java.deb
RUN rm -f java.deb