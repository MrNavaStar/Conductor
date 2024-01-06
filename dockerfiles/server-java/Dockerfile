FROM amazoncorretto:21-alpine
RUN mkdir /conductor

# installing eudev-dev stops minecraft from throwing a warning
RUN apk add jq eudev-dev

RUN adduser -H -D conductor conductor
RUN chown -R conductor:conductor /conductor
USER conductor