FROM debian:buster-slim
RUN mkdir /conductor

# Packages
RUN apt-get update
RUN apt-get install jq -y

# User setup
RUN adduser --no-create-home --disabled-password conductor conductor
RUN chown -R conductor:conductor /conductor