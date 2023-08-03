package bjwt

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

type BJwt interface {
	GenerateToken(claims jwt.Claims) (*string, error)
	GetJwtHandle() fiber.Handler
}

type bjwt struct {
	secretKey string
}

func NewBjwt(secretKey string) BJwt {
	return &bjwt{
		secretKey: secretKey,
	}
}

func (b *bjwt) GenerateToken(claims jwt.Claims) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(b.secretKey))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (b *bjwt) GetJwtHandle() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(b.secretKey),

		Filter: func(c *fiber.Ctx) bool {
			if c.Method() == "OPTIONS" {
				return true
			}
			return false
		},
	})
}
