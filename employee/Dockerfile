#Build stage
FROM golang:1.17-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o employee cmd/app/main.go

#Run stage

FROM alpine

WORKDIR /app

COPY --from=builder /app/employee .

COPY --from=builder /app/configs/.env .

EXPOSE 8082
CMD ["./employee","--env=./.env", "serve"]


