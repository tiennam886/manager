FROM golang:1.17

WORKDIR /employee

COPY go.mod .

COPY go.sum .

RUN go mod tidy

COPY . .

RUN go build -o employee cmd/app/main.go

CMD ["./employee", "--env=configs/.env", "serve"]


