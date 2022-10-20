test:
	go test -v -count=1 ./...
test100:
	go test -v -count=100 ./...
migrateup:
	migrate -path schemes/postgres/up -database "postgresql://postgres:docker@localhost:5432/postgres?sslmode=disable" -verbose up
migratedown:
	migrate -path schemes/postgres/down -database "postgresql://postgres:docker@localhost:5432/postgres?sslmode=disable" -verbose down
run:
	docker compose --env-file .\configs\app.env up --build post-app
