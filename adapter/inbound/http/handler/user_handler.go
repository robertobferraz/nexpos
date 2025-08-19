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
	user.Get("/:id", h.FindUserByID)
	user.Put("/", h.UpdateUser)
}

// GetUsers godoc
// @Summary List all-users
// @Description Returns a list of all users registered in the system
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.Base{data=[]dto.GetUsersOutDto}
// @Failure 500 {object} dto.BaseError
// @Router /user [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	user, err := h.uc.GetUsers(c.Context(), &dto.GetUsersInDto{
		Protocol: utils.PString(c.Protocol()),
		HostName: utils.PString(c.Hostname()),
	})
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

// FindUserByID godoc
// @Summary Find user by ID
// @Description Retrieves the details of a specific user by their ID
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} dto.Base{data=dto.FindUserOutDto}
// @Failure 404 {object} dto.BaseError
// @Failure 500 {object} dto.BaseError
// @Router /user/{id} [get]
func (h *UserHandler) FindUserByID(c *fiber.Ctx) error {
	id := utils.PString(c.Params("id"))

	user, err := h.uc.FindUserByID(c.Context(), &dto.FindUserInDto{
		ID:       id,
		Protocol: utils.PString(c.Protocol()),
		HostName: utils.PString(c.Hostname()),
	})

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

// UpdateUser godoc
// @Summary Update user
// @Description Updates the details of an existing user, including optional profile image
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param email query string false "Email address"
// @Param username query string false "Username"
// @Param phone_number query string false "Phone number"
// @Param birth_date query string false "Birth date in format YYYY-MM-DD"
// @Param name query string false "Full name"
// @Param image formData file false "Profile image"
// @Security ApiKeyAuth
// @Success 200 {object} dto.Base{data=entity.User}
// @Failure 400 {object} dto.BaseError
// @Failure 500 {object} dto.BaseError
// @Router /user [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Locals("userID").(*string)

	email := utils.PString(c.Query("email"))
	username := utils.PString(c.Query("username"))
	phoneNumber := utils.PString(c.Query("phone_number"))
	birthDate := utils.PString(c.Query("birth_date"))
	name := utils.PString(c.Query("name"))

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    errors.ErrInvalidInput,
				Message: utils.PString(err.Error()),
			},
		})
	}

	var imageData []byte
	var imageName *string
	var contentType *string

	if file != nil {
		imageData, imageName, contentType, err = utils.UploadImage(
			utils.PInt(400),
			utils.PInt(400),
			utils.PInt(500*1024),
			utils.PInt64(5*1024*1024),
			file,
		)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Base{
				Success: utils.PBool(false),
				Error: &dto.BaseError{
					Code:    errors.ErrInternalServer,
					Message: utils.PString(err.Error()),
				},
			})
		}
	}

	in := dto.SaveUserInDto{
		ID:          id,
		Email:       email,
		Username:    username,
		PhoneNumber: phoneNumber,
		Name:        name,
		Birthdate:   birthDate,
	}

	if len(imageData) > 0 {
		in.Image = &dto.Image{
			Name:        imageName,
			Data:        utils.PByte(imageData),
			ContentType: contentType,
		}
	}

	user, err := h.uc.SaveUser(c.Context(), &in)
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
		Message: utils.PString("User updated with success!"),
		Data:    user,
	})
}
