package http

import (
	"github.com/user2410/rrms-backend/internal/domain/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	service auth.Service
}

func NewAdapter(service auth.Service) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) RegisterServer(router *fiber.Router, tokenMaker token.Maker) {
	authRoute := (*router).Group("/auth")

	credentialGroup := authRoute.Group("/credential")
	credentialGroup.Post("/register", a.credentialRegister())
	credentialGroup.Post("/login", GetAuthorizationMiddleware(tokenMaker), a.credentialLogin())
	credentialGroup.Get("/me", AuthorizedMiddleware(tokenMaker), a.credentialGetCurrentUser())
	credentialGroup.Get("/ids", AuthorizedMiddleware(tokenMaker), a.credentialGetUserByIds())
	credentialGroup.Put("/refresh", a.credentialRefresh())
	credentialGroup.Patch("/update", AuthorizedMiddleware(tokenMaker), a.credentialUpdateUser())
	credentialGroup.Delete("/logout", AuthorizedMiddleware(tokenMaker), a.credentialLogout())

	bffGroup := authRoute.Group("/bff")
	bffUserGroup := bffGroup.Group("/user")
	bffUserGroup.Post("/create", a.bffCreateUser())
	bffUserGroup.Get("/get-by-id/:id", a.bffGetUserById())
	bffUserGroup.Get("/get-by-email/:email", a.bffGetUserByEmail())
	bffUserGroup.Get("/get-by-account/:id", a.bffGetUserByAccount())
	bffUserGroup.Patch("/update/:id", a.bffUpdateUser())
	bffUserGroup.Delete("/delete/:id", a.bffDeleteUser())
	bffGroup.Patch("link-account", a.bffLinkAccount())
	bffGroup.Patch("unlink-account", a.bffUnlinkAccount())
	bffSessionGroup := bffGroup.Group("/session")
	bffSessionGroup.Post("/create", a.bffCreateSession())
	bffSessionGroup.Get("/user/:token", a.bffGetSessionAndUser())
	bffSessionGroup.Patch("/update", a.bffUpdateSession())
	bffSessionGroup.Delete("/delete/:token", a.bffDeleteSession())
	bffGroup.Post("/create-verification-token", a.bffCreateVerificationToken())
	bffGroup.Put("/use-verification-token", a.bffUseVerificationToken())
}
