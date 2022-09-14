# migrate -path db/migration -database "postgresql://root:Thong@123@127.0.0.1:5432/BankSystem?sslmode=disable" -verbose up 
# migrate create -ext sql -dir db/migration -seq init_schema (new file migration)

postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=Vk7dYrptO5AMuHKQUydk -d 	postgres:12-alpine 

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:Vk7dYrptO5AMuHKQUydk@simple-bank.ci5bs4qfk9kq.ap-southeast-1.rds.amazonaws.com:5432/simple_bank" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:Vk7dYrptO5AMuHKQUydk@simple-bank.ci5bs4qfk9kq.ap-southeast-1.rds.amazonaws.com:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:Vk7dYrptO5AMuHKQUydk@simple-bank.ci5bs4qfk9kq.ap-southeast-1.rds.amazonaws.com:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:Vk7dYrptO5AMuHKQUydk@simple-bank.ci5bs4qfk9kq.ap-southeast-1.rds.amazonaws.com:5432/simple_bank?sslmode=disable" -verbose down 1


sqlc:
	sqlc generate

server:
	go run main.go

mock: 
	mockgen -package mockdb -destination db/mock/store.go  github.com/bank/simple-bank/db/sqlc Store

test: 
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock