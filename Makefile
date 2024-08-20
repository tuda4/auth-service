DB_URL=postgresql://root:12345678@localhost:5432/mb_backend?sslmode=disable
postgres:
	docker run --name mb_backend --network mb-backend -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=12345678 -d postgres:15-alpine

createdb:
	docker exec -it mb_backend createdb --username=root --owner=root  mb_backend

dropdb:
	docker exec -it mb_backend dropdb mb_backend

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

test:
	go test -v ./...

server:
	go run main.go

PHONY: postgres createdb dropdb migrateup migratedown sqlc new_migration migrateup1 migratedown1
