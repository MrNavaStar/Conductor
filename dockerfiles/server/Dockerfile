FROM debian:buster-slim
RUN mkdir /conductor

# Packages
RUN apt-get update
RUN apt-get install jq -y
RUN rm -rf /var/lib/apt/lists/*

# User setup
RUN adduser --no-create-home --disabled-password conductor
RUN chown -R conductor:conductor /conductor