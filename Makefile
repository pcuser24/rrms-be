# Variables:
# DB_URL=postgresql://root:mysecret@localhost:32755/rrms?sslmode=disable

# Migration
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


# SQLC
sqlcgen:
	sqlc generate


# Operations
build:
	go build -o rrmsd
serve:
	go run main.go serve


# Mocking:
mock_repo:
	@directory_path="internal/domain"; \
	subdirectories=$$(find "$$directory_path" -mindepth 1 -maxdepth 1 -type d); \
	for subdir in $$subdirectories; do \
		mockgen -package repo -destination $$subdir/repo/mock.go github.com/user2410/rrms-backend/$$subdir/repo Repo; \
	done; \
	true;

mock_asynctask:
	@directory_path="internal/domain"; \
	subdirectories=$$(find "$$directory_path" -mindepth 1 -maxdepth 1 -type d); \
	for subdir in $$subdirectories; do \
		mockgen -package asynctask -destination $$subdir/asynctask/mock.go github.com/user2410/rrms-backend/$$subdir/asynctask TaskDistributor,TaskProcessor; \
	done; \
	true;


# Test
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


# Localstack: for development and testing only
# Start s3
ls_s3:
	localstack start -d && \
	awslocal s3api create-bucket --bucket rrms-image --region ap-southeast-1 --create-bucket-configuration LocationConstraint=ap-southeast-1 && \
	awslocal s3api put-bucket-cors --bucket rrms-image --cors-configuration file://$$PWD/internal/infrastructure/aws/s3/cors-config.json

.PHONY: sqlcgen build serve migratecreate migrateup migrateup1 migratedown migratedown1 mock_repo mock_asynctask test test_db test_pkg ls_s3
