package models

type Song struct {
    ID          int    `json:"id"`
    GroupName   string `json:"group"`
    SongTitle   string `json:"song"`
    ReleaseDate string `json:"releaseDate"`
    Text        string `json:"text"`
    Link        string `json:"link"`
}
