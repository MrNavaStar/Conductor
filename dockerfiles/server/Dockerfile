FROM debian:buster-slim
RUN mkdir /conductor

# Packages
RUN apt-get update && apt-get install jq wget xz-utils -y && apt clean && rm -rf /var/lib/apt/lists/*

# User setup
RUN adduser --no-create-home --disabled-password conductor
RUN chown -R conductor:conductor /conductor