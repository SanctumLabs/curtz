# Builder
FROM golang:1.18.1-alpine3.15 as builder

# hadolint ignore=DL3017,DL3018
RUN apk update && apk upgrade && \
    apk --update --no-cache add git make

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o curtz app/cmd/main.go

# Distribution
FROM alpine:3.15.4

# hadolint ignore=DL3017,DL3018
RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata

WORKDIR /app

EXPOSE 8085

COPY --from=builder /app/curtz .

CMD [ "/app/curtz" ]