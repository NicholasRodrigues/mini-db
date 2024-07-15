FROM golang:1.22

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o mini-db ./cmd

EXPOSE 8080

CMD ["./mini-db"]