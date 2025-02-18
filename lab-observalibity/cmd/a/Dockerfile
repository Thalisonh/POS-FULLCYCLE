FROM golang:1.23 AS build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app /app/cmd/a/main.go

FROM alpine:3.18

ENV PORT_SERVICE_A=:9000
ENV PORT_SERVICE_B=:8000
ENV OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:8889
ENV SERVICE_NAME_A=service_a
ENV SERVICE_NAME_B=service_b
ENV SERVICE_B_URL=http://serviceb:8000/
ENV ZIPKIN_URL=http://zipkin:9411/api/v2/spans

WORKDIR /app

# Copiar o binário compilado e o arquivo .env para a imagem final
COPY --from=build /app/app .
COPY .env .

EXPOSE 9000

CMD ["./app"]