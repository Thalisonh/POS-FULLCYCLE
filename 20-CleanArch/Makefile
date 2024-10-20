createmigration:
	./migrate create -ext=sql -dir=sql/migrations -seq init

migration: 
	./migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose up

migrationdown: 
	./migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose down

.PHONY: migrate migrationdown createmigration

run:
	go run cmd/main.go cmd/wire_gen.go