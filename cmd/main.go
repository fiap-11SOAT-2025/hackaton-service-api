package main

import (
	"fiapx-api/internal/entity"
	"fiapx-api/internal/handler"
	"fiapx-api/internal/infra/service"
	"fiapx-api/internal/infra/database"
	"fiapx-api/internal/middleware"
	"fiapx-api/internal/usecase"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName:= os.Getenv("DB_NAME")
    if dbHost == "" {
        dbHost = "localhost"
    }
	if dbUser == "" {
        dbUser = "user"
    }
	if dbPassword == "" {
        dbPassword = "password"
    }
	if dbName == "" {
        dbName = "fiapx_db"
    }
	
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", dbHost, dbUser, dbPassword, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar no banco")
	}
	db.AutoMigrate(&entity.User{}, &entity.Video{})

	awsService, err := service.NewAWSService()
	if err != nil {
		log.Fatalf("Erro ao configurar AWS: %v", err)
	}
	tokenService := service.NewTokenService()

	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	videoRepo := database.NewVideoRepository(db)
	userRepo := database.NewUserRepository(db)

	videoUC := usecase.NewVideoUseCase(videoRepo, awsService, awsService)
	userUC := usecase.NewUserUseCase(userRepo, tokenService)

	videoHandler := handler.NewVideoHandler(videoUC)
	authHandler := handler.NewAuthHandler(userUC)

	r := gin.Default()
	r.MaxMultipartMemory = 50 << 20 // 50MB

	r.Static("/static", "./web")
	r.GET("/", func(c *gin.Context) {
		if _, err := os.Stat("./web/login.html"); os.IsNotExist(err) {
			c.String(404, "Frontend não encontrado")
			return
		}
		c.File("./web/login.html")
	})
	r.GET("/dashboard", func(c *gin.Context) {
		c.File("./web/upload.html")
	})

	api := r.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		
		protected := api.Group("/")
		protected.Use(authMiddleware.Handle())
		{
			protected.POST("/upload", videoHandler.UploadVideo)
			protected.GET("/videos", videoHandler.ListVideos)
			protected.GET("/videos/:id/download", videoHandler.GetDownloadLink)
		}
	}

	fmt.Println("🚀 API rodando na porta 8080")
	r.Run(":8080")
}