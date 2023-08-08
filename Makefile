# Variables:
DB_URL=postgresql://root:mysecret@localhost:32755/rrms?sslmode=disable

# migration
migratecreate:
	migrate create -ext sql -dir internal/infrastructure/database/migrations -seq ${MIGRATION_NAME}

migrateup:
	migrate -path internal/infrastructure/database/migrations -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path internal/infrastructure/database/migrations -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path internal/infrastructure/database/migrations -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path internal/infrastructure/database/migrations -database "$(DB_URL)" -verbose down 1


sqlcgen:
	sqlc generate


build:
	go build -o rrmsd


serve:
	go run main.go serve

.PHONY: sqlcgen build serve migratecreate migrateup migrateup1 migratedown migratedown1