FROM golang:1.22-alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o employee-api

EXPOSE 4000

ENTRYPOINT ["/app/employee-api","migrate"]
ENTRYPOINT ["/app/employee-api","server"]
