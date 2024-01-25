# Variables:
# DB_URL=postgresql://root:mysecret@localhost:32755/rrms?sslmode=disable

# migration
migratecreate:
	migrate create -ext sql -dir internal/infrastructure/database/migrations -seq ${NAME}

migrateup:
	migrate -path internal/infrastructure/database/migrations -database "${DB_URL}" -verbose up

migrateupn:
	migrate -path internal/infrastructure/database/migrations -database "${DB_URL}" -verbose up ${n}

migratedown:
	migrate -path internal/infrastructure/database/migrations -database "${DB_URL}" -verbose down

migratedownn:
	migrate -path internal/infrastructure/database/migrations -database "${DB_URL}" -verbose down ${n}


sqlcgen:
	sqlc generate


build:
	go build -o rrmsd

# Operations

serve:
	go run main.go serve

# Test db container for local development
test_db:
	docker start rrms_test; \
    if [ $$? -ne 0 ]; then \
        echo "[$$(date)]: Test db not found or failed to start, creating new db container"; \
				docker run --name rrms_test -p 32760:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecret -e POSTGRES_DB=rrms_test -d postgres:15.2; \
    fi

test:
	go test -v -cover -short ./...

test_pkg:
	go test -v -cover github.com/user2410/rrms-backend/${PKG}

.PHONY: sqlcgen build serve migratecreate migrateup migrateup1 migratedown migratedown1 test test_db test_pkg
