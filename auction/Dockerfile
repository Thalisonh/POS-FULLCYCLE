FROM golang:1.24 AS build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app /app/cmd/auction/main.go

FROM alpine:3.18

ENV MONGODB_URL=mongodb://localhost:27017
ENV MONGODB_DB=auctions
ENV BATCH_INSERT_INTERVAL=7m
ENV MAX_BATCH_SIZE=10
ENV AUCTION_INTERVAL=10m

WORKDIR /app

COPY --from=build /app/app .
COPY .env .

EXPOSE 8080

CMD ["./app"]