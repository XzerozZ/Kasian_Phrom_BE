package main

import (
	"log"
	"math"

	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/servers"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/database"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config := configs.LoadConfigs()
	database.InitDB(config.PostgreSQL)
	app := fiber.New(fiber.Config{
		BodyLimit: math.MaxInt64,
	})

	utils.StartScheduler()
	servers.SetupRoutes(app, config.JWT, config.Supabase, config.Mail, config.Recommend)
	serverAddress := config.App.Host + ":" + config.App.Port
	log.Printf("Server is running on %s", serverAddress)
	log.Fatal(app.Listen(serverAddress))
}
