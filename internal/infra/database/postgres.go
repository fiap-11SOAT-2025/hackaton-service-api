package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase(host, user, password, name string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, name)
	
	// Configuração de logs para mostrar apenas erros graves
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	
	if err != nil {
		log.Fatalf("❌ Erro ao conectar no Banco de Dados: %v", err)
	}

	log.Println("✅ Conexão com o banco de dados estabelecida!")
	return db
}