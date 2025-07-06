FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN go build -o app .

FROM alpine

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
