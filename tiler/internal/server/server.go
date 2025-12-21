package server

import (
	"net/http"
	"strings"

	"tiler/internal/services"

	"github.com/rs/cors"
)

func NewServer(imageService services.ImageService) http.Handler {
	handler := NewHandler(imageService)

	mux := http.NewServeMux()

	// DZI endpoint
	mux.HandleFunc("/dzi/", func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, это запрос на DZI XML или на тайл
		// Если путь содержит "/files/", то это запрос на тайл
		if strings.Contains(r.URL.Path, "/files/") {
			// Это запрос на тайл
			handler.GetTile(w, r)
		} else {
			// Это запрос на DZI XML
			handler.GetDZI(w, r)
		}
	})

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	return c.Handler(mux)
}
