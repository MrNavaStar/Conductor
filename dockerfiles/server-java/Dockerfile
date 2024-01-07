FROM mrnavastar/conductor:server
ARG JAVA_VERSION=21

# Install amazon corretto
RUN apt-get update && apt-get install java-common -y && apt clean && rm -rf /var/lib/apt/lists/*
RUN wget --output-document=java.deb https://corretto.aws/downloads/latest/amazon-corretto-${JAVA_VERSION}-x64-linux-jdk.deb && dpkg --install java.deb && rm -f java.deb