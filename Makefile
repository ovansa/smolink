.PHONY: migrate-up migrate-down migrate-create seed

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-create:
	migrate create -ext sql -dir migrations -seq init

seed:
	migrate -path migrations -database "$(DATABASE_URL)" up 2
