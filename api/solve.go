package handler

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/graceevelyns/Tubes2_BE_ian/src/cmd/pkg/algorithm"
	originalApi "github.com/graceevelyns/Tubes2_BE_ian/src/cmd/pkg/api"
	"github.com/graceevelyns/Tubes2_BE_ian/src/cmd/pkg/scraper"
)

var (
	solveHandlerInstance *originalApi.SolveHandler
	initOnce             sync.Once

	processedGraphData []*scraper.Element
	elementNameToID    map[string]int
	elementIDToName    map[int]string
)

func initializeDependencies() {
	log.Println("Vercel Handler: Initializing dependencies...")
	var err error

	processedGraphData, err = scraper.FetchAndProcessData()
	if err != nil {
		log.Fatalf("Vercel Handler: Error fetching and processing data: %v", err)
		return
	}
	if processedGraphData == nil || len(processedGraphData) == 0 {
		log.Fatalf("Vercel Handler: Invalid or empty data from scraper.")
		return
	}

	elementNameToID = make(map[string]int)
	elementIDToName = make(map[int]string)
	for _, el := range processedGraphData {
		normalizedName := strings.ToLower(el.Name)
		elementNameToID[normalizedName] = el.ID
		elementIDToName[el.ID] = el.Name
	}

	if scraper.GetProcessedElements() == nil || len(scraper.GetProcessedElements()) == 0 {
		algorithm.InitializeAlgorithmElements(scraper.GetProcessedElements())
	}

	solveHandlerInstance = originalApi.NewSolveHandler(elementNameToID, elementIDToName)
	log.Println("Vercel Handler: Dependencies initialized.")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	initOnce.Do(initializeDependencies)

	origin := r.Header.Get("Origin")
	allowedOrigins := []string{
		"https://tubes2-fe-ian.vercel.app",
		"http://localhost:3000",       
	}
	for _, o := range allowedOrigins {
		if origin == o {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			break
		}
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if solveHandlerInstance == nil {
		log.Println("Vercel Handler: solveHandlerInstance is nil, initialization might have failed.")
		http.Error(w, "Service not available due to initialization error", http.StatusInternalServerError)
		return
	}

	solveHandlerInstance.ServeHTTP(w, r)
}
