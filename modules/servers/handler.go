package servers

import (
	"log"
	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/database"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/middlewares"
	nhControllers "github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/controllers"
	nhRepositories "github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/repositories"
	nhUseCases "github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/usecases"
	userControllers "github.com/XzerozZ/Kasian_Phrom_BE/modules/user/controllers"
	userRepositories "github.com/XzerozZ/Kasian_Phrom_BE/modules/user/repositories"
	userUseCases "github.com/XzerozZ/Kasian_Phrom_BE/modules/user/usecases"
	newsControllers "github.com/XzerozZ/Kasian_Phrom_BE/modules/news/controllers"
	newsRepositories "github.com/XzerozZ/Kasian_Phrom_BE/modules/news/repositories"
	newsUseCases "github.com/XzerozZ/Kasian_Phrom_BE/modules/news/usecases"

	"gorm.io/gorm"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App, jwt configs.JWT ,supa configs.Supabase) {
	db := database.GetDB()
	if db == nil {
		log.Fatal("Failed to initialize database")
	}

	app.Use(cors.New(cors.Config{
        AllowOrigins: "http://localhost:3000",
        AllowMethods: "GET, POST, PUT, DELETE",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    }))

	setupNursingHouseRoutes(app, db, supa)
	SetupNewsRoutes(app, db)
	setupUserRoutes(app, db, jwt)
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status":  "Success",
			"message": "Welcome to the Nursing House System!",
		})
	})
}

func SetupNewsRoutes(app *fiber.App, db *gorm.DB) {
	newsRepository := newsRepositories.NewGormNewsRepository(db)
	newsUseCase := newsUseCases.NewNewsUseCase(newsRepository)
	newsController := newsControllers.NewNewsController(newsUseCase)

	app.Post("/news", newsController.CreateNewsHandler)
}

func setupUserRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	userRepository := userRepositories.NewGormUserRepository(db)
	userUseCase := userUseCases.NewUserUseCase(userRepository, jwt)
	userController := userControllers.NewUserController(userUseCase)

	app.Post("/register", userController.RegisterHandler)
	app.Post("/login", userController.LoginHandler)
	app.Post("/logout", middlewares.JWTMiddleware(jwt), userController.LogoutHandler)
}

func setupNursingHouseRoutes(app *fiber.App, db *gorm.DB, supa configs.Supabase) {
	nhRepository := nhRepositories.NewGormNhRepository(db)
	nhUseCase := nhUseCases.NewNhUseCase(nhRepository, supa)
	nhController := nhControllers.NewNhController(nhUseCase)

	nhGroup := app.Group("/nursinghouses")
	nhGroup.Post("/", nhController.CreateNhHandler)
	nhGroup.Get("/", nhController.GetAllNhHandler)
	nhGroup.Get("/active", nhController.GetAllActiveNhHandler)
	nhGroup.Get("/inactive", nhController.GetAllInactiveNhHandler)
	nhGroup.Get("/id" , nhController.GetNhNextIDHandler)
	nhGroup.Get("/:id", nhController.GetNhByIDHandler)
	nhGroup.Put("/:id", nhController.UpdateNhByIDHandler)
	nhGroup.Get("/:id/images", nhController.GetNhImages)
}