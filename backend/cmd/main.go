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
	"github.com/xiajiayi/ai-motion/internal/infrastructure/ai/gemini"
	"github.com/xiajiayi/ai-motion/internal/infrastructure/ai/sora"
	"github.com/xiajiayi/ai-motion/internal/infrastructure/config"
	"github.com/xiajiayi/ai-motion/internal/infrastructure/database"
	infra_middleware "github.com/xiajiayi/ai-motion/internal/infrastructure/middleware"
	"github.com/xiajiayi/ai-motion/internal/infrastructure/repository/supabase"
	"github.com/xiajiayi/ai-motion/internal/infrastructure/storage/local"
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
	var generationHandler *handler.GenerationHandler
	var mangaWorkflowHandler *handler.MangaWorkflowHandler

	geminiBaseURL := os.Getenv("GEMINI_BASE_URL")
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	soraBaseURL := os.Getenv("SORA_BASE_URL")
	soraAPIKey := os.Getenv("SORA_API_KEY")
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}

	var geminiClient *gemini.Client
	var soraClient *sora.Client

	if geminiBaseURL != "" && geminiAPIKey != "" {
		client, clientErr := gemini.NewClient(geminiBaseURL, geminiAPIKey)
		if clientErr != nil {
			log.Printf("Warning: Failed to initialize Gemini client: %v", clientErr)
		} else {
			geminiClient = client
			log.Printf("Gemini client initialized (baseURL: %s)", geminiBaseURL)
		}
	} else {
		log.Println("GEMINI_BASE_URL or GEMINI_API_KEY not set, AI image generation will be unavailable")
	}

	if soraBaseURL != "" && soraAPIKey != "" {
		client, clientErr := sora.NewClient(soraBaseURL, soraAPIKey)
		if clientErr != nil {
			log.Printf("Warning: Failed to initialize Sora client: %v", clientErr)
		} else {
			soraClient = client
			log.Printf("Sora client initialized (baseURL: %s)", soraBaseURL)
		}
	} else {
		log.Println("SORA_BASE_URL or SORA_API_KEY not set, AI video generation will be unavailable")
	}

	fileStorage, storageErr := local.NewFileStorage(storagePath)
	if storageErr != nil {
		log.Printf("Warning: Failed to initialize file storage: %v", storageErr)
	} else {
		log.Printf("File storage initialized at %s", storagePath)
		_ = fileStorage
	}

	log.Println("=== Service Initialization ===")
	log.Printf("Supabase URL configured: %v", cfg.Supabase.URL != "")
	log.Printf("Supabase API Key configured: %v", cfg.Supabase.APIKey != "")
	log.Printf("Gemini client available: %v", geminiClient != nil)
	log.Printf("Sora client available: %v", soraClient != nil)

	if cfg.Supabase.URL != "" && cfg.Supabase.APIKey != "" {
		supabaseCfg := &database.SupabaseConfig{
			URL:    cfg.Supabase.URL,
			APIKey: cfg.Supabase.APIKey,
		}

		supabaseClient, err := database.NewSupabaseClient(supabaseCfg)
		if err != nil {
			log.Printf("ERROR: Failed to connect to Supabase: %v", err)
			log.Println("Starting server without database connection...")
		} else {
			log.Println("✓ Supabase connection established")

			novelRepo := supabase.NewNovelRepository(supabaseClient)
			chapterRepo := supabase.NewChapterRepository(supabaseClient)
			characterRepo := supabase.NewCharacterRepository(supabaseClient)
			sceneRepo := supabase.NewSceneRepository(supabaseClient)
			mediaRepo := supabase.NewMediaRepository(supabaseClient)
			taskRepo := supabase.NewTaskRepository(supabaseClient)

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

			if geminiClient != nil && soraClient != nil {
				generationService := service.NewGenerationService(mediaRepo, sceneRepo, geminiClient, soraClient)
				generationHandler = handler.NewGenerationHandler(generationService)
				log.Println("Generation service initialized")
			} else {
				log.Println("AI clients not available, generation service disabled")
			}

			if geminiClient != nil {
				mangaWorkflowService := service.NewMangaWorkflowService(
					taskRepo,
					novelRepo,
					chapterRepo,
					characterRepo,
					sceneRepo,
					mediaRepo,
					parserService,
					extractorService,
					dividerService,
					geminiClient,
				)
				mangaWorkflowHandler = handler.NewMangaWorkflowHandler(mangaWorkflowService)
				log.Println("✓ Manga workflow service initialized")
			} else {
				log.Println("✗ Manga workflow DISABLED: Gemini client not available")
			}
		}
	} else {
		log.Println("✗ Supabase configuration not found, starting without database...")
	}

	log.Println("=== Handler Status ===")
	log.Printf("Novel Handler: %v", novelHandler != nil)
	log.Printf("Character Handler: %v", characterHandler != nil)
	log.Printf("Scene Handler: %v", sceneHandler != nil)
	log.Printf("Generation Handler: %v", generationHandler != nil)
	log.Printf("Manga Workflow Handler: %v", mangaWorkflowHandler != nil)
	log.Println("=============================")

	r := gin.New()

	rateLimiter := middleware.NewRateLimiter(100, 200)
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.ErrorHandler())
	r.Use(rateLimiter.Middleware())

	// 初始化认证中间件
	var authMiddleware *infra_middleware.AuthMiddleware
	if cfg.Supabase.JWTSecret != "" {
		authMiddleware = infra_middleware.NewAuthMiddleware(cfg.Supabase.JWTSecret)
		log.Println("JWT authentication middleware initialized")
	} else {
		log.Println("Warning: SUPABASE_JWT_SECRET not configured, authentication disabled")
	}

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

		if generationHandler != nil {
			generateGroup := v1.Group("/generate")
			{
				generateGroup.POST("/image", generationHandler.GenerateImage)
				generateGroup.POST("/video", generationHandler.GenerateVideo)
				generateGroup.POST("/batch", generationHandler.BatchGenerate)
				generateGroup.GET("/status/:scene_id", generationHandler.GetStatus)
			}
		} else {
			v1.POST("/generate/image", func(c *gin.Context) {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"error": "AI services not configured",
				})
			})
		}

		if mangaWorkflowHandler != nil {
			mangaGroup := v1.Group("/manga")
			if authMiddleware != nil {
				mangaGroup.Use(authMiddleware.SupabaseAuth()) // 应用认证中间件
			}
			{
				mangaGroup.POST("/generate", mangaWorkflowHandler.GenerateManga)
				mangaGroup.GET("/task/:task_id", mangaWorkflowHandler.GetTaskStatus)
				mangaGroup.GET("/tasks", mangaWorkflowHandler.GetTaskList)
				mangaGroup.POST("/task/:task_id/cancel", mangaWorkflowHandler.CancelTask)
			}
		} else {
			v1.POST("/manga/generate", func(c *gin.Context) {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"error": "Manga workflow service not configured",
				})
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
