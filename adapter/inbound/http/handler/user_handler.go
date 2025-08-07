package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/robertobff/nexpos/adapter/inbound/http/middleware"
	"github.com/robertobff/nexpos/application/dto"
	"github.com/robertobff/nexpos/application/usecase"
	"github.com/robertobff/nexpos/domain/errors"
	"github.com/robertobff/nexpos/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var UserHandlerModule = fx.Module(
	"user_handler",
	fx.Provide(NewUserHandler),
)

type UserHandler struct {
	logger *zap.SugaredLogger
	authMd *middleware.AuthMiddleware
	userMd *middleware.UserMiddleware
	uc     *usecase.Usecase
}

func NewUserHandler(
	logger *zap.SugaredLogger,
	authMd *middleware.AuthMiddleware,
	userMd *middleware.UserMiddleware,
	uc *usecase.Usecase,
) (*UserHandler, error) {
	return &UserHandler{
		logger: logger,
		authMd: authMd,
		userMd: userMd,
		uc:     uc,
	}, nil
}

func (h *UserHandler) RegisterRoutes(r fiber.Router) {
	user := r.Group("/user", h.authMd.Require, h.userMd.CheckUser)
	user.Get("/", h.GetUsers)
}

// GetUsers godoc
// @Summary List all-users
// @Description Returns a list of all users registered in the system
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.Base{data=[]entity.User}
// @Failure 500 {object} dto.BaseError
// @Router /user [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	user, err := h.uc.GetUsers(c.Context())
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
		Message: utils.PString("Success"),
		Data:    user,
	})
}
