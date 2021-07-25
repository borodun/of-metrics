FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

COPY . $GOPATH/src/borodun/collector
WORKDIR $GOPATH/src/borodun/collector/cmd

RUN go get -d -v
RUN go build -o /tmp/collector

FROM alpine

RUN addgroup -S appgroup && adduser -S appuser -G appgroup && mkdir -p /app
COPY --from=builder /tmp/collector /app
RUN chmod a+rx /app/collector

USER appuser
WORKDIR /app

ENV mongo_uri="mongodb://username:password@host:port/defaultauthdb"
ENV LISTENING_PORT 8080

CMD ["./collector"] 