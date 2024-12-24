package servers

import (
	"log"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/database"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/controllers"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/repositories"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App) {
	db := database.GetDB()
	if db == nil {
		log.Fatal("Failed to initialize database")
	}

	NhRepository := repositories.NewGormNhRepository(db)
	NhUsecase := usecases.NewNhUseCase(NhRepository)
	NhController := controllers.NewNhController(NhUsecase)
	app.Use(cors.New(cors.Config{
        AllowOrigins: "http://localhost:3000",
        AllowMethods: "GET, POST, PUT, DELETE",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    }))
	app.Post("/nursinghouses", NhController.CreateNhHandler)
	app.Get("/nursinghouses", NhController.GetAllNhHandler)
	app.Get("/nursinghouses/active", NhController.GetAllActiveNhHandler)
	app.Get("/nursinghouses/:id", NhController.GetNhByIDHandler)
	app.Put("/nursinghouses/:id", NhController.UpdateNhByIDHandler)
	app.Delete("/nursinghouses/:id", NhController.DeleteNhByIDHandler)

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status":  "Success",
			"message": "Welcome to the Nursing House System!",
		})
	})
}