package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/graceevelyns/Tubes2_BE_ian/src/cmd/docs"
	"github.com/graceevelyns/Tubes2_BE_ian/src/cmd/internal/api"
	"github.com/graceevelyns/Tubes2_BE_ian/src/cmd/internal/scraper"
	"github.com/graceevelyns/Tubes2_BE_ian/src/cmd/internal/algorithm"
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
	log.Println("Memulai aplikasi Backend Little Alchemy Solver...")
	log.Println("Memulai scraping dan pemrosesan data...")
	var err error
	processedGraphData, err = scraper.FetchAndProcessData()
	if err != nil {
		log.Fatalf("Gagal melakukan scraping dan pemrosesan data: %v", err)
	}
	if processedGraphData == nil || len(processedGraphData) == 0 {
		log.Fatalf("Scraping dan pemrosesan tidak menghasilkan data graf yang valid atau datanya kosong.")
	}
	log.Printf("Scraping dan pemrosesan data berhasil. Jumlah elemen terproses untuk output: %d.\n", len(processedGraphData))

	initializeElementMaps(processedGraphData)

	algorithm.InitializeAlgorithmElements(scraper.GetProcessedElements())

	solveRecipe := api.NewSolveHandler(elementNameToID, elementIDToName)

	// router mux untuk menangani permintaan HTTP
	// dan menyajikan data graf yang sudah diproses

	//	@Summary		Get All Processed Graph Data
	//	@Description	Mengembalikan seluruh data elemen dan resep yang valid dalam format ID terstruktur.
	//	@Tags			Graph Data
	//	@Produce		json
	//	@Success		200	{array}		scraper.Element	"Array data elemen dalam format ID dengan Tier"
	//	@Failure		500	{string}	string			"Error jika data graf belum siap atau tidak valid"
	//	@Router			/graph-data [get]
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	})
	log.Println("Router Mux dibuat.")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/graph-data", serveGraphDataHandler).Methods(http.MethodGet)

	r.Handle("/solve-recipe", solveRecipe).Methods(http.MethodGet)


	// menyajikan API Swagger
	port := "8080"
	log.Printf("Server API siap berjalan di http://localhost:%s\n", port)
	log.Printf("-> Akses /graph-data untuk melihat data graf JSON.")
	log.Printf("-> Akses /solve-recipe untuk menggunakan solver resep (contoh: /solve-recipe?element=Brick&algorithm=dfs&count=1).")
	log.Printf("-> Akses /swagger/index.html untuk dokumentasi API interaktif.")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

// serveGraphDataHandler menyajikan data graf yang sudah diproses.
// Anotasi Swagger
//	@Summary		Get All Processed Graph Data
//	@Description	Mengembalikan seluruh data elemen dan resep yang valid dalam format ID terstruktur, terurut berdasarkan penemuan saat scraping. Termasuk info Tier, FromPair, dan CanMake. Hanya elemen dasar atau yang punya resep valid dan tier terhitung yang disertakan.
//	@Tags			Graph Data
//	@Produce		json
//	@Success		200	{array}		scraper.Element	"Array data elemen dalam format ID dengan Tier"	//	Tipe	di	swagger	diupdate
//	@Failure		500	{string}	string			"Error jika data graf belum siap atau tidak valid"
//	@Router			/graph-data [get]
func serveGraphDataHandler(w http.ResponseWriter, r *http.Request) {
	if processedGraphData == nil || len(processedGraphData) == 0 {
		log.Println("Handler /graph-data: processedGraphData belum siap atau kosong.")
		http.Error(w, "Data graf belum siap atau tidak ada data yang valid.", http.StatusInternalServerError)
		return
	}

	log.Printf("Handler /graph-data: Menyajikan %d elemen terproses.", len(processedGraphData))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(processedGraphData); err != nil {
		log.Printf("Error encoding processed graph data to JSON: %v", err)
		http.Error(w, "Gagal memformat data graf.", http.StatusInternalServerError)
	}
}
