package main

import (
	"log"
	"time"

	"currencies/internal/home"

	"currencies/config"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"

	"currencies/internal/importer"
	"currencies/pkg/database"

	v1 "currencies/internal/api/v1"

	"currencies/internal/middleware"
)

func main() {

	config.Init()

	dbConfig := config.NewDatabaseConfig()
	dbpool, err := database.CreateDbPool(dbConfig)
	if err != nil {
		log.Fatalf("failed to create DB pool: %v", err)
	}
	defer dbpool.Close()

	app := fiber.New()

	// Restore app
	app.Use(recover.New())

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("ok")
	})

	// Static files
	app.Get("/*", static.New("./web"))

	// Repositories
	currenciesRepo := home.NewCurrenciesRepository(dbpool)
	importRepo := importer.NewImportRepository(dbpool)
	v1Repo := v1.NewV1Repository(dbpool)

	// Настройки для импорта
	config.NewImportConfig()
	importConfig := config.NewImportConfig()

	// Port
	port := config.GetAppPort()

	// Routes
	home.NewHandler(app, currenciesRepo)               // Главная
	importer.NewHandler(app, importConfig, importRepo) // Импорт JSON

	// API v1
	trustedIPs := []string{"127.0.0.1"}
	demoLimiter := middleware.NewIPRateLimiter(1, time.Minute, trustedIPs)    // 1 запрос в минуту
	normalLimiter := middleware.NewIPRateLimiter(30, time.Minute, trustedIPs) // 30 запросов в минуту

	// 1. Общая группа /api/v1
	v1Group := app.Group("/api/v1", middleware.TypeServiceMiddleware())
	{
		// 2. Обычная подгруппа: normalLimiter
		normalGroup := v1Group.Group("", middleware.RateLimiter(normalLimiter))
		{
			v1.NewHandler(normalGroup, v1Repo)
		}

		// 3. Демо подгруппа: demoLimiter
		demoGroup := v1Group.Group("/demo", middleware.RateLimiter(demoLimiter))
		{
			v1.NewHandler(demoGroup, v1Repo)
		}
	}

	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}

}
