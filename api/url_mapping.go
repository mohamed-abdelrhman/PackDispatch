package api

import (
	"github.com/gofiber/fiber/v2"
	pack_size "github.com/mohamed-abdelrhman/pack-dispatch/api/handlers/pack-size"
)

// MapUrl abstracted function to map and register all the url for the application
// helps to keep all the endpoints' definition in one place
// the only place to interact with handlers and middlewares
func MapUrl(app *fiber.App, handlers ...fiber.Handler) {
	// routing groups
	api := app.Group("api")
	v1 := api.Group("v1")
	for i := 0; i < len(handlers); i++ {
		v1.Use(handlers[i])
	}
	orders := v1.Group("orders")
	orders.Post("/calculate-packs", pack_size.Packs)
}
