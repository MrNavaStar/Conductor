FROM amazoncorretto:21-alpine
RUN mkdir /conductor

RUN apk add jq

RUN adduser -H -D conductor conductor
USER conductor