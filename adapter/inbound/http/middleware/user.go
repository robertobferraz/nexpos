package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/robertobff/food-service/application/dto"
	"github.com/robertobff/food-service/application/usecase"
	dErr "github.com/robertobff/food-service/domain/errors"
	"github.com/robertobff/food-service/utils"
	"go.uber.org/fx"
)

var UserMiddlewareModule = fx.Module(
	"user_middleware",
	AuthMiddlewareModule,
	fx.Provide(NewUserMiddleware),
)

type UserMiddleware struct {
	userUc *usecase.Usecase
}

func NewUserMiddleware(userUc *usecase.Usecase) *UserMiddleware {
	return &UserMiddleware{
		userUc: userUc,
	}
}

func (m *UserMiddleware) CheckUser(c *fiber.Ctx) error {
	tokenClaims, err := GetTokenClaims(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    dErr.ErrInternalServer,
				Message: utils.PString(err.Error()),
			},
		})
	}

	_, err = m.userUc.CreateUser(
		c.Context(),
		&dto.CreateUserInDto{
			Name:       tokenClaims.Name,
			Email:      tokenClaims.Email,
			ExternalID: tokenClaims.UserID,
		},
	)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    dErr.ErrForbidden,
				Message: utils.PString(err.Error()),
			},
		})
	}

	return c.Next()
}
