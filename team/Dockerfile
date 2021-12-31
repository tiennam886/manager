#Build stage
FROM golang:1.17-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o team cmd/app/main.go

#Run stage

FROM alpine
WORKDIR /app

COPY --from=builder /app/team .

COPY --from=builder /app/configs/.env .

EXPOSE 8081

CMD ["./team","--env=./.env", "serve"]