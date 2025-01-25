go run cmd/ordersystem/main.go cmd/ordersystem/wire_gen.go

mysql -uroot -p orders

cd cmd/
wire


sudo docker-compose build
docker-compose up -d