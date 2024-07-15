postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres
createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb simple_bank

truncateall:
	docker exec -it postgres psql -U root -d simple_bank -c "TRUNCATE TABLE accounts, entries, transactions RESTART IDENTITY;"

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

createnewmigration/%:# you can run createnewmigration/{new_name_schema}
	migrate create -ext sql -dir db/migration -seq $(shell echo $@ | cut -d '/' -f2-)

sqlc:
	sqlc generate

gotest:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown createnewmigration/% sqlc gotest