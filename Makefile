test:
	go test -v -count=1 ./...
test100:
	go test -v -count=100 ./...
migrateup:
	migrate -path schemes/postgres -database "postgresql://postgres:docker@localhost:5432/try?sslmode=disable" -verbose up
migratedown:
	migrate -path schemes/postgres -database "postgresql://postgres:docker@localhost:5432/try?sslmode=disable" -verbose down
