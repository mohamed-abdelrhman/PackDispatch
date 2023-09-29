package config

import (
	"github.com/gofiber/fiber/v2"
	restErrors "github.com/mohamed-abdelrhman/pack-dispatch/pkg/errors"
	"strconv"
	"time"
)

func FiberConfig() fiber.Config {
	readTimeoutSecondsCount, _ := strconv.Atoi(Environment.ServerReadTimeout)
	return fiber.Config{
		ReadTimeout:  time.Second * time.Duration(readTimeoutSecondsCount),
		ErrorHandler: defaultErrorHandler,
	}
}

// defaultErrorHandler used to catch all unhandled run time errors mainly panics
var defaultErrorHandler = func(c *fiber.Ctx, err error) error {
	internalErr := restErrors.NewInternalServerError(err.Error())
	return c.Status(internalErr.StatusCode()).JSON(internalErr)
}
