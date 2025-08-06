package middleware

import (
	"errors"
	"github.com/robertobff/food-service/adapter/outbound/auth"
	"github.com/robertobff/food-service/application/dto"
	dErr "github.com/robertobff/food-service/domain/errors"
	"github.com/robertobff/food-service/utils"
	"go.uber.org/fx"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
)

var AuthMiddlewareModule = fx.Module(
	"auth_middleware",
	fx.Provide(NewAuthMiddleware),
)

type TokenClaims struct {
	UserID  *string `mapstructure:"user_id"`
	Name    *string `mapstructure:"name"`
	Email   *string `mapstructure:"email"`
	Picture *string `mapstructure:"picture"`
}

type AuthMiddleware struct {
	fbAuth *auth.AuthFirebase
}

func NewAuthMiddleware(fbAuth *auth.AuthFirebase) *AuthMiddleware {
	return &AuthMiddleware{
		fbAuth: fbAuth,
	}
}

func (a *AuthMiddleware) Require(c *fiber.Ctx) error {
	idToken, err := getAuthorization(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    dErr.ErrUnauthorized,
				Message: utils.PString(err.Error()),
			},
		})
	}

	token, err := a.fbAuth.VerifyIDToken(c.Context(), idToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    dErr.ErrUnauthorized,
				Message: nil,
			},
		})
	}

	tokenClaims, err := getClaims(&token.Claims)
	if err != nil {
		err = errors.New("error parsing token")
		return c.Status(fiber.StatusBadRequest).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    dErr.ErrUnauthorized,
				Message: utils.PString(err.Error()),
			},
		})
	}

	c.Locals("tokenClaims", tokenClaims)
	return c.Next()
}

func getAuthorization(c *fiber.Ctx) (*string, error) {
	var idToken string

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header is required")
	}

	if len(strings.Split(authHeader, " ")) != 2 {
		return nil, errors.New("bad authorization header")
	}

	idToken = strings.Split(authHeader, " ")[1]
	return &idToken, nil
}

func getClaims(token *map[string]interface{}) (*TokenClaims, error) {
	var tokenClaims TokenClaims
	if err := mapstructure.Decode(token, &tokenClaims); err != nil {
		return nil, err
	}

	return &tokenClaims, nil
}

func GetTokenClaims(c *fiber.Ctx) (*TokenClaims, error) {
	tokenClaims, ok := c.Locals("tokenClaims").(*TokenClaims)
	if !ok {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    dErr.ErrInternalServer,
				Message: utils.PString("internal server error"),
			},
		})
	}
	return tokenClaims, nil
}
