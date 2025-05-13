package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/graceevelyns/Tubes2_BE_ian/src/cmd/docs"
	"github.com/graceevelyns/Tubes2_BE_ian/src/cmd/pkg/algorithm"
	"github.com/graceevelyns/Tubes2_BE_ian/src/cmd/pkg/api"
	"github.com/graceevelyns/Tubes2_BE_ian/src/cmd/pkg/scraper"
)

var (
	processedGraphData []*scraper.Element
	elementNameToID    map[string]int
	elementIDToName    map[int]string
)

func initializeElementMaps(elements []*scraper.Element) {
	elementNameToID = make(map[string]int)
	elementIDToName = make(map[int]string)
	for _, el := range elements {
		normalizedName := strings.ToLower(el.Name)
		elementNameToID[normalizedName] = el.ID
		elementIDToName[el.ID] = el.Name
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	processedGraphData, err = scraper.FetchAndProcessData()
	if err != nil {
		log.Fatalf("error in scraping: %v", err)
	}
	if processedGraphData == nil || len(processedGraphData) == 0 {
		log.Fatalf("invalid or empty data.")
	}

	initializeElementMaps(processedGraphData)

	algorithm.InitializeAlgorithmElements(scraper.GetProcessedElements())

	solveRecipe := api.NewSolveHandler(elementNameToID, elementIDToName)

	// router mux for handling API requests

	//	@Summary		Get All Processed Graph Data
	//	@Description	For testing purposes, this endpoint returns all processed graph data in JSON format
	//	@Tags			Graph Data
	//	@Produce		json
	//	@Success		200	{array}		scraper.Element "Array element data in JSON format"
	//	@Failure		500	{string}	string			"Error if graph data is not ready or invalid"
	//	@Router			/graph-data [get]
	r := mux.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow specific origins in production
			origin := r.Header.Get("Origin")
			allowedOrigins := []string{
				"https://tubes2-fe-ian.vercel.app",
				"http://localhost:3000", // for local development
			}

			for _, o := range allowedOrigins {
				if origin == o {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/graph-data", serveGraphDataHandler).Methods(http.MethodGet)

	r.Handle("/solve-recipe", solveRecipe).Methods(http.MethodGet)

	// menyajikan API Swagger
	port := "8080"
	log.Printf("API server ready in http://localhost:%s\n", port)
	log.Printf("-> access /graph-data to see graph data in JSON.")
	log.Printf("-> access /solve-recipe to use solver for a recipe.")
	log.Printf("-> access /swagger/index.html for interactive API documentation.")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("error: %v", err)
	}
}

// serveGraphDataHandler menyajikan data graf yang sudah diproses.
// Anotasi Swagger
//
//	@Summary		Get All Processed Graph Data
//	@Description	For testing purposes, this endpoint returns all processed graph data in JSON format
//	@Tags			Graph Data
//	@Produce		json
//	@Success		200	{array}		scraper.Element "Array element data in JSON format"
//	@Failure		500	{string}	string			"Error if graph data is not ready or invalid"
//	@Router			/graph-data [get]
func serveGraphDataHandler(w http.ResponseWriter, r *http.Request) {
	if processedGraphData == nil || len(processedGraphData) == 0 {
		log.Println("/graph-data: processedGraphData is not ready or empty.")
		http.Error(w, "graph data is not ready.", http.StatusInternalServerError)
		return
	}

	log.Printf("/graph-data: %d proccessed.", len(processedGraphData))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(processedGraphData); err != nil {
		log.Printf("error encoding processed graph data ke JSON: %v", err)
		http.Error(w, "error in formatting graph data.", http.StatusInternalServerError)
	}
}
