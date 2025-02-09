package servers

import (
	"log"

	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	assetControllers "github.com/XzerozZ/Kasian_Phrom_BE/modules/asset/controllers"
	assetRepositories "github.com/XzerozZ/Kasian_Phrom_BE/modules/asset/repositories"
	assetUseCases "github.com/XzerozZ/Kasian_Phrom_BE/modules/asset/usecases"
	favControllers "github.com/XzerozZ/Kasian_Phrom_BE/modules/favorite/controllers"
	favRepositories "github.com/XzerozZ/Kasian_Phrom_BE/modules/favorite/repositories"
	favUseCases "github.com/XzerozZ/Kasian_Phrom_BE/modules/favorite/usecases"
	loanControllers "github.com/XzerozZ/Kasian_Phrom_BE/modules/loan/controllers"
	loanRepositories "github.com/XzerozZ/Kasian_Phrom_BE/modules/loan/repositories"
	loanUseCases "github.com/XzerozZ/Kasian_Phrom_BE/modules/loan/usecases"
	newsControllers "github.com/XzerozZ/Kasian_Phrom_BE/modules/news/controllers"
	newsRepositories "github.com/XzerozZ/Kasian_Phrom_BE/modules/news/repositories"
	newsUseCases "github.com/XzerozZ/Kasian_Phrom_BE/modules/news/usecases"
	nhControllers "github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/controllers"
	nhRepositories "github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/repositories"
	nhUseCases "github.com/XzerozZ/Kasian_Phrom_BE/modules/nursing_house/usecases"
	quizControllers "github.com/XzerozZ/Kasian_Phrom_BE/modules/quiz/controllers"
	quizRepositories "github.com/XzerozZ/Kasian_Phrom_BE/modules/quiz/repositories"
	quizUseCases "github.com/XzerozZ/Kasian_Phrom_BE/modules/quiz/usecases"
	retirementControllers "github.com/XzerozZ/Kasian_Phrom_BE/modules/retirement_plan/controllers"
	retirementRepositories "github.com/XzerozZ/Kasian_Phrom_BE/modules/retirement_plan/repositories"
	retirementUseCases "github.com/XzerozZ/Kasian_Phrom_BE/modules/retirement_plan/usecases"
	transControllers "github.com/XzerozZ/Kasian_Phrom_BE/modules/transaction/controllers"
	transRepositories "github.com/XzerozZ/Kasian_Phrom_BE/modules/transaction/repositories"
	transUseCases "github.com/XzerozZ/Kasian_Phrom_BE/modules/transaction/usecases"
	userControllers "github.com/XzerozZ/Kasian_Phrom_BE/modules/user/controllers"
	userRepositories "github.com/XzerozZ/Kasian_Phrom_BE/modules/user/repositories"
	userUseCases "github.com/XzerozZ/Kasian_Phrom_BE/modules/user/usecases"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/database"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, jwt configs.JWT, supa configs.Supabase, mail configs.Mail) {
	db := database.GetDB()
	if db == nil {
		log.Fatal("Failed to initialize database")
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET, POST, PUT, DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	userRepository := userRepositories.NewGormUserRepository(db)
	retirementRepository := retirementRepositories.NewGormRetirementRepository(db)
	nhRepository := nhRepositories.NewGormNhRepository(db)
	assetRepository := assetRepositories.NewGormAssetRepository(db)
	userUseCase := userUseCases.NewUserUseCase(userRepository, retirementRepository, assetRepository, nhRepository, jwt, supa, mail)

	setupNursingHouseRoutes(app, db, supa)
	SetupNewsRoutes(app, db, supa)
	setupFavoriteRoutes(app, jwt, db)
	setupAssetRoutes(app, jwt, db)
	setupUserRoutes(app, db, jwt, supa, mail)
	setupRetirementRoutes(app, jwt, userUseCase, db)
	setupLoanRoutes(app, jwt, db)
	setupQuizRoutes(app, jwt, db)
	setupTransactionRoutes(app, jwt, db)

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
	newsGroup.Get("/id", newsController.GetNewsNextIDHandler)
	newsGroup.Get("/:id", newsController.GetNewsByIDHandler)
	newsGroup.Put("/:id", newsController.UpdateNewsByIDHandler)
	newsGroup.Delete("/:id", newsController.DeleteNewsByIDHandler)
}

func setupUserRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT, supa configs.Supabase, mail configs.Mail) {
	userRepository := userRepositories.NewGormUserRepository(db)
	retirementRepository := retirementRepositories.NewGormRetirementRepository(db)
	nhRepository := nhRepositories.NewGormNhRepository(db)
	assetRepository := assetRepositories.NewGormAssetRepository(db)
	userUseCase := userUseCases.NewUserUseCase(userRepository, retirementRepository, assetRepository, nhRepository, jwt, supa, mail)
	userController := userControllers.NewUserController(userUseCase)

	authGroup := app.Group("/auth")
	authGroup.Post("/register", userController.RegisterHandler)
	authGroup.Post("/login", userController.LoginHandler)
	authGroup.Post("/google/login", userController.LoginWithGoogleHandler)
	authGroup.Post("/admin/login", userController.LoginAdminHandler)
	authGroup.Post("/forgotpassword", userController.ForgotPasswordHandler)
	authGroup.Post("/forgotpassword/otp", userController.VerifyOTPHandler)
	authGroup.Put("/forgotpassword/changepassword", userController.ChangedPasswordHandler)
	authGroup.Put("/resetpassword", middlewares.JWTMiddleware(jwt), userController.ResetPasswordHandler)
	authGroup.Post("/logout", middlewares.JWTMiddleware(jwt), userController.LogoutHandler)

	userGroup := app.Group("/user")
	userGroup.Get("/", middlewares.JWTMiddleware(jwt), userController.GetUserByIDHandler)
	userGroup.Get("/plan", middlewares.JWTMiddleware(jwt), userController.GetRetirementPlanHandler)
	userGroup.Get("/selected", middlewares.JWTMiddleware(jwt), userController.GetSelectedHouseHandler)
	userGroup.Put("/user", middlewares.JWTMiddleware(jwt), userController.UpdateUserByIDHandler)
	userGroup.Put("/:nh_id", middlewares.JWTMiddleware(jwt), userController.UpdateSelectedHouseHandler)

	historyGroup := app.Group("/history")
	historyGroup.Post("/", middlewares.JWTMiddleware(jwt), userController.CreateHistoryHandler)
	historyGroup.Get("/", middlewares.JWTMiddleware(jwt), userController.GetHistoryByUserIDHandler)
	historyGroup.Get("/summary", middlewares.JWTMiddleware(jwt), userController.GetSummaryHistoryByUserIDHandler)
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
	nhGroup.Get("/id", nhController.GetNhNextIDHandler)
	nhGroup.Get("/:id", nhController.GetNhByIDHandler)
	nhGroup.Put("/:id", nhController.UpdateNhByIDHandler)
}

