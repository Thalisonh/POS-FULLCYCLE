FROM golang:1.22.3-alpine as build
WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lab-cloud-run main.go

FROM scratch

WORKDIR /app

COPY . .
COPY --from=build /app/lab-cloud-run .
ENTRYPOINT [ "./lab-cloud-run" ]
