FROM golang:1.23 AS build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app /app/cmd/b/main.go

FROM alpine:3.18

ENV PORT_SERVICE_A=:9000
ENV PORT_B=:8000
ENV OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:8889
ENV SERVICE_NAME_A=service_a
ENV SERVICE_NAME_B=service_b
ENV ZIPKIN_URL=http://zipkin:9411/api/v2/spans

WORKDIR /app

COPY --from=build /app/app .
COPY .env .

EXPOSE 8000

CMD ["./app"]