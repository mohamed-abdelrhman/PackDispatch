package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mohamed-abdelrhman/pack-dispatch/api"
	"github.com/mohamed-abdelrhman/pack-dispatch/pkg/config"
	"github.com/mohamed-abdelrhman/pack-dispatch/pkg/migration"
	"github.com/mohamed-abdelrhman/pack-dispatch/pkg/seeder"
	"github.com/mohamed-abdelrhman/pack-dispatch/pkg/server"
	"github.com/mohamed-abdelrhman/pack-dispatch/pkg/sqlclient"
)

func main() {
	app := fiber.New(config.FiberConfig())

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())
	api.MapUrl(app)
	dbClient := sqlclient.OpenDBConnection()
	migrationService := migration.NewService(dbClient)
	migrationService.Run()

	seederService := seeder.NewService(dbClient)
	seederService.Run()

	server.StartServerWithGracefulShutdown(app)
}
