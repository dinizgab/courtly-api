run:
	go run cmd/main.go

compose-up:
	docker compose up -d --build

migration-up:
	@goose -dir ./migrations postgres "user=postgres password=booking-pass host=127.0.0.1 port=5432 dbname=booking-local sslmode=disable" up 
