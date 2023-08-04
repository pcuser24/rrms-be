package unit

import db "github.com/user2410/rrms-backend/internal/infrastructure/database"

type Repo interface{}

type repo struct {
	dao db.DAO
}

func NewRepo(d db.DAO) Repo {
	return &repo{
		dao: d,
	}
}
