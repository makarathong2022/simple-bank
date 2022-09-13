# migrate -path db/migration -database "postgresql://root:Thong@123@127.0.0.1:5432/BankSystem?sslmode=disable" -verbose up 
# migrate create -ext sql -dir db/migration -seq init_schema (new file migration)

postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=Thong@123 -d 	postgres:12-alpine 

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root BankSystem

dropdb:
	docker exec -it postgres12 dropdb BankSystem

migrateup:
	migrate -path db/migration -database "postgresql://root:Thong@123@127.0.0.1:5432/BankSystem?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:Thong@123@127.0.0.1:5432/BankSystem?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:Thong@123@127.0.0.1:5432/BankSystem?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:Thong@123@127.0.0.1:5432/BankSystem?sslmode=disable" -verbose down 1


sqlc:
	sqlc generate

server:
	go run main.go

mock: 
	mockgen -package mockdb -destination db/mock/store.go  github.com/bank/simple-bank/db/sqlc Store

test: 
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock