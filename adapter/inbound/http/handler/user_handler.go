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
	user.Put("/", h.UpdateUser)

	address := user.Group("/address", h.authMd.Require, h.userMd.CheckUser)
	address.Post("/:id", h.CreateUserAddress)
	address.Get("/", h.GetUserAddress)
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

// GetUserAddress godoc
// @Summary Get user address
// @Description Retrieves the address of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.Base{data=entity.User}
// @Failure 500 {object} dto.BaseError
// @Router /user/address [get]
func (h *UserHandler) GetUserAddress(c *fiber.Ctx) error {
	id := c.Locals("userID").(*string)

	resp, err := h.uc.GetUserAddress(c.Context(), &dto.GetUserAddressInDto{
		Protocol: utils.PString(c.Protocol()),
		HostName: utils.PString(c.Hostname()),
		UserID:   id,
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
		Message: utils.PString("user address!"),
		Data:    resp,
	})
}

// CreateUserAddress godoc
// @Summary Create user address
// @Description Creates a new address for the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "country ID"
// @Param request body dto.CreateUserAddressInDto true "User address data. If the selected country is Brazil, you may set action_others to null and only fill in the action_brazil fields."
// @Success 201 {object} dto.Base{data=entity.UserAddress}
// @Failure 400 {object} dto.BaseError
// @Failure 500 {object} dto.BaseError
// @Router /user/address/{id} [post]
func (h *UserHandler) CreateUserAddress(c *fiber.Ctx) error {
	id := c.Locals("userID").(*string)
	countryId := utils.PString(c.Params("id"))
	var in dto.CreateUserAddressInDto

	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    errors.ErrInvalidInput,
				Message: utils.PString(err.Error()),
			},
		})
	}

	in.UserID = id
	in.CountryID = countryId

	if in.ActionOthers != nil {
		in.ActionOthers.UserID = id
		in.ActionOthers.CountryID = countryId
	}
	if in.ActionBrazil != nil {
		in.ActionBrazil.UserID = id
		in.ActionBrazil.CountryID = countryId
	}

	resp, err := h.uc.CreateUserAddressCondition(c.Context(), &in)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    errors.ErrInternalServer,
				Message: utils.PString(err.Error()),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Base{
		Success: utils.PBool(true),
		Message: utils.PString("user address created with success!"),
		Data:    resp,
	})
}
