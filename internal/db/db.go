package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// Делаем подключение к базе данных
func InitDB() {
	var err error
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Connection not opened: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database: ", err)
	}

	log.Println("Database connect successful")

	applyMigrations()
}

func applyMigrations() {
    // Получаем драйвер PostgreSQL для миграций
    driver, err := postgres.WithInstance(DB, &postgres.Config{})
    if err != nil {
        log.Fatal("Failed to create migration driver:", err)
    }

    // Указываем путь к миграциям
    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations",
        "postgres",
        driver,
    )
    if err != nil {
        log.Fatal("Failed to initialize migrate:", err)
    }

    // Применяем миграции
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal("Failed to apply migrations:", err)
    }

    log.Println("Migrations applied successfully")
}

