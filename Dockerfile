FROM golang:1.16-alpine As builder

RUN apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o go-jwt-project .

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/go-jwt-project /app/go-jwt-project

EXPOSE 8000

CMD ["/app/go-jwt-project"]