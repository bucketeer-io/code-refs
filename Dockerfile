FROM alpine:3.19

RUN apk update
RUN apk add --no-cache git
RUN apk add --no-cache openssh

COPY bucketeer-find-code-refs /usr/local/bin/bucketeer-find-code-refs

ENTRYPOINT ["bucketeer-find-code-refs"]
