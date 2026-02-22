package main

import (
	"context"
	"fmt"
	"hackaton-service-api/internal/entity"
	"hackaton-service-api/internal/handler"
	"hackaton-service-api/internal/infra/database"
	"hackaton-service-api/internal/infra/service"
	"hackaton-service-api/internal/middleware"
	"hackaton-service-api/internal/usecase"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "hackaton-service-api/docs"
)

// @title FIAP X - API de Processamento de Vídeos
// @version 1.0
// @description API para upload e gerenciamento de processamento de vídeos.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	awsRegion := getEnv("AWS_REGION", "us-east-1")
	awsEndpoint := getEnv("AWS_ENDPOINT", "") 

	awsBucket := getEnv("AWS_BUCKET", "fiap-videos-aldo-hackaton-2025")
	awsQueueURL := getEnv("AWS_QUEUE_URL", "https://sqs.us-east-1.amazonaws.com/629000537837/video-processing-queue")

	ctx := context.TODO()

	awsFactory := service.NewAWSClientFactory(
		ctx,
		awsRegion,
		awsEndpoint,
		"",
		"",
	)

	secretName := getEnv("DB_SECRET_NAME", "database-credentials20260218011702627300000001")
	var dbHost, dbUser, dbPassword, dbName, dbSslmode string

	creds, err := service.GetDatabaseSecrets(awsFactory, secretName)

	if err == nil {
		fmt.Println("✅ Credenciais carregadas do AWS Secrets Manager")
		dbHost = creds.Host
		dbUser = creds.Username
		dbPassword = creds.Password
		dbName = creds.Name
		dbSslmode = creds.Sslmode
	} else {
		fmt.Printf("⚠️ Erro ao acessar secret (%v). Usando variáveis locais.\n", err)
		dbHost = getEnv("DB_HOST", "localhost")
		dbUser = getEnv("DB_USER", "user")
		dbPassword = getEnv("DB_PASSWORD", "password")
		dbName = getEnv("DB_NAME", "fiapx_db")
		dbSslmode = getEnv("DB_SSL_MODE", "disable")
	}

	// Certifique-se que seu postgres.go está com sslmode=require para AWS e sslmode=disable para Local 
	db := database.SetupDatabase(dbHost, dbUser, dbPassword, dbName, dbSslmode)
	if db == nil {
		panic("❌ Falha crítica: Banco de dados não inicializado.")
	}
	db.AutoMigrate(&entity.User{}, &entity.Video{})

	storageService := service.NewStorageService(
		awsFactory.NewS3Client(),
		awsFactory.NewSQSClient(),
		awsBucket,
		awsQueueURL,
	)

	tokenService := service.NewTokenService()
	videoRepo := database.NewVideoRepository(db)
	userRepo := database.NewUserRepository(db)

	videoUC := usecase.NewVideoUseCase(videoRepo, userRepo, storageService, storageService)
	userUC := usecase.NewUserUseCase(userRepo, tokenService)

	authMiddleware := middleware.NewAuthMiddleware(tokenService)
	videoHandler := handler.NewVideoHandler(videoUC)
	authHandler := handler.NewAuthHandler(userUC)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "alive"})
	})

	r.GET("/ready", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(503, gin.H{"status": "unready", "database": "down"})
			return
		}
		c.JSON(200, gin.H{"status": "ready", "database": "up"})
	})

	r.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	setupRoutes(r, authHandler, videoHandler, authMiddleware)

	fmt.Printf("🚀 API rodando na porta 8080. Banco: %s\n", dbHost)
	r.Run(":8080")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func setupRoutes(r *gin.Engine, auth *handler.AuthHandler, video *handler.VideoHandler, mid *middleware.AuthMiddleware) {
	r.MaxMultipartMemory = 50 << 20
	r.Static("/static", "./web")

	r.GET("/", func(c *gin.Context) { c.File("./web/login.html") })
	r.GET("/dashboard", func(c *gin.Context) { c.File("./web/upload.html") })

	api := r.Group("/api")
	{
		api.POST("/register", auth.Register)
		api.POST("/login", auth.Login)

		protected := api.Group("/")
		protected.Use(mid.Handle())
		{
			protected.POST("/upload", video.UploadVideo)
			protected.GET("/videos", video.ListVideos)
			protected.GET("/videos/:id/download", video.GetDownloadLink)
		}
	}
}