func setupFavoriteRoutes(app *fiber.App, jwt configs.JWT, db *gorm.DB) {
	favRepository := favRepositories.NewGormFavRepository(db)
	favUseCase := favUseCases.NewFavUseCase(favRepository)
	favController := favControllers.NewFavController(favUseCase)

	favGroup := app.Group("/favorite")
	favGroup.Post("/", middlewares.JWTMiddleware(jwt), favController.CreateFavHandler)
	favGroup.Get("/user", middlewares.JWTMiddleware(jwt), favController.GetFavByUserIDHandler)
	favGroup.Get("/:nh_id", middlewares.JWTMiddleware(jwt), favController.CheckFavHandler)
	favGroup.Delete("/:nh_id", middlewares.JWTMiddleware(jwt), favController.DeleteFavByIDHandler)
}

func setupAssetRoutes(app *fiber.App, jwt configs.JWT, db *gorm.DB) {
	assetRepository := assetRepositories.NewGormAssetRepository(db)
	assetUseCase := assetUseCases.NewAssetUseCase(assetRepository)
	assetController := assetControllers.NewAssetController(assetUseCase)

	assetGroup := app.Group("/asset")
	assetGroup.Post("/", middlewares.JWTMiddleware(jwt), assetController.CreateAssetHandler)
	assetGroup.Get("/:id", assetController.GetAssetByIDHandler)
	assetGroup.Get("/", middlewares.JWTMiddleware(jwt), assetController.GetAssetByUserIDHandler)
	assetGroup.Put("/:id", middlewares.JWTMiddleware(jwt), assetController.UpdateAssetByIDHandler)
	assetGroup.Delete("/:id", assetController.DeleteAssetByIDHandler)
}

