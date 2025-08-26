package database

import (
	"chatbot/pkg/utils"
	"chatbot/pkg/utils/go-utils/encryptDecrypt"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Declare the variable for database
var DB *gorm.DB
var EsystemDB *gorm.DB

// ConnectDB connect to db
func ConnectToCAGABAYDB() {

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Only log slow SQL queries
			LogLevel:      logger.Silent, // Disable SQL logs
			Colorful:      false,         // Disable color
		},
	)

	var err error
	secretKey := utils.GetEnv("SECRET_KEY")
	//fmt.Println("SECRET KEY: ", secretKey)
	p, err := encryptDecrypt.Decrypt(utils.GetEnv("DB_PORT"), secretKey)
	if err != nil {
		return
	}
	h, err := encryptDecrypt.Decrypt(utils.GetEnv("DB_HOST"), secretKey)
	if err != nil {
		return
	}
	user, err := encryptDecrypt.Decrypt(utils.GetEnv("DB_USER"), secretKey)
	if err != nil {
		return
	}
	password, err := encryptDecrypt.Decrypt(utils.GetEnv("DB_PASSWORD"), secretKey)
	if err != nil {
		return
	}
	db, err := encryptDecrypt.Decrypt(utils.GetEnv("DB_NAME"), secretKey)
	if err != nil {
		return
	}

	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		log.Fatalf("Failed to parse DB_PORT: %v", err)
	}

	//Connection URL to connect to Postgres DB
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", h, port, user, password, db)

	//Connect to the DB and initialize the DB variable
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	// //Connect to the DB and initialize the DB variable
	// DB, err = gorm.Open(postgres.Open(dsn))

	if err != nil {
		log.Println(err)
		panic("Failed to connect to CA-GABAY database")
	}

	fmt.Print("Successfully connected to CA-GABAY database \n")
}

func ConnectToEsystemDB() {

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Only log slow SQL queries
			LogLevel:      logger.Silent, // Disable SQL logs
			Colorful:      false,         // Disable color
		},
	)

	var err error
	secretKey := utils.GetEnv("SECRET_KEY")

	esystemP, err := encryptDecrypt.Decrypt(utils.GetEnv("ESYSTEM_DB_PORT"), secretKey)
	if err != nil {
		return
	}
	esystemHost, err := encryptDecrypt.Decrypt(utils.GetEnv("ESYSTEM_DB_HOST"), secretKey)
	if err != nil {
		return
	}
	esystemUserAccount, err := encryptDecrypt.Decrypt(utils.GetEnv("ESYSTEM_DB_USER"), secretKey)
	if err != nil {
		return
	}
	esystemtPassword, err := encryptDecrypt.Decrypt(utils.GetEnv("ESYSTEM_DB_PASSWORD"), secretKey)
	if err != nil {
		return
	}
	esystemDbName, err := encryptDecrypt.Decrypt(utils.GetEnv("ESYSTEM_DB_NAME"), secretKey)
	if err != nil {
		return
	}

	esystemPort, err := strconv.ParseUint(esystemP, 10, 32)
	if err != nil {
		log.Fatalf("Failed to parse ESYSTEM_DB_PORT: %v", err)
	}

	//Connection URL to connect to Postgres DB
	esystemDsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", esystemHost, esystemPort, esystemUserAccount, esystemtPassword, esystemDbName)

	//Connect to the DB and initialize the DB variable
	EsystemDB, err = gorm.Open(postgres.Open(esystemDsn), &gorm.Config{
		Logger: gormLogger,
	})

	// //Connect to the DB and initialize the DB variable
	// DB, err = gorm.Open(postgres.Open(dsn))

	if err != nil {
		log.Println(err)
		panic("Failed to connect to eSystem database")
	}

	fmt.Print("Successfully connected to eSystem database \n")
}
