INSERT INTO users (email, password, created_at, updated_at)
VALUES (:email, :password, NOW(), NOW())
RETURNING *