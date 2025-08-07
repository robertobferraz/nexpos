package handler

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/robertobff/nexpos/adapter/inbound/http/middleware"
	"github.com/robertobff/nexpos/application/dto"
	"github.com/robertobff/nexpos/application/usecase"
	"github.com/robertobff/nexpos/domain/errors"
	"github.com/robertobff/nexpos/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var AuthHandlerModule = fx.Module(
	"auth_handler",
	fx.Provide(NewAuthHandler),
)
var validate = validator.New()

type AuthHandler struct {
	logger *zap.SugaredLogger
	authMd *middleware.AuthMiddleware
	userMd *middleware.UserMiddleware
	uc     *usecase.Usecase
}

func NewAuthHandler(
	logger *zap.SugaredLogger,
	authMd *middleware.AuthMiddleware,
	userMd *middleware.UserMiddleware,
	uc *usecase.Usecase,
) (*AuthHandler, error) {
	return &AuthHandler{
		logger: logger,
		uc:     uc,
		authMd: authMd,
		userMd: userMd,
	}, nil
}

func (h *AuthHandler) RegisterRoutes(r fiber.Router) {
	noAuth := r.Group("/auth")
	noAuth.Post("signup", h.Signup)

	auth := r.Group("/auth", h.authMd.Require, h.userMd.CheckUser)
	auth.Delete("disable/:id", h.Delete)
}

func formatValidationErrors(err error) string {
	var sb strings.Builder
	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()

		sb.WriteString("[")
		sb.WriteString(field)
		sb.WriteString(": ")
		switch tag {
		case "required":
			sb.WriteString("it's mandatory")
		case "email":
			sb.WriteString("must be a valid email")
		case "min":
			sb.WriteString("min value not reached")
		case "max":
			sb.WriteString("max value not reached")
		default:
			sb.WriteString("invalid")
		}
		sb.WriteString("], ")
	}
	return sb.String()
}

// Signup godoc
// @Summary      Create a new user
// @Description  Endpoint to register a new user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.CreateUserInDto  true  "User registration input"
// @Success      200      {object}  dto.Base{data=dto.CreateUserOutDto}
// @Failure      500      {object}  dto.Base
// @Router       /auth/signup [post]
func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	var in *dto.CreateUserInDto
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    errors.ErrInvalidInput,
				Message: utils.PString(err.Error()),
			},
		})
	}

	if err := validate.Struct(in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    errors.ErrInvalidInput,
				Message: utils.PString("Validation: " + formatValidationErrors(err)),
			},
		})
	}

	userOut, err := h.uc.CreateUser(c.Context(), in)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    errors.ErrInternalServer,
				Message: utils.PString(err.Error()),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Base{
		Success: utils.PBool(true),
		Error:   nil,
		Message: utils.PString("user created successfully"),
		Data:    userOut,
	})
}

// Delete godoc
// @Summary Delete a user
// @Description Endpoint to delete a user by ID
// @Security ApiKeyAuth
// @Tags Auth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.Base
// @Failure 500 {object} dto.Base
// @Router /auth/disable/{id} [delete]
func (h *AuthHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.uc.DeleteUser(c.Context(), &dto.DeleteUserInDto{ID: utils.PString(id)})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    errors.ErrInternalServer,
				Message: utils.PString(err.Error()),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Base{
		Success: utils.PBool(true),
		Error:   nil,
		Message: utils.PString("user deleted successfully"),
	})
}
