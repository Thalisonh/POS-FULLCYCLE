### Contexto
Esta listagem precisa ser feita com:
- Endpoint REST (GET /order) PORT=8080
- Service ListOrders com GRPC PORT=8081
- Query ListOrders GraphQL PORT=8082
Não esqueça de criar as migrações necessárias e o arquivo api.http com a request para criar e listar as orders.

Para a criação do banco de dados, utilize o Docker (Dockerfile / docker-compose.yaml), com isso ao rodar o comando docker compose up tudo deverá subir, preparando o banco de dados.
Inclua um README.md com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.

### Instalação e rodar
docker-compose up --build

### Wire
go generate on cmd

