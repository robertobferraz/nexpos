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

var CategoryHandlerModule = fx.Module(
	"category_handler",
	fx.Provide(NewCategoryHandler),
)

type CategoryHandler struct {
	logger *zap.SugaredLogger
	authMd *middleware.AuthMiddleware
	userMd *middleware.UserMiddleware
	uc     *usecase.Usecase
}

func NewCategoryHandler(
	logger *zap.SugaredLogger,
	authMd *middleware.AuthMiddleware,
	userMd *middleware.UserMiddleware,
	uc *usecase.Usecase,
) (*CategoryHandler, error) {
	return &CategoryHandler{
		logger: logger,
		authMd: authMd,
		userMd: userMd,
		uc:     uc,
	}, nil
}

func (h *CategoryHandler) RegisterRoutes(r fiber.Router) {
	category := r.Group("/category", h.authMd.Require, h.userMd.CheckUser)
	category.Post("/", h.CreateCategory)
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Creates a new category with optional image upload
// @Tags Category
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param name query string true "Category name"
// @Param description query string false "Category description"
// @Param image formData file false "Category image file"
// @Success 201 {object} dto.Base{data=entity.Category}
// @Failure 400 {object} dto.Base
// @Failure 500 {object} dto.Base
// @Router /category [post]
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	name := utils.PString(c.Query("name"))
	description := utils.PString(c.Query("description"))
	parentID := utils.PString(c.Query("parent_id"))

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
			utils.PInt(1080),
			utils.PInt(1080),
			utils.PInt(1*1024*1024),
			utils.PInt64(5*1024*1024),
			file,
		)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Base{
				Success: utils.PBool(false),
				Error: &dto.BaseError{
					Code:    errors.ErrInvalidInput,
					Message: utils.PString(err.Error()),
				},
			})
		}
	}

	in := dto.CreateCategoryInDto{
		Name:        name,
		Description: description,
		ParentID:    parentID,
	}

	if len(imageData) > 0 {
		in.Image = &dto.Image{
			Name:        imageName,
			Data:        utils.PByte(imageData),
			ContentType: contentType,
		}
	}

	category, err := h.uc.CreateCategory(c.Context(), &in)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    errors.ErrInvalidInput,
				Message: utils.PString(err.Error()),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Base{
		Success: utils.PBool(true),
		Message: utils.PString("category created successfully"),
		Data:    category,
	})
}

// GetCategories godoc
// @Summary List all categories
// @Description Returns a list of all categories with optional parent and image data
// @Tags Category
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.Base{data=[]dto.GetCategoriesOutDto}
// @Failure 500 {object} dto.BaseError
// @Router /category [get]
func (h *CategoryHandler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.uc.GetCategories(c.Context(), &dto.GetCategoriesInDto{
		Protocol: utils.PString(c.Protocol()),
		HostName: utils.PString(c.Hostname()),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    errors.ErrInvalidInput,
				Message: utils.PString(err.Error()),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Base{
		Success: utils.PBool(true),
		Message: utils.PString("categories successfully retrieved"),
		Data:    categories,
	})
}

// FindCategoryByID godoc
// @Summary Find category by ID
// @Description Returns details of a specific category by its ID
// @Tags Category
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Category ID"
// @Success 200 {object} dto.Base{data=dto.GetCategoriesOutDto}
// @Failure 404 {object} dto.BaseError
// @Failure 500 {object} dto.BaseError
// @Router /category/{id} [get]
func (h *CategoryHandler) FindCategoryByID(c *fiber.Ctx) error {
	id := utils.PString(c.Params("id"))

	category, err := h.uc.FindCategory(c.Context(), &dto.FindCategoryInDto{
		ID:       id,
		Protocol: utils.PString(c.Protocol()),
		HostName: utils.PString(c.Hostname()),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    errors.ErrInvalidInput,
				Message: utils.PString(err.Error()),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.Base{
		Success: utils.PBool(true),
		Message: utils.PString("category successfully retrieved"),
		Data:    category,
	})
}
