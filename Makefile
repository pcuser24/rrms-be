# Variables:
# DB_URL=postgresql://root:mysecret@localhost:32755/rrms?sslmode=disable

TESTDB_CONTAINER=rrms_test
LSS3_CONTAINER=rrms_s3
LSS3_BUCKET=rrms-image

# Migration
migratecreate:
	@{ \
    set -e ;\
    output=$$(migrate create -ext sql -dir internal/infrastructure/database/migrations -seq ${NAME} 2>&1) ;\
    echo "$$output" | while read -r file ; do \
      echo "BEGIN;\n\nEND;" > $$file ;\
    done ;\
    }

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
./rrmsd: build

build:
	go build -o rrmsd .

serve: ./rrmsd
	./rrmsd serve

payment:
	go run main.go payment

dev:
	air serve

# Tunnel for local development, use localtunnel (`npm install -g localtunnel`)
tunnel:
	lt --port 8000 --subdomain rrms

# Mocking:
mock_repo:
	mockgen -package repo -destination internal/domain/${DOMAIN}/repo/mock.go github.com/user2410/rrms-backend/internal/domain/${DOMAIN}/repo Repo; \

mock_repos:
	@directory_path="internal/domain"; \
	subdirectories=$$(find "$$directory_path" -mindepth 1 -maxdepth 1 -type d); \
	for subdir in $$subdirectories; do \
		mockgen -package repo -destination $$subdir/repo/mock.go github.com/user2410/rrms-backend/$$subdir/repo Repo; \
	done; \
	true;

# Test
# Test db container for local development
test_db:
	docker start $(TESTDB_CONTAINER); \
	if [ $$? -ne 0 ]; then \
		echo "[$$(date)]: Test db not found or failed to start, creating new db container"; \
		docker run --name $(TESTDB_CONTAINER) -p 32760:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecret -e POSTGRES_DB=$(TESTDB_CONTAINER) -d postgres:15.2; \
	fi

test:
	go test -v -cover -short ./...

test_pkg:
	go test -v -cover github.com/user2410/rrms-backend/${PKG}


# Localstack: for development and testing only
# Start s3
ls_s3:
	docker start $(LSS3_CONTAINER); \
	if [ $$? -ne 0 ]; then \
		echo "[$$(date)]: Localstack s3 not found or failed to start, creating new localstack s3 container" && \
		docker run -d -p 4566:4566 --name $(LSS3_CONTAINER) localstack/localstack:s3-latest && \
		sleep 5 && \
		awslocal s3api create-bucket --bucket $(LSS3_BUCKET) --region ap-southeast-1 --create-bucket-configuration LocationConstraint=ap-southeast-1 && \
		awslocal s3api put-bucket-cors --bucket $(LSS3_BUCKET) --cors-configuration file://$$PWD/internal/infrastructure/aws/s3/cors-config.json; \
	fi


.PHONY: sqlcgen build serve tunnel dev migratecreate migrateup migrateup1 migratedown migratedown1 mock_repo mock_asynctask test test_db test_pkg ls_s3
