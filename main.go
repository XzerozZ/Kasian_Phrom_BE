package main

import (
	"log"
	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/servers"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/database"
	
	"github.com/gofiber/fiber/v2"
)

func main() {
	config := configs.LoadConfigs()
	database.InitDB(config.PostgreSQL)
	app := fiber.New()
	servers.SetupRoutes(app, config.JWT, config.Supabase)
	serverAddress := config.App.Host + ":" + config.App.Port
	log.Printf("Server is running on %s", serverAddress)
	log.Fatal(app.Listen(serverAddress))
}