package main

import (
	"Kopsis-Spensa/internal/config"
	"Kopsis-Spensa/internal/handlers"
	"Kopsis-Spensa/internal/middleware"
	"Kopsis-Spensa/internal/models"
	"Kopsis-Spensa/internal/seeder"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	config.Init()
	config.InitSession()

	err := config.DB.AutoMigrate(
		&models.BackupRestore{},
		&models.Category{},
		&models.MutasiStock{},
		&models.Product{},
		&models.Transactions{},
		&models.TransactionsItems{},
		&models.User{},
	)
	if err != nil {
		log.Fatal(err)
	}
	// Seeder
	seeder.SeedUsers()
	seeder.SeedCategory()

	engine := html.New("./web/views", ".html")
	engine.AddFunc("subtract", func(a, b float64) float64 {
		return a - b
	})
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	// Static files
	app.Static("/static", "./web/static")
	app.Static("/uploads", "./uploads")
	app.Static("/storage", "./storage")
	app.Static("/storage/backups", "./uploads/backups")

	//Kasir
	PosRoutes := app.Group("/pos", middleware.IsAuthenticated)
	PosRoutes.Get("/", handlers.GetPos)
	app.Get("/api/pos/products", handlers.GetProductsPOS)
	app.Post("/api/pos/checkout", handlers.CheckOutProcess)
	// Routes Authenticated
	app.Get("/login", middleware.RedirectIfLoggedIn, handlers.ShowLogin)
	app.Post("/loginProccess", handlers.ProccessLogin)
	app.Get("/logout", handlers.Logout)

	//Routes Dashboard
	DashboardRoutes := app.Group("/", middleware.IsAuthenticated)
	DashboardRoutes.Get("/", handlers.GetDashboard)
	DashboardRoutes.Get("/dashboard", handlers.GetDashboard)
	//Users
	DashboardRoutes.Get("/manage_users", handlers.GetUsers)
	DashboardRoutes.Post("/users", handlers.CreateUser)
	DashboardRoutes.Get("/users/:uid", handlers.GetUserByUID)
	DashboardRoutes.Put("/users/:uid", handlers.UpdateUser)
	DashboardRoutes.Delete("/users/:uid", handlers.DeleteUser)
	//Category
	DashboardRoutes.Get("/manage_category", handlers.GetCategories)
	//Product
	DashboardRoutes.Get("/manage_product", handlers.GetManageProducts)
	DashboardRoutes.Post("/products", handlers.CreateProduct)
	DashboardRoutes.Get("/products/:uid", handlers.GetProductByUID)
	DashboardRoutes.Put("/products/:uid", handlers.UpdateProduct)
	DashboardRoutes.Delete("/products/:uid", handlers.DeleteProduct)
	//Stock
	DashboardRoutes.Get("/manage_stock", handlers.GetStock)
	DashboardRoutes.Post("/stock/add", handlers.AddStock)
	DashboardRoutes.Post("/stock/reduce", handlers.ReduceStock)
	//Reports
	DashboardRoutes.Get("/reports", handlers.GetReports)
	DashboardRoutes.Get("/reports/export", handlers.ExportReport)
	//BackupRestore
	DashboardRoutes.Get("/backup_restore", handlers.GetBackupRestore)
	DashboardRoutes.Post("/backup/create", handlers.CreateBackup)
	DashboardRoutes.Post("/backup/restore", handlers.RestoreBackup)
	DashboardRoutes.Delete("/backup/delete", handlers.DeleteBackup)
	DashboardRoutes.Get("/backup/download/:filename", handlers.DownloadBackup)
	DashboardRoutes.Post("/backup/upload-restore", handlers.UploadAndRestore)
	//Profile
	DashboardRoutes.Get("/profile", handlers.GetProfile)
	DashboardRoutes.Post("/profile/update", handlers.UpdateProfile)
	DashboardRoutes.Post("/profile/password", handlers.UpdatePassword)
	DashboardRoutes.Post("/profile/photo", handlers.UploadPhoto)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server berjalan di http://localhost:%s\n", port)
	log.Fatal(app.Listen(":" + port))
}
