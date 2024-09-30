package app

import (
	"log"
	"net/http"

	"github.com/YoungGoofy/MusicLib/internal/db"
	"github.com/YoungGoofy/MusicLib/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func App() {
	// Загружаем .env
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Инициализация базы данных
    db.InitDB()

	r := gin.Default()

	r.POST("/songs", handlers.AddSong)
	r.DELETE("/songs/:id", handlers.DeleteSong)
	r.GET("/songs/:id", handlers.GetSongById)
	r.PUT("/songs/:id", handlers.UpdateSong)
	r.GET("/songs", handlers.GetSongs)

	log.Println("Starting server on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}