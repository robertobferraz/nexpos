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

var AddressHandlerModule = fx.Module(
	"address_handler",
	fx.Provide(NewAddressHandler),
)

type AddressHandler struct {
	logger *zap.SugaredLogger
	authMd *middleware.AuthMiddleware
	userMd *middleware.UserMiddleware
	uc     *usecase.Usecase
}

func NewAddressHandler(
	logger *zap.SugaredLogger,
	authMd *middleware.AuthMiddleware,
	userMd *middleware.UserMiddleware,
	uc *usecase.Usecase,
) (*AddressHandler, error) {
	return &AddressHandler{
		logger: logger,
		authMd: authMd,
		userMd: userMd,
		uc:     uc,
	}, nil
}

func (h *AddressHandler) RegisterRoutes(r fiber.Router) {
	address := r.Group("/address", h.authMd.Require, h.userMd.CheckUser)
	address.Post("/:id", h.CreateAddress)
	address.Get("/", h.GetAddress)
}

// GetAddress godoc
// @Summary Get address
// @Description Retrieves the address of the authenticated user
// @Tags Address
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.Base{data=dto.GetAddressOutDto}
// @Failure 500 {object} dto.BaseError
// @Router /address [get]
func (h *AddressHandler) GetAddress(c *fiber.Ctx) error {
	id := c.Locals("userID").(*string)

	resp, err := h.uc.GetAddress(c.Context(), &dto.GetAddressInDto{
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

// CreateAddress godoc
// @Summary Create address
// @Description Creates a new address for the authenticated user
// @Tags Address
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "country ID"
// @Param request body dto.CreateAddressInDto true "Address data. If the selected country is Brazil, you may set action_others to null and only fill in the action_brazil fields."
// @Success 201 {object} dto.Base{data=entity.Address}
// @Failure 400 {object} dto.BaseError
// @Failure 500 {object} dto.BaseError
// @Router /address/{id} [post]
func (h *AddressHandler) CreateAddress(c *fiber.Ctx) error {
	id := c.Locals("userID").(*string)
	countryId := utils.PString(c.Params("id"))
	var in dto.CreateAddressInDto

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
