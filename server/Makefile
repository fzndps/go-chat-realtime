postgres:
	docker exec -it postgres_docker psql

createdb:
	docker exec -it postgres_docker createdb --username=root --owner=root go_chat

dropdb:
	docker exec -it postgres_docker dropdb go_chat

migrateup:
	migrate -path db/migrations -database "postgresql://root:fizonenda150902@localhost:5432/go_chat?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:fizonenda150902@localhost:5432/go_chat?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb