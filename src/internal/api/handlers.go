package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	alg "github.com/graceevelyns/Tubes2_BE_ian/src/internal/algorithm"
	"github.com/graceevelyns/Tubes2_BE_ian/src/internal/model"
)

var (
	allNodesData     map[string]*model.RecipeNode
	baseElementsData []*model.RecipeNode
)

func InitData(nodes map[string]*model.RecipeNode, base []*model.RecipeNode) {
	allNodesData = nodes
	baseElementsData = base
	if len(allNodesData) > 0 {
		log.Println("API Handler: Data graf siap digunakan.")
	} else {
		log.Println("API Handler: Peringatan - data graf kosong saat inisialisasi.")
	}
}

// RecipeResponse adalah struktur untuk JSON API response
type RecipeResponse struct {
	ElementName  string                 `json:"elementName" example:"Steam"`  // nama elemen yang dicari resepnya
	SearchParams map[string]interface{} `json:"searchParams"`                 // parameter pencarian yang digunakan (echo)
	Found        bool                   `json:"found" example:"true"`         // status apakah resep ditemukan
	SearchTimeMs float64                `json:"searchTimeMs" example:"15.23"` // waktu eksekusi pencarian dalam milidetik
	NodesVisited int                    `json:"nodesVisited" example:"85"`    // jumlah node graf yang dikunjungi algoritma
	Recipes      []*model.RecipeNode    `json:"recipes"`                      // array berisi pohon resep yang ditemukan (struktur RecipeNode)
}

// GetRecipesHandler godoc
// @Summary      Mencari Resep Elemen
// @Description  Mencari satu atau lebih jalur resep untuk membuat elemen target menggunakan algoritma BFS atau DFS
// @Tags         Recipes
// @Accept       json
// @Produce      json
// @Param        elementName path string true "Nama Elemen Target (case-insensitive)" example(Steam)
// @Param        algorithm query string true "Algoritma Pencarian ('bfs' atau 'dfs')" Enums(bfs, dfs)
// @Param        mode query string false "Mode Pencarian ('shortest' atau 'multiple')" Enums(shortest, multiple) default(shortest)
// @Param        count query int false "Jumlah resep maks (jika mode = multiple, default = 5)" default(1) minimum(1)
// @Success      200 {object} api.RecipeResponse "Respon sukses dengan hasil pencarian resep"
// @Failure      400 {object} map[string]string "Error: Input parameter tidak valid"
// @Failure      404 {object} map[string]string "Error: Elemen target tidak ditemukan di data graf"
// @Failure      500 {object} map[string]string "Error: Masalah internal server"
// @Router       /api/recipes/{elementName} [get]

func GetRecipesHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// 1. ambil path parameter {elementName} dari Mux
	vars := mux.Vars(r)
	elementNameReq := vars["elementName"]
	elementName := strings.Title(strings.ToLower(elementNameReq))

	// periksa apakah elemen target ada di data graf
	if _, exists := allNodesData[elementName]; !exists {
		errMsg := fmt.Sprintf("Elemen '%s' tidak ditemukan dalam data.", elementName)
		log.Println("API Handler:", errMsg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
		return
	}

	// 2. ambil query parameters & validasi
	algorithm := strings.ToLower(r.URL.Query().Get("algorithm"))
	mode := strings.ToLower(r.URL.Query().Get("mode"))
	countStr := r.URL.Query().Get("count")

	if mode == "" {
		mode = "shortest"
	}
	maxCount := 1
	var err error
	if mode == "multiple" {
		if countStr != "" {
			maxCount, err = strconv.Atoi(countStr)
			if err != nil || maxCount <= 0 {
				errMsg := "Parameter 'count' harus angka positif jika mode = 'multiple'"
				log.Println("API Handler: Error -", errMsg)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
				return
			}
		} else {
			maxCount = 5
			log.Printf("Parameter 'count' kosong untuk mode multiple, menggunakan default %d.", maxCount)
		}
	}
	if mode == "shortest" {
		maxCount = 1
	}

	if algorithm != "bfs" && algorithm != "dfs" {
		errMsg := "Parameter 'algorithm' harus 'bfs' atau 'dfs'"
		log.Println("API Handler: Error, ", errMsg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
		return
	}
	if mode != "shortest" && mode != "multiple" {
		errMsg := "Parameter 'mode' harus 'shortest' atau 'multiple'"
		log.Println("API Handler: Error, ", errMsg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
		return
	}

	if len(allNodesData) == 0 {
		log.Println("API Handler: Error! Data graf (allNodesData) kosong.")
		http.Error(w, "Data server belum siap.", http.StatusInternalServerError)
		return
	}

	log.Printf("API Handler: Memproses permintaan: Elemen = %s, Algoritma = %s, Mode = %s, Count = %d\n", elementName, algorithm, mode, maxCount)

	// 3. panggil algoritma pencarian resep
	var found bool
	var nodesVisited int
	var recipesResult []*model.RecipeNode = make([]*model.RecipeNode, 0)

	switch algorithm {
	case "bfs":
		tree, visited, success := alg.FindShortestRecipeBFS(elementName, allNodesData, baseElementsData)
		nodesVisited = visited
		found = success
		if found && tree != nil {
			recipesResult = []*model.RecipeNode{tree}
		}
		if mode == "multiple" { // idk apakah bfs bisa multiple
			log.Printf("Peringatan: Mode 'multiple' diminta dengan BFS. Hanya 1 resep terpendek dikembalikan.")
		}

	case "dfs":
		trees, visited, success := alg.FindMultipleRecipesDFS(elementName, maxCount, allNodesData, baseElementsData)
		nodesVisited = visited
		found = success
		if found && trees != nil {
			recipesResult = trees
		}
	}

	duration := time.Since(startTime)

	// 4. response JSON
	response := RecipeResponse{
		ElementName: elementName,
		SearchParams: map[string]interface{}{
			"algorithm": algorithm, "mode": mode, "count": maxCount,
		},
		Found:        found,
		SearchTimeMs: float64(duration.Microseconds()) / 1000.0,
		NodesVisited: nodesVisited,
		Recipes:      recipesResult,
	}
	if response.Recipes == nil {
		response.Recipes = make([]*model.RecipeNode, 0)
	}

	// 5. kirim response JSON
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("API Handler: Error encoding JSON response: %v", err)
	}
}
