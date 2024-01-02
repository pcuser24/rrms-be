# Variables:
DB_URL=postgresql://root:mysecret@localhost:32755/rrms?sslmode=disable

# migration
migratecreate:
	migrate create -ext sql -dir internal/infrastructure/database/migrations -seq ${NAME}

migrateup:
	migrate -path internal/infrastructure/database/migrations -database "$(DB_URL)" -verbose up

migrateupn:
	migrate -path internal/infrastructure/database/migrations -database "$(DB_URL)" -verbose up ${n}

migratedown:
	migrate -path internal/infrastructure/database/migrations -database "$(DB_URL)" -verbose down

migratedownn:
	migrate -path internal/infrastructure/database/migrations -database "$(DB_URL)" -verbose down ${n}


sqlcgen:
	sqlc generate


build:
	go build -o rrmsd


serve:
	go run main.go serve

.PHONY: sqlcgen build serve migratecreate migrateup migrateup1 migratedown migratedown1
