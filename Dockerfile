FROM golang:1.16-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main cmd/webserver/main.go

EXPOSE 5000

CMD ["./main"]