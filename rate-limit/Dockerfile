FROM golang:1.24 AS build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app /app/cmd/limit/main.go

FROM alpine:3.18

WORKDIR /app

# Copiar o binário compilado e o arquivo .env para a imagem final
COPY --from=build /app/app .
COPY .env .

EXPOSE 8080

CMD ["./app"]