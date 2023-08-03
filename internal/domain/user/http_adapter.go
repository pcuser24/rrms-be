package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/user2410/rrms-backend/internal/bjwt"
	"github.com/user2410/rrms-backend/internal/domain/user/model"
	appHttp "github.com/user2410/rrms-backend/internal/infrastructure/http"
	"github.com/user2410/rrms-backend/pkg/utils/types"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

type Adapter interface {
	RegisterServer(server appHttp.Server)
}

type adapter struct {
	service UserService
	server  appHttp.Server
	bj      bjwt.BJwt
}

func NewAdapter(service UserService, bj bjwt.BJwt) Adapter {
	return &adapter{
		service: service,
		bj:      bj,
	}
}

func (a *adapter) RegisterServer(server appHttp.Server) {
	a.server = server
	a.server.GetAuthRouter().Post("/register", a.getRegisterHandle())
	a.server.GetAuthRouter().Post("/login", a.getLoginHandle())
	a.server.GetApiRouter().Post("/get-user-by-token", a.getUserFromTokenHandle())
}

func (a *adapter) getRegisterHandle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			log.Println(err.Error())
			return c.SendStatus(http.StatusBadRequest)
		}
		user, err := a.service.GetUserByEmail(payload.Email)
		if err != nil {
			return c.Status(http.StatusUnauthorized).SendString("Wrong username & password " + err.Error())
		}
		if user == nil {
			hash, _ := hashPassword(payload.Password)
			user = &model.UserModel{
				Email:    &payload.Email,
				Password: types.Ptr[string](hash),
			}
			user, err = a.service.InsertUser(user)
			if err != nil {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"code":    http.StatusBadRequest,
					"message": err.Error(),
				})
			}
		} else {
			//Update
		}
		token, err := a.bj.GenerateToken(&jwt.MapClaims{
			"email": payload.Email,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		})
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString("Something went wrong")
		}
		return c.JSON(fiber.Map{
			"email":     user.Email,
			"api_token": token,
		})
	}
}

func (a *adapter) getLoginHandle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			log.Println(err.Error())
			return c.SendStatus(http.StatusBadRequest)
		}
		user, err := a.service.Login(payload.Email, payload.Password)
		if err != nil {
			return c.Status(http.StatusUnauthorized).SendString("something went wrong" + err.Error())
		}
		if user == nil {
			return c.Status(http.StatusUnauthorized).SendString("Wrong username & password")
		}
		token, err := a.bj.GenerateToken(&jwt.MapClaims{
			"email": payload.Email,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		})
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString("Something went wrong" + err.Error())
		}
		return c.JSON(fiber.Map{
			"email":     user.Email,
			"api_token": token,
		})
	}
}

func (a *adapter) getUserFromTokenHandle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		luser := c.Locals("user")
		if luser == nil {
			return c.Status(http.StatusUnauthorized).SendString("Could not find this user")
		}
		user := luser.(*jwt.Token)
		if user == nil {
			return c.Status(http.StatusUnauthorized).SendString("Could not find this user")
		}
		claims := user.Claims.(jwt.MapClaims)
		log.Println(claims)
		//claims := u.bjwt.
		email := claims["email"].(string)
		member, err := a.service.GetUserByEmail(email)
		if err != nil {
			return c.Status(http.StatusUnauthorized).SendString("Could not find this user")
		}
		return c.JSON(member)
	}
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
