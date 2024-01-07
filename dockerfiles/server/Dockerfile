FROM alpine:latest
RUN mkdir /conductor

RUN apk add jq gcompat

RUN adduser -H -D conductor conductor
RUN chown -R conductor:conductor /conductor
USER conductor