func setupRetirementRoutes(app *fiber.App, jwt configs.JWT, userUseCase userUseCases.UserUseCase, db *gorm.DB) {
	retirementRepository := retirementRepositories.NewGormRetirementRepository(db)
	retirementUseCase := retirementUseCases.NewRetirementUseCase(retirementRepository, userUseCase)
	retirementController := retirementControllers.NewRetirementController(retirementUseCase)

	retirementGroup := app.Group("/retirement")
	retirementGroup.Post("/", middlewares.JWTMiddleware(jwt), retirementController.CreateRetirementHandler)
	retirementGroup.Get("/", middlewares.JWTMiddleware(jwt), retirementController.GetRetirementByUserIDHandler)
	retirementGroup.Put("/", middlewares.JWTMiddleware(jwt), retirementController.UpdateRetirementHandler)
}

func setupLoanRoutes(app *fiber.App, jwt configs.JWT, db *gorm.DB) {
	loanRepository := loanRepositories.NewGormLoanRepository(db)
	transRepository := transRepositories.NewGormTransRepository(db)
	loanUseCase := loanUseCases.NewLoanUseCase(loanRepository, transRepository)
	transUseCase := transUseCases.NewTransactionUseCase(transRepository, loanRepository)
	loanController := loanControllers.NewLoanController(loanUseCase, transUseCase)

	loanGroup := app.Group("/loan")
	loanGroup.Post("/", middlewares.JWTMiddleware(jwt), loanController.CreateLoanHandler)
	loanGroup.Get("/:id", loanController.GetLoanByIDHandler)
	loanGroup.Get("/", middlewares.JWTMiddleware(jwt), loanController.GetLoanByUserIDHandler)
	loanGroup.Put("/:id/status", middlewares.JWTMiddleware(jwt), loanController.UpdateLoanStatusByIDHandler)
	loanGroup.Delete("/:id", loanController.DeleteLoanHandler)
}

func setupQuizRoutes(app *fiber.App, jwt configs.JWT, db *gorm.DB) {
	quizRepository := quizRepositories.NewGormQuizRepository(db)
	quizUseCase := quizUseCases.NewQuizUseCase(quizRepository)
	quizController := quizControllers.NewQuizController(quizUseCase)

	quizGroup := app.Group("/quiz")
	quizGroup.Post("/", middlewares.JWTMiddleware(jwt), quizController.CreateQuizHandler)
	quizGroup.Get("/", middlewares.JWTMiddleware(jwt), quizController.GetQuizByUserIDHandler)
}

func setupTransactionRoutes(app *fiber.App, jwt configs.JWT, db *gorm.DB) {
	transRepository := transRepositories.NewGormTransRepository(db)
	loanRepository := loanRepositories.NewGormLoanRepository(db)
	transUseCase := transUseCases.NewTransactionUseCase(transRepository, loanRepository)
	transController := transControllers.NewTransactionController(transUseCase)

	transGroup := app.Group("/transaction")
	transGroup.Post("/all", transController.CreateTransactionsForAllUsersHandler)
	transGroup.Get("/", middlewares.JWTMiddleware(jwt), transController.GetTransactionByUserIDHandler)
	transGroup.Put("/:id", middlewares.JWTMiddleware(jwt), transController.MarkTransactiontoPaidHandler)
}
