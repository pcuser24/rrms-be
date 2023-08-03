package queries

import _ "embed"

var (
	//go:embed get_user_by_email.sql
	GetUserByEmailSQL string
	//go:embed login.sql
	LoginSQL string

	//go:embed insert_user.sql
	InsertUserSQL string
)
