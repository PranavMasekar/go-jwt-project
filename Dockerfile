FROM golang:1.16-alpine

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go mod download

RUN go build -o go-jwt-project .

EXPOSE 8000

CMD ["/app/go-jwt-project"]