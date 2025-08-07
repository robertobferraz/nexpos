package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/robertobff/nexpos/adapter/outbound/auth"
	"github.com/robertobff/nexpos/application/dto"
	"github.com/robertobff/nexpos/application/usecase"
	dErr "github.com/robertobff/nexpos/domain/errors"
	"github.com/robertobff/nexpos/utils"
	"go.uber.org/fx"
)

var UserMiddlewareModule = fx.Module(
	"user_middleware",
	fx.Provide(NewUserMiddleware),
)

type UserMiddleware struct {
	userUc *usecase.Usecase
	fb     *auth.Firebase
}

func NewUserMiddleware(userUc *usecase.Usecase, fb *auth.Firebase) *UserMiddleware {
	return &UserMiddleware{
		userUc: userUc,
		fb:     fb,
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

	user, err := m.userUc.GetUserByUID(c.Context(), &dto.GetUserByUIDInDto{UID: tokenClaims.UserID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    dErr.ErrInternalServer,
				Message: utils.PString(err.Error()),
			},
		})
	}

	if user == nil {
		_, err = m.userUc.CreateUserIfNotExist(
			c.Context(),
			&dto.CreateUserInDto{
				Username:    utils.CreateRandomUsername(tokenClaims.Name),
				PhoneNumber: tokenClaims.PhoneNumber,
				Name:        tokenClaims.Name,
				Email:       tokenClaims.Email,
				Birthdate:   nil,
				ExternalID:  tokenClaims.UserID,
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
	} else {
		if user.DeletedAt != nil {
			user.DeletedAt = nil
			err = m.userUc.SaveUser(c.Context(), user)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(dto.Base{
					Success: utils.PBool(false),
					Error: &dto.BaseError{
						Code:    dErr.ErrInternalServer,
						Message: utils.PString(err.Error()),
					},
				})
			}

			err = m.userUc.CheckUserDeletion(c.Context(), user)
			if err != nil {
				return c.Status(fiber.StatusForbidden).JSON(dto.Base{
					Success: utils.PBool(false),
					Error: &dto.BaseError{
						Code:    dErr.ErrForbidden,
						Message: utils.PString(err.Error()),
					},
				})
			}
		}
	}

	return c.Next()
}
