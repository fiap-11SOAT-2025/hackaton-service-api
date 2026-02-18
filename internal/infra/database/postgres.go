package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func SetupDatabase(host, user, password, dbName string) *gorm.DB {
	// ATERAÇÃO AQUI: Mudamos sslmode de "disable" para "require"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=require TimeZone=America/Sao_Paulo", 
		host, user, password, dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[error] failed to initialize database, got error %v", err)
		// É importante retornar nil ou dar panic aqui se não conectar,
		// mas seguindo seu padrão atual:
		return nil 
	}
	
	return db
}