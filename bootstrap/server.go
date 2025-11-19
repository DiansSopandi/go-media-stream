package bootstrap

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/DiansSopandi/media_stream/db"
	"github.com/DiansSopandi/media_stream/docs"
	"github.com/DiansSopandi/media_stream/middlewares"
	"github.com/DiansSopandi/media_stream/pkg"
	"github.com/DiansSopandi/media_stream/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"golang.org/x/oauth2"
)

var (
	oauthConfig     *oauth2.Config
	appJWTSecret    []byte
	frontendURL     string
	appJWTExpiresIn time.Duration
)

func ServerInitialize() {
	database := db.InitDatabase()
	// defer database.Close()

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\n🔌 Closing database connection...")
		db.CloseDB()
		pkg.CloseRedis()
		os.Exit(0)
	}()

	expires := pkg.Cfg.Application.AppJWTAccessExpiresIn
	if expires <= 0 {
		expires = 3600
	}
	appJWTExpiresIn = time.Duration(expires) * time.Second

	// global error handler
	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.ErrorHandler,
	})

	// global middleware panic handler
	app.Use(middlewares.GlobalRecoveryMiddleware)

	allowed := strings.Split(pkg.Cfg.Application.CorsOrigins, ",")
	app.Use(cors.New(cors.Config{
		// AllowOrigins: "http://localhost:3000",
		AllowOrigins:     strings.Join(allowed, ","),
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))

	// global guard JWT authentication middleware
	app.Use(middlewares.JwtAuthGuard)

	// middlewares.InitRateLimiter()
	pkg.InitRedis()
	// apply global rate limit middleware all routes
	// duration := time.Minute
	// app.Use(middlewares.RateLimitMiddleware(&pkg.Cfg.Application.DefaultMaxRequestPerMinute, &duration))

	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server for Go Ride API"
	docs.SwaggerInfo.Version = "1.0.0"
	// docs.SwaggerInfo.Host = pkg.GetEnv("APP_URL")
	docs.SwaggerInfo.Host = pkg.Cfg.Application.AppUrl
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Route untuk Swagger UI
	app.Get("/swagger/*", swagger.HandlerDefault) // akses di /swagger/index.html

	// app.Get("/favicon.ico", func(c *fiber.Ctx) error {
	// 	return c.SendStatus(fiber.StatusNoContent)
	// })
	app.Static("/favicon.ico", "./public/favicon.ico")

	// Set up API root routes
	// app.Get("/", handlers.RootHandler)
	routes.SetupRoutes(app)

	// port := pkg.GetEnv("APP_PORT")
	port := pkg.Cfg.Application.AppPort

	log.Println("Server started successfully on port:", port)
	log.Printf("📊 Database connection status: %v", database != nil)

	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
