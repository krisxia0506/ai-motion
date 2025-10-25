package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xiajiayi/ai-motion/internal/application/service"
	"github.com/xiajiayi/ai-motion/internal/domain/character"
	"github.com/xiajiayi/ai-motion/internal/domain/novel"
	"github.com/xiajiayi/ai-motion/internal/domain/scene"
	"github.com/xiajiayi/ai-motion/internal/infrastructure/config"
	"github.com/xiajiayi/ai-motion/internal/infrastructure/database"
	"github.com/xiajiayi/ai-motion/internal/infrastructure/repository/mysql"
	"github.com/xiajiayi/ai-motion/internal/interfaces/http/handler"
	"github.com/xiajiayi/ai-motion/internal/interfaces/http/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	var novelHandler *handler.NovelHandler
	var characterHandler *handler.CharacterHandler
	var sceneHandler *handler.SceneHandler

	if cfg.Database.Host != "" && cfg.Database.Password != "" {
		dbCfg := &database.Config{
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			Database: cfg.Database.Database,
		}

		dbConn, err := database.NewMySQLConnection(dbCfg)
		if err != nil {
			log.Printf("Warning: Failed to connect to database: %v", err)
			log.Println("Starting server without database connection...")
		} else {
			defer database.CloseMySQLConnection(dbConn)
			log.Println("Database connection established")

			migrationsPath := os.Getenv("MIGRATIONS_PATH")
			if migrationsPath == "" {
				migrationsPath = "./internal/infrastructure/database/migrations"
			}

			if err := database.RunMigrations(dbConn, migrationsPath); err != nil {
				log.Printf("Warning: Failed to run migrations: %v", err)
			} else {
				log.Println("Database migrations completed")
			}

			novelRepo := mysql.NewNovelRepository(dbConn)
			chapterRepo := mysql.NewChapterRepository(dbConn)
			characterRepo := mysql.NewMySQLCharacterRepository(dbConn)
			sceneRepo := mysql.NewMySQLSceneRepository(dbConn)

			parserService := novel.NewParserService()
			novelService := service.NewNovelService(novelRepo, chapterRepo, parserService)
			novelHandler = handler.NewNovelHandler(novelService)

			extractorService := character.NewCharacterExtractorService(characterRepo)
			characterService := service.NewCharacterService(characterRepo, novelRepo, extractorService)
			characterHandler = handler.NewCharacterHandler(characterService)

			dividerService := scene.NewSceneDividerService(sceneRepo)
			promptGeneratorService := scene.NewPromptGeneratorService(sceneRepo)
			sceneService := service.NewSceneService(sceneRepo, chapterRepo, characterRepo, dividerService, promptGeneratorService)
			sceneHandler = handler.NewSceneHandler(sceneService)
		}
	} else {
		log.Println("Database configuration not found, starting without database...")
	}

	r := gin.New()

	rateLimiter := middleware.NewRateLimiter(100, 200)
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.ErrorHandler())
	r.Use(rateLimiter.Middleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "ai-motion",
		})
	})

	v1 := r.Group("/api/v1")
	{
		if novelHandler != nil {
			novelGroup := v1.Group("/novel")
			{
				novelGroup.POST("/upload", novelHandler.Upload)
				novelGroup.GET("/:id", novelHandler.Get)
				novelGroup.GET("", novelHandler.List)
				novelGroup.DELETE("/:id", novelHandler.Delete)
				novelGroup.GET("/:id/chapters", novelHandler.GetChapters)
			}
		} else {
			v1.POST("/novel/upload", func(c *gin.Context) {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"error": "Database not configured",
				})
			})
		}

		if characterHandler != nil {
			characterGroup := v1.Group("/characters")
			{
				characterGroup.POST("/novel/:novel_id/extract", characterHandler.Extract)
				characterGroup.GET("/:id", characterHandler.Get)
				characterGroup.GET("/novel/:novel_id", characterHandler.ListByNovel)
				characterGroup.PUT("/:id", characterHandler.Update)
				characterGroup.DELETE("/:id", characterHandler.Delete)
				characterGroup.POST("/merge", characterHandler.Merge)
			}
		}

		if sceneHandler != nil {
			sceneGroup := v1.Group("/scenes")
			{
				sceneGroup.POST("/chapter/:chapter_id/divide", sceneHandler.DivideChapter)
				sceneGroup.GET("/:id", sceneHandler.Get)
				sceneGroup.GET("/chapter/:chapter_id", sceneHandler.ListByChapter)
				sceneGroup.GET("/novel/:novel_id", sceneHandler.ListByNovel)
				sceneGroup.DELETE("/:id", sceneHandler.Delete)
			}

			promptGroup := v1.Group("/prompts")
			{
				promptGroup.POST("/generate", sceneHandler.GeneratePrompt)
				promptGroup.POST("/generate/batch", sceneHandler.GenerateBatchPrompts)
			}
		}

		generate := v1.Group("/generate")
		{
			generate.POST("/scene", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "generate scene endpoint - coming soon"})
			})
			generate.POST("/voice", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "generate voice endpoint - coming soon"})
			})
		}

		v1.POST("/anime/:id/export", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "export anime endpoint - coming soon"})
		})
	}

	serverAddr := ":" + cfg.Server.Port
	log.Printf("Starting AI-Motion server on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
