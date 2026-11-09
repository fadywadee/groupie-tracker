package main

import (
	"groupie-tracker/data"
	"groupie-tracker/fetcher"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// App Struct
type App struct {
	homeTmpl   *template.Template
	artistTmpl *template.Template
	artistsMap map[int]data.DetailedArtist
}

// Main Function
func main() {
	const baseURL = "https://groupietrackers.01edu.ai/api"
	app := &App{}

	// We'll initialize the app and router here

	//_ = app

	appData, err := fetcher.FetchAllData(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	app.artistsMap = data.WeaveData(
		appData.Artists,
		appData.Locations.Index,
		appData.Dates.Index,
		appData.Relations.Index,
	)

	app.homeTmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	app.artistTmpl, err = template.ParseFiles("templates/artist.html")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.homeHandler)

	mux.HandleFunc("GET /artist", app.artistHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))

}

func (app *App) homeHandler(

	w http.ResponseWriter,
	r *http.Request,

) {
	artistsList := make([]data.DetailedArtist, 0, len(app.artistsMap))

	for _, artist := range app.artistsMap {
		artistsList = append(artistsList, artist)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := app.homeTmpl.Execute(w, artistsList); err != nil {
		log.Printf("Home template execution error: %v", err)
		return
	}
}

func (app *App) artistHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	// 1. Get artist ID
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		http.Error(w, "Missing artist ID", http.StatusBadRequest)
		return
	}

	// 2. Convert ID to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	// 3. Find artist in masterMap
	artist, exists := app.artistsMap[id]
	if !exists {
		http.NotFound(w, r)
		return
	}

	// 4. Set Content-Type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// 5. Execute template
	if err := app.artistTmpl.Execute(w, artist); err != nil {
		log.Printf("Artist template execution error: %v", err)
		return
	}
}
