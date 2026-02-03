package main

import (
	"hackaton-service-api/internal/entity"
	"hackaton-service-api/internal/handler"
	"hackaton-service-api/internal/infra/service"
	"hackaton-service-api/internal/infra/database"
	"hackaton-service-api/internal/middleware"
	"hackaton-service-api/internal/usecase"
	"os"
	"fmt"
	"context"

	"github.com/gin-gonic/gin"
)

func main() {

	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "user")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "fiapx_db")
	
	awsRegion := getEnv("AWS_REGION", "us-east-1")
	awsEndpoint := getEnv("AWS_ENDPOINT", "http://localhost:4566")
	awsAccessKeyID := getEnv("AWS_ACCESS_KEY_ID", "teste")
	awsSecretAccessKey := getEnv("AWS_SECRET_ACCESS_KEY", "teste")
	awsBucket := getEnv("AWS_BUCKET", "fiap-videos")
	awsQueueURL := getEnv("AWS_QUEUE_URL", "http://localhost:4566/000000000000/video-processing-queue")
	
	db := database.SetupDatabase(dbHost, dbUser, dbPassword, dbName)
	db.AutoMigrate(&entity.User{}, &entity.Video{})

	awsFactory := service.NewAWSClientFactory(
        context.TODO(),
        awsRegion,
        awsEndpoint,
        awsAccessKeyID,
        awsSecretAccessKey,
    )

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
	setupRoutes(r, authHandler, videoHandler, authMiddleware)

	fmt.Println("🚀 API rodando na porta 8080")
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