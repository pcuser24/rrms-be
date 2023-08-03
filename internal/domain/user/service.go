package user

import (
	"errors"
	"github.com/user2410/rrms-backend/internal/domain/user/model"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUserByEmail(email string) (*model.UserModel, error)
	Login(email string, password string) (*model.UserModel, error)
	InsertUser(user *model.UserModel) (*model.UserModel, error)
}

type userService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) UserService {
	return &userService{
		repo: repo,
	}
}

func (u *userService) GetUserByEmail(email string) (*model.UserModel, error) {
	return u.repo.GetUserByEmail(email)
}

func (u *userService) Login(email string, password string) (*model.UserModel, error) {
	user, err := u.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user != nil && !checkPassword(password, *user.Password) {
		return nil, errors.New("Invalid password")
	}
	return user, nil
}

func checkPassword(password, hashedPassword string) bool {
	hashedBytes := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(hashedBytes, []byte(password))
	return err == nil
}

func (u *userService) InsertUser(user *model.UserModel) (*model.UserModel, error) {
	return u.repo.InsertUser(user)
}
