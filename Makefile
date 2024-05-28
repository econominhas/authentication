start:
	go run cmd/main.go

migrate:
	migrate -source file://migrations -database postgres://localhost:5432/authentication up

gen-migration:
	./scripts/gen-migration.sh
