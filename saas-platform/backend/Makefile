DB_URL=postgres://postgres:postgres@localhost:5432/saas?sslmode=disable

run:
	go run cmd/server/main.go

migrate:
	psql -h localhost -U postgres -d saas -f schema.sql

docker-db:
	docker run --name saas-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=saas -d postgres