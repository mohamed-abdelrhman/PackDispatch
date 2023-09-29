package pack_size

import (
	"github.com/gofiber/fiber/v2"
	pack_size "github.com/mohamed-abdelrhman/pack-dispatch/internal/pack-size"
	restErrors "github.com/mohamed-abdelrhman/pack-dispatch/pkg/errors"
	"github.com/mohamed-abdelrhman/pack-dispatch/pkg/shared"
	"net/http"
)

var (
	packSizeService = pack_size.NewService()
)

func Packs(c *fiber.Ctx) error {
	dto := new(pack_size.OrderRequestDto)
	if intErr := c.BodyParser(dto); intErr != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.StatusCode()).JSON(badReq)
	}

	packsNeeded, err := packSizeService.CalculatePacks(dto.Quantity)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(shared.NewResponse(packsNeeded))
}
