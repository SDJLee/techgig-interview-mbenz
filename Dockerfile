# multi-stage dockerfile for go module

# stage 1 - build golang binary
FROM golang:1.15-buster as builder

ENV GO111MODULE=on

COPY . /app/

WORKDIR /app

# unit test case coverage
RUN go test -v -coverpkg=./... -coverprofile=benz.cov ./...
RUN go tool cover -func benz.cov

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o benz "main.go"

## stage 2 - use lighter alpine base and expose entry
FROM alpine:latest

# app metada
LABEL app="benz"
LABEL version="0.0.1"

COPY --from=builder /app/benz /app/
COPY --from=builder /app/app-dev.env /app/
COPY --from=builder /app/app-prod.env /app/

ENV APP_ENV=dev
ENV BASE_PATH=/app

EXPOSE 8080
CMD /app/benz serve
