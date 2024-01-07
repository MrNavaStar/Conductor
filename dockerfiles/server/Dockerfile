FROM debian:buster-slim
RUN mkdir /conductor

RUN apt update
RUN apt install jq

# User setup
RUN adduser --no-create-home -d conductor conductor
RUN chown -R conductor:conductor /conductor