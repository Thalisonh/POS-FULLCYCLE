FROM golang:1.23 AS build

# Definir o diretório de trabalho
WORKDIR /app

# Copiar os arquivos do projeto para dentro do contêiner
COPY . .

# Baixar dependências e compilar a aplicação
#RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o app /app/cmd/b/main.go

# Imagem final
FROM alpine:3.18

# Definir variáveis de ambiente
ENV DB_HOST=mysql
ENV DB_PORT=3306
ENV DB_USER=root
ENV DB_PASSWORD=root
ENV DB_DRIVER=mysql
ENV DB_NAME=orders
ENV WEB_SERVER_PORT=8080
ENV GRPC_SERVER_PORT=8081
ENV GRAPHQL_SERVER_PORT=8082



# Definir diretório de trabalho
WORKDIR /app

# Copiar o binário compilado e o arquivo .env para a imagem final
COPY --from=build /app/app .
COPY .env .

# Expor a porta onde a aplicação será executada
EXPOSE 8080
EXPOSE 8082
EXPOSE 8083

# Comando para executar a aplicação
CMD ["./app"]