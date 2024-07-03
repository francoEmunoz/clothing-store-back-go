package db

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB // Variable global para la instancia de la base de datos

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Configura el logger para que imprima en STDOUT
		logger.Config{
			SlowThreshold:             time.Second,   // Umbral para queries lentas
			LogLevel:                  logger.Silent, // Nivel de log (Silent para desactivar logs)
			IgnoreRecordNotFoundError: true,          // Ignora errores de "record not found"
			Colorful:                  true,          // Salida colorida
		},
	)

	// Conecta a la base de datos MySQL
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	} else {
		fmt.Println("Conexion exitosa")
	}

	// Asigna la instancia de la base de datos a la variable global
	Database = database
}
