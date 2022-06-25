# Builder
FROM golang:1.18.1-alpine3.15 as builder

# hadolint ignore=DL3017,DL3018
RUN apk update && apk upgrade && \
    apk --update --no-cache add git make

WORKDIR /app

COPY . .

RUN make build

# Distribution
FROM alpine:3.14.0

# hadolint ignore=DL3017,DL3018
RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata

WORKDIR /app

EXPOSE 8080

COPY --from=builder /app/bin/main .

CMD [ "/app/main" ]