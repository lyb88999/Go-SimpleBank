DB_URL_LOCAL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
DB_URL_REMOTE=postgresql://lyb8889999:Lyb1217@@@rm-cn-28t3uyw2y0008wso.rwlb.rds.aliyuncs.com:5432/simple_bank?sslmode=disable
postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL_LOCAL)" -verbose up

migrateupR:
	migrate -path db/migration -database "$(DB_URL_REMOTE)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL_LOCAL)" -verbose up 1

migrateupR1:
	migrate -path db/migration -database "$(DB_URL_REMOTE)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL_LOCAL)" -verbose down

migratedownR:
	migrate -path db/migration -database "$(DB_URL_REMOTE)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL_LOCAL)" -verbose down 1

migratedownR1:
	migrate -path db/migration -database "$(DB_URL_REMOTE)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/lyb88999/Go-SimpleBank/db/sqlc Store

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

.PHONY: postgres createdb dropdb migrateup migrateupR migrateup1 migratedownR1 migratedown migratedownR migratedown1 migratedownR1 sqlc server mock db_docs db_schema