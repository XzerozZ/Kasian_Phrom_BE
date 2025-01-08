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
	CrpControllers "github.com/XzerozZ/Kasian_Phrom_BE/modules/Calculate_retirement_plan/controllers"
	CrpRepositories "github.com/XzerozZ/Kasian_Phrom_BE/modules/Calculate_retirement_plan/repositories"
	CrpUseCases "github.com/XzerozZ/Kasian_Phrom_BE/modules/Calculate_retirement_plan/usecases"

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
	SetupNewsRoutes(app, db, supa)
	setupUserRoutes(app, db, jwt, supa)
	setupCrpRoutes(app, db)

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status":  "Success",
			"message": "Welcome to the Nursing House System!",
		})
	})
}

func SetupNewsRoutes(app *fiber.App, db *gorm.DB, supa configs.Supabase) {
	newsRepository := newsRepositories.NewGormNewsRepository(db)
	newsUseCase := newsUseCases.NewNewsUseCase(newsRepository, supa)
	newsController := newsControllers.NewNewsController(newsUseCase)

	newsGroup := app.Group("/news")
	newsGroup.Post("/", newsController.CreateNewsHandler)
	newsGroup.Get("/", newsController.GetAllNewsHandler)
	newsGroup.Get("/id" , newsController.GetNewsNextIDHandler)
	newsGroup.Get("/:id", newsController.GetNewsByIDHandler)
	newsGroup.Put("/:id", newsController.UpdateNewsByIDHandler)
	newsGroup.Delete("/:id", newsController.DeleteNewsByIDHandler)
}

func setupUserRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT, supa configs.Supabase) {
	userRepository := userRepositories.NewGormUserRepository(db)
	userUseCase := userUseCases.NewUserUseCase(userRepository, jwt, supa)
	userController := userControllers.NewUserController(userUseCase)

	authGroup := app.Group("/auth")
	authGroup.Post("/register", userController.RegisterHandler)
	authGroup.Post("/login", userController.LoginHandler)
	authGroup.Post("/admin/login", userController.LoginAdminHandler)
	authGroup.Post("/logout", middlewares.JWTMiddleware(jwt), userController.LogoutHandler)

	userGroup := app.Group("/user")
	userGroup.Put("/:id", userController.UpdateUserByIDHandler)
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
}

func setupCrpRoutes(app *fiber.App, db *gorm.DB) {
	finRepository := CrpRepositories.NewGormFinRepository(db)
	finUseCase := CrpUseCases.NewFinUseCase(finRepository)
	finController := CrpControllers.NewFinController(finUseCase)

	finGroup := app.Group("/financial")
	finGroup.Post("/", finController.CreateFinHandler)
    finGroup.Get("/:id", finController.GetFinByIDHandler)
	///////

	assRepository := CrpRepositories.NewGormAssRepository(db)
	assUseCase := CrpUseCases.NewAssUseCase(assRepository)
	assController := CrpControllers.NewAssController(assUseCase)

	assGroup := app.Group("/asset")
	assGroup.Post("/", assController.CreateAssHandler)
    assGroup.Get("/:id", assController.GetAssByIDHandler)
    assGroup.Get("/user/:username", assController.GetAssByUsernameHandler)
    assGroup.Put("/:id", assController.UpdateAssByIDHandler)
    assGroup.Delete("/:id", assController.DeleteAssByIDHandler)


	///////

	//เพิ่ม user ด้วย
	userRepository := userRepositories.NewGormUserRepository(db)

	retRepository := CrpRepositories.NewGormRetRepository(db)
	retUseCase := CrpUseCases.NewRetUseCase(retRepository, finRepository, assRepository, userRepository)
	retController := CrpControllers.NewRetController(retUseCase)

	retireGroup := app.Group("retire")
	retireGroup.Post("/", retController.CreateRetHandler)
	
}