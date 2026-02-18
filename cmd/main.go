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
	// 1. Configurações de Ambiente para acesso à AWS
	// Para rodar sem LocalStack e sem Docker, garanta que AWS_ACCESS_KEY_ID e AWS_SECRET_ACCESS_KEY
	// estejam exportadas no seu terminal.
	awsRegion := getEnv("AWS_REGION", "us-east-1")
	awsEndpoint := getEnv("AWS_ENDPOINT", "") // Vazio para conectar diretamente à AWS oficial
	awsAccessKeyID := getEnv("AWS_ACCESS_KEY_ID", "")
	awsSecretAccessKey := getEnv("AWS_SECRET_ACCESS_KEY", "")
	awsBucket := getEnv("AWS_BUCKET", "fiap-videos")
	awsQueueURL := getEnv("AWS_QUEUE_URL", "")

	ctx := context.TODO()
	awsFactory := service.NewAWSClientFactory(
		ctx,
		awsRegion,
		awsEndpoint,
		awsAccessKeyID,
		awsSecretAccessKey,
	)

	// 2. Busca credenciais do Banco de Dados no Secrets Manager
	// Como você quer usar o RDS remoto, o segredo fornecerá o HOST, USER e PASSWORD corretos.
	secretName := "database-credentials20260218011702627300000001"
	var dbHost, dbUser, dbPassword, dbName string

	creds, err := service.GetDatabaseSecrets(awsFactory, secretName)
	if err == nil {
		fmt.Println("✅ Conectando ao RDS remoto via AWS Secrets Manager")
		dbHost = creds.Host
		dbUser = creds.Username
		dbPassword = creds.Password
		dbName = creds.Name
	} else {
		// Log de erro caso não consiga buscar as credenciais da nuvem
		fmt.Printf("❌ Falha crítica: Não foi possível carregar os segredos da AWS: %v\n", err)
		// Caso ainda queira um fallback local por segurança:
		dbHost = getEnv("DB_HOST", "localhost")
		dbUser = getEnv("DB_USER", "user")
		dbPassword = getEnv("DB_PASSWORD", "password")
		dbName = getEnv("DB_NAME", "fiapx_db")
	}

	// 3. Inicialização do Banco de Dados (RDS Remoto)
	// Certifique-se de que o Security Group do RDS permite conexão externa (Porta 5432) do seu IP.
	db := database.SetupDatabase(dbHost, dbUser, dbPassword, dbName)
	db.AutoMigrate(&entity.User{}, &entity.Video{})

	// 4. Inicialização dos Serviços (S3 e SQS)
	storageService := service.NewStorageService(
		awsFactory.NewS3Client(),
		awsFactory.NewSQSClient(),
		awsBucket,
		awsQueueURL,
	)

	tokenService := service.NewTokenService()
	videoRepo := database.NewVideoRepository(db)
	userRepo := database.NewUserRepository(db)

	// 5. Casos de Uso e Handlers
	videoUC := usecase.NewVideoUseCase(videoRepo, userRepo, storageService, storageService)
	userUC := usecase.NewUserUseCase(userRepo, tokenService)

	authMiddleware := middleware.NewAuthMiddleware(tokenService)
	videoHandler := handler.NewVideoHandler(videoUC)
	authHandler := handler.NewAuthHandler(userUC)

	// 6. Configuração do Servidor Gin
	r := gin.Default()

	// Endpoints de Saúde
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

	// Swagger e Rotas
	r.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	setupRoutes(r, authHandler, videoHandler, authMiddleware)

	fmt.Printf("🚀 API iniciada com sucesso. Conectada ao banco: %s em %s\n", dbName, dbHost)
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