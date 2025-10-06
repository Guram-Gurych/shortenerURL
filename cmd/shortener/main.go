package main

import (
	"github.com/Guram-Gurych/shortenerURL.git/internal/handler"
	"github.com/Guram-Gurych/shortenerURL.git/internal/repository"
	"github.com/Guram-Gurych/shortenerURL.git/internal/service"
	"net/http"
)

func main() {
	var baseURL = "http://localhost:8080"
	var serverAddr = ":8080"

	repo := repository.NewMemoryRepository()
	serv := service.NewShortenerService(repo)
	hndl := handler.NewHandler(serv, baseURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			hndl.PostHandler(w, r)
		case http.MethodGet:
			hndl.GetHandler(w, r)
		default:
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
	})

	err := http.ListenAndServe(serverAddr, mux)
	if err != nil {
		panic(err)
	}
}
