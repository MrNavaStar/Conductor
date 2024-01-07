FROM mrnavastar/conductor:server
ARG JAVA_VERSION=21

# Install corretto key
RUN wget -O - https://apt.corretto.aws/corretto.key | sudo gpg --dearmor -o /usr/share/keyrings/corretto-keyring.gpg && echo "deb [signed-by=/usr/share/keyrings/corretto-keyring.gpg] https://apt.corretto.aws stable main" | sudo tee /etc/apt/sources.list.d/corretto.list

# installing eudev-dev stops minecraft from throwing a warning
RUN apt-get update && apt-get install java-${JAVA_VERSION}-amazon-corretto-jdk -y && apt clean && rm -rf /var/lib/apt/lists/*