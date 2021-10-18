# multi-stage dockerfile for go module

# stage 1 - build golang binary
FROM golang:1.15-buster as builder
ARG MODE
ARG SHIPLOGS

ENV GO111MODULE=on

COPY . /app/

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/benz "main.go"

## stage 2 - use lighter alpine base and expose entry
FROM alpine:latest
ARG MODE
ARG SHIPLOGS

LABEL app="benz"
LABEL version="0.0.1"

COPY --from=builder /app/dist/benz /app/
COPY --from=builder /app/app-dev.env /app/
COPY --from=builder /app/app-prod.env /app/

ENV APP_ENV=$MODE
ENV SHIPLOGS=$SHIPLOGS
ENV BASE_PATH=/app

EXPOSE 8080
CMD /app/benz serve
