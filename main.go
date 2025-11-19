package main

import (
	"github.com/DiansSopandi/media_stream/bootstrap"
	"github.com/DiansSopandi/media_stream/cmd"
	_ "github.com/DiansSopandi/media_stream/docs"
)

// @title GoRide API
// @version 1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @BasePath /v1
func main() {
	// Load environment variables
	// pkg.LoadEnv()
	cmd.Execute()

	// Connect to the database
	// db := db.Connect()
	// defer db.Close()

	// Initialize database (create if not exists and connect)
	// database := db.InitDatabase()
	// defer database.Close()

	// app := fiber.New()

	// docs.SwaggerInfo.Title = "Swagger Example API"
	// docs.SwaggerInfo.Description = "This is a sample server for Go Ride API"
	// docs.SwaggerInfo.Version = "1.0.0"
	// // docs.SwaggerInfo.Host = pkg.GetEnv("APP_URL")
	// docs.SwaggerInfo.Host = pkg.Cfg.Application.AppUrl
	// docs.SwaggerInfo.BasePath = "/"
	// docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// // Route untuk Swagger UI
	// app.Get("/swagger/*", swagger.HandlerDefault) // akses di /swagger/index.html

	// // app.Get("/favicon.ico", func(c *fiber.Ctx) error {
	// // 	return c.SendStatus(fiber.StatusNoContent)
	// // })
	// app.Static("/favicon.ico", "./public/favicon.ico")

	// // Set up API root routes
	// // app.Get("/", handlers.RootHandler)
	// routes.SetupRoutes(app)

	// port := pkg.GetEnv("APP_PORT")

	// log.Println("Server started successfully on port:", port)

	// if err := app.Listen(":" + port); err != nil {
	// 	log.Fatalf("Error starting server: %v", err)
	// }

	bootstrap.ServerInitialize()
}
