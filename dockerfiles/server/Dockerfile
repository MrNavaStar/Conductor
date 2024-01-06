FROM alpine:latest
RUN mkdir /conductor

RUN apk add jq

RUN adduser -H -D conductor conductor
RUN chown chown -R conductor:conductor /conductor
USER conductor