FROM golang:1.17

WORKDIR /team

COPY go.mod .

COPY go.sum .

RUN go mod tidy

COPY . .

RUN go build -o team cmd/app/main.go

CMD ["./team", "--env=configs/.env", "serve"]