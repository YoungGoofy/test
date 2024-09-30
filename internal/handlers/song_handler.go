package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/YoungGoofy/MusicLib/internal/db"
	"github.com/YoungGoofy/MusicLib/internal/models"
	"github.com/gin-gonic/gin"
)

// AddSong godoc
// @Summary Add a new song
// @Description Add a new song to the library and enrich it with data from an external API.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body models.Song true "Song to add"
// @Success 201 {object} models.Song
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Server error"
// @Router /songs [post]
func AddSong(c *gin.Context) {
	var song models.Song
	if err := c.ShouldBindBodyWithJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	
	// Обрабатываем запрос к внешнему API для получения песен
	externalAPI := os.Getenv("EXTERNAL_API_URL")
	resp, err := http.Get(externalAPI + "&group=" + song.GroupName + "?song=" + song.SongTitle)
	if err != nil {
		log.Println("Failed to call external API: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call external API"})
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("External API returns non-200 status: ", resp.Status)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get song details"})
		return
	}

	var songDetail struct {
        ReleaseDate string `json:"releaseDate"`
        Text        string `json:"text"`
        Link        string `json:"link"`
    }

	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
        log.Println("Failed to decode external API response:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode external API response"})
        return
    }

	// Добавляем данные песни в модель базы данных
	song.ReleaseDate = songDetail.ReleaseDate
	song.Text = songDetail.Text
	song.Link = songDetail.Link

	// Добавляем данные в базу данных
	_, err = db.DB.Exec("INSERT INTO songs (group_name, song_title, release_date, text, link) VALUES ($1, $2, $3, $4, $5)",
        song.GroupName, song.SongTitle, song.ReleaseDate, song.Text, song.Link)

    if err != nil {
        log.Println("Failed to insert song:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert song"})
        return
    }

    c.Status(http.StatusCreated)
}

// DeleteSong godoc
// @Summary Delete a song by ID
// @Description Delete a song from the library by its ID.
// @Tags songs
// @Param id path int true "Song ID"
// @Success 200 {string} string "Song deleted"
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Song not found"
// @Failure 500 {string} string "Server error"
// @Router /songs/{id} [delete]
func DeleteSong(c *gin.Context) {
	id := c.Param(":id")

	_, err := db.DB.Exec("delete from songs where id=$1", id)
	if err != nil {
		log.Println("Failed to delete song:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
        return
	}
	c.Status(http.StatusOK)
}

// GetSongByID godoc
// @Summary Get a song by ID
// @Description Get a song from the library by its ID, including the song text.
// @Tags songs
// @Produce  json
// @Param id path int true "Song ID"
// @Success 200 {object} models.Song
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Song not found"
// @Failure 500 {string} string "Server error"
// @Router /songs/{id} [get]
func GetSongById(c *gin.Context) {
	id := c.Param(":id")

	var song models.Song
	err := db.DB.QueryRow("SELECT id, group_name, song_title, release_date, text, link FROM songs WHERE id=$1", id).
	Scan(&song.ID, &song.GroupName, &song.SongTitle, &song.ReleaseDate, &song.Text, &song.Link)
	
	if err != nil {
		log.Println("Failed to fetch song:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song"})
		return
	}
	
	c.JSON(http.StatusOK, song)
}

// UpdateSong godoc
// @Summary Update a song by ID
// @Description Update the details of a song in the library by its ID.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Song data to update"
// @Success 200 {object} models.Song
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Song not found"
// @Failure 500 {string} string "Server error"
// @Router /songs/{id} [put]
func UpdateSong(c *gin.Context) {
	id := c.Param("id")

    var song models.Song
    if err := c.ShouldBindJSON(&song); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    _, err := db.DB.Exec("UPDATE songs SET group_name=$1, song_title=$2, release_date=$3, text=$4, link=$5 WHERE id=$6",
        song.GroupName, song.SongTitle, song.ReleaseDate, song.Text, song.Link, id)

    if err != nil {
        log.Println("Failed to update song:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
        return
    }
	
    c.Status(http.StatusOK)
}

// GetSongs godoc
// @Summary Get list of songs
// @Description Get songs with optional filtering and pagination.
// @Tags songs
// @Produce  json
// @Success 200 {array} models.Song
// @Failure 500 {string} string "Server error"
// @Router /songs [get]
func GetSongs(c *gin.Context) {
	var songs []models.Song

    rows, err := db.DB.Query("SELECT id, group_name, song_title, release_date, text, link FROM songs")
    if err != nil {
        log.Println("Failed to fetch songs:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch songs"})
        return
    }
    defer rows.Close()

    for rows.Next() {
        var song models.Song
        if err := rows.Scan(&song.ID, &song.GroupName, &song.SongTitle, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
            log.Println("Failed to scan song:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan song"})
            return
        }
        songs = append(songs, song)
    }

    c.JSON(http.StatusOK, songs)
}