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

var ImageHandlerModule = fx.Module(
	"image_handler",
	fx.Provide(NewImageHandler),
)

type ImageHandler struct {
	logger *zap.SugaredLogger
	authMd *middleware.AuthMiddleware
	userMd *middleware.UserMiddleware
	uc     *usecase.Usecase
}

func NewImageHandler(
	logger *zap.SugaredLogger,
	authMd *middleware.AuthMiddleware,
	userMd *middleware.UserMiddleware,
	uc *usecase.Usecase,
) (*ImageHandler, error) {
	return &ImageHandler{
		logger: logger,
		authMd: authMd,
		userMd: userMd,
		uc:     uc,
	}, nil
}

func (h *ImageHandler) RegisterRoutes(r fiber.Router) {
	image := r.Group("/image", h.authMd.Require, h.userMd.CheckUser)
	image.Get("/:id", h.GeImage)
}

// GeImage godoc
// @Summary Get image by ID
// @Description Retrieves an image file by its ID, with validation, type restrictions, rate limiting, and size checks
// @Tags Image
// @Accept json
// @Produce octet-stream
// @Security ApiKeyAuth
// @Param id path string true "Image ID"
// @Success 200 {file} binary
// @Failure 400 {object} dto.Base
// @Failure 403 {object} dto.Base
// @Failure 404 {object} dto.Base
// @Failure 413 {object} dto.Base
// @Failure 429 {object} dto.BaseError
// @Failure 500 {object} dto.Base
// @Router /image/{id} [get]
func (h *ImageHandler) GeImage(c *fiber.Ctx) error {
	id := utils.PString(c.Params("id"))

	img, contentType, err := h.uc.GetImage(c.Context(), id)
	if err != nil {
		h.logger.Errorf("error fetching image %s: %v", *id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Base{
			Success: utils.PBool(false),
			Message: utils.PString("internal server error"),
		})
	}

	if contentType == nil || img == nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.Base{
			Success: utils.PBool(false),
			Message: utils.PString("image not found"),
		})
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
	if !allowedTypes[*contentType] {
		h.logger.Warnf("image type not allowed for ID %s: %s", *id, *contentType)
		return c.Status(fiber.StatusForbidden).JSON(dto.Base{
			Success: utils.PBool(false),
			Message: utils.PString("unsupported image type"),
		})
	}

	if exceeded := utils.SimpleRateLimit(c.IP(), *id); *exceeded {
		h.logger.Warnf("rate limit exceeded for IP %s and image %s", c.IP(), *id)
		return c.Status(fiber.StatusTooManyRequests).JSON(dto.Base{
			Success: utils.PBool(false),
			Error: &dto.BaseError{
				Code:    errors.ErrTooManyRequests,
				Message: utils.PString("too many requests"),
			},
		})
	}

	if len(img) > 1*1024*1024 {
		h.logger.Warnf("image too large for ID %s (%d bytes)", *id, len(img))
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(dto.Base{
			Success: utils.PBool(false),
			Message: utils.PString("image too large"),
		})
	}

	c.Set("Content-Type", *contentType)
	c.Set("Cache-Control", "public, max-age=86400")
	c.Set("X-Content-Type-Options", "nosniff")

	return c.Send(img)
}
