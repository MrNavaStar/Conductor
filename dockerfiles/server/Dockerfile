FROM debian:buster-slim
RUN mkdir /conductor

RUN apt-get update
RUN apt-get install jq -y

# User setup
RUN adduser --no-create-home -d conductor conductor
RUN chown -R conductor:conductor /conductor