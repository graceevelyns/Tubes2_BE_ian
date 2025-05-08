package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/graceevelyns/Tubes2_BE_ian/src/cmd/docs"
	"github.com/graceevelyns/Tubes2_BE_ian/src/internal/api"
	"github.com/graceevelyns/Tubes2_BE_ian/src/internal/model"
	"github.com/graceevelyns/Tubes2_BE_ian/src/internal/scraper"
)

var (
	globalAllNodes        map[string]*model.RecipeNode
	globalBaseElements    []*model.RecipeNode
	globalOrderedNodeKeys []string
)

type ElementOutputData struct {
	ID       int     `json:"Id"`
	Name     string  `json:"Name"`
	Tier     int     `json:"Tier"`
	FromPair [][]int `json:"FromPair"`
	CanMake  []int   `json:"CanMake"`
}

var standardBaseElementsList = []string{"Air", "Earth", "Fire", "Water", "Time"}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Memulai aplikasi Backend Little Alchemy Solver...")

	log.Println("Memulai scraping data dari Fandom Wiki...")
	var err error
	globalAllNodes, globalBaseElements, globalOrderedNodeKeys, err = scraper.FetchAndParseData()
	if err != nil {
		log.Fatalf("Gagal melakukan scraping data awal: %v", err)
	}
	if len(globalAllNodes) == 0 {
		log.Fatalf("Scraping tidak menghasilkan data node.")
	}
	if len(globalOrderedNodeKeys) != len(globalAllNodes) {
		log.Printf("PERINGATAN: Jumlah kunci terurut (%d) tidak sama dengan jumlah total node di peta (%d).", len(globalOrderedNodeKeys), len(globalAllNodes))
	}
	log.Printf("Scraping awal berhasil. Total node unik: %d. Jumlah kunci terurut: %d.\n", len(globalAllNodes), len(globalOrderedNodeKeys))

	api.InitData(globalAllNodes, globalBaseElements)

	r := mux.NewRouter()
	log.Println("Router Mux dibuat.")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/api/recipes/{elementName}", api.GetRecipesHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/graph-data", serveGraphDataAsIDScrapingOrder).Methods(http.MethodGet)

	port := "8080"
	log.Printf("Server API siap berjalan di http://localhost:%s\n", port)
	log.Printf("-> Akses /api/recipes/{nama}?algorithm=... untuk mencari resep.")
	log.Printf("-> Akses /graph-data untuk melihat data graf JSON (format ID, urutan scraping, dengan tier).")
	log.Printf("-> Akses /swagger/index.html untuk dokumentasi API interaktif.")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

// serveGraphDataAsIDScrapingOrder membuat JSON format ID, urut sesuai penemuan scraper
// Anotasi Swagger
// @Summary      Get All Processed Graph Data (Scraping Order with Tiers)
// @Description  Mengembalikan seluruh data elemen dan resep yang valid dalam format ID terstruktur, terurut berdasarkan penemuan saat scraping. Termasuk info Tier, FromPair, dan CanMake. Hanya elemen dasar atau yang punya resep yang disertakan.
// @Tags         Graph Data
// @Produce      json
// @Success      200 {array} main.ElementOutputData "Array data elemen dalam format ID dengan Tier"
// @Failure      500 {string} string "Error jika data graf belum siap"
// @Router       /graph-data [get]
func serveGraphDataAsIDScrapingOrder(w http.ResponseWriter, r *http.Request) {
	if len(globalAllNodes) == 0 {
		http.Error(w, "Data graf belum siap.", http.StatusInternalServerError)
		return
	}
	if len(globalOrderedNodeKeys) == 0 && len(globalAllNodes) > 0 {
		log.Println("Peringatan Kritis di Handler: globalOrderedNodeKeys kosong!")
		http.Error(w, "Data urutan node tidak tersedia.", http.StatusInternalServerError)
		return
	}

	nameToID := make(map[string]int)
	idDataMap := make(map[int]*ElementOutputData)
	nextID := 1

	log.Printf("Handler /graph-data: Memulai assignment ID berdasarkan globalOrderedNodeKeys (jumlah: %d)", len(globalOrderedNodeKeys))
	elementCountInMap := 0
	for _, name := range globalOrderedNodeKeys {
		if _, nodeExists := globalAllNodes[name]; nodeExists {
			currentID := nextID
			nameToID[name] = currentID
			idDataMap[currentID] = &ElementOutputData{ID: currentID, Name: name, Tier: -1, FromPair: make([][]int, 0), CanMake: make([]int, 0)}
			nextID++
			elementCountInMap++
		}
	}
	log.Printf("Handler /graph-data: Selesai assignment ID. %d node valid mendapatkan ID.", elementCountInMap)
	if elementCountInMap != len(globalAllNodes) {
		log.Printf("Peringatan di Handler: Jumlah node yang dapat di-map (%d) tidak sama dengan globalAllNodes (%d)!", elementCountInMap, len(globalAllNodes))
	}

	elementTiersByName := make(map[string]int)
	for name := range globalAllNodes {
		elementTiersByName[name] = -1
	}

	normalizedBaseElements := make(map[string]bool)
	for _, baseName := range standardBaseElementsList {
		normalizedBaseElements[strings.Title(strings.ToLower(baseName))] = true
	}

	for name := range globalAllNodes {
		if normalizedBaseElements[name] {
			elementTiersByName[name] = 0
		}
	}

	maxIterations := len(globalAllNodes) + 5
	for i := 0; i < maxIterations; i++ {
		changedInPass := false
		for elementName, nodeData := range globalAllNodes {
			if elementTiersByName[elementName] == 0 {
				continue
			}

			minTierForThisElement := -1

			if len(nodeData.DibuatDari) == 0 && !normalizedBaseElements[elementName] {
				continue
			}

			for _, recipePair := range nodeData.DibuatDari {
				if recipePair[0] == nil || recipePair[1] == nil {
					continue
				}
				ing1Name := recipePair[0].NamaElemen
				ing2Name := recipePair[1].NamaElemen

				tierIng1, ok1 := elementTiersByName[ing1Name]
				tierIng2, ok2 := elementTiersByName[ing2Name]

				if ok1 && tierIng1 != -1 && ok2 && tierIng2 != -1 {
					calculatedRecipeTier := 1 + max(tierIng1, tierIng2)
					if minTierForThisElement == -1 || calculatedRecipeTier < minTierForThisElement {
						minTierForThisElement = calculatedRecipeTier
					}
				}
			}

			if minTierForThisElement != -1 {
				if elementTiersByName[elementName] == -1 || minTierForThisElement < elementTiersByName[elementName] {
					elementTiersByName[elementName] = minTierForThisElement
					changedInPass = true
				}
			}
		}
		if !changedInPass {
			log.Printf("Kalkulasi tier konvergen pada iterasi ke-%d.", i+1)
			break
		}
		if i == maxIterations-1 {
			log.Println("Peringatan: Kalkulasi tier mencapai batas iterasi maksimum.")
		}
	}

	for currentIDFromMap, data := range idDataMap {
		if tier, ok := elementTiersByName[data.Name]; ok {
			data.Tier = tier
		} else {
			data.Tier = -1
		}
		idDataMap[currentIDFromMap] = data
	}

	canMakeTemp := make(map[int]map[int]bool)
	for resultName, resultNodeData := range globalAllNodes {
		resultID, rOk := nameToID[resultName]
		if !rOk {
			continue
		}
		elementOutputData, eOk := idDataMap[resultID]
		if !eOk {
			continue
		}

		isResultBase := false
		if node, exists := globalAllNodes[resultName]; exists {
			isResultBase = node.IsBaseElement
		}

		if !isResultBase && len(resultNodeData.DibuatDari) > 0 {
			processedPairs := make(map[string]bool)
			for _, recipePair := range resultNodeData.DibuatDari {
				if recipePair[0] == nil || recipePair[1] == nil {
					continue
				}
				ing1Name := recipePair[0].NamaElemen
				ing2Name := recipePair[1].NamaElemen
				ing1ID, i1OK := nameToID[ing1Name]
				ing2ID, i2OK := nameToID[ing2Name]
				if i1OK && i2OK {
					pairIDs := []int{ing1ID, ing2ID}
					if pairIDs[0] > pairIDs[1] {
						pairIDs[0], pairIDs[1] = pairIDs[1], pairIDs[0]
					}
					pairKey := fmt.Sprintf("%d-%d", pairIDs[0], pairIDs[1])
					if !processedPairs[pairKey] {
						elementOutputData.FromPair = append(elementOutputData.FromPair, pairIDs)
						processedPairs[pairKey] = true
					}
					if _, ok := canMakeTemp[ing1ID]; !ok {
						canMakeTemp[ing1ID] = make(map[int]bool)
					}
					canMakeTemp[ing1ID][resultID] = true
					if ing1ID != ing2ID {
						if _, ok := canMakeTemp[ing2ID]; !ok {
							canMakeTemp[ing2ID] = make(map[int]bool)
						}
						canMakeTemp[ing2ID][resultID] = true
					}
				}
			}
		}
	}
	for bahanID, hasilSet := range canMakeTemp {
		if elementOutputData, ok := idDataMap[bahanID]; ok {
			elementOutputData.CanMake = make([]int, 0, len(hasilSet))
			for hasilID := range hasilSet {
				elementOutputData.CanMake = append(elementOutputData.CanMake, hasilID)
			}
			sort.Ints(elementOutputData.CanMake)
		}
	}

	outputSlice := make([]*ElementOutputData, 0, len(idDataMap))
	for _, name := range globalOrderedNodeKeys {
		id, nameOk := nameToID[name]
		if !nameOk {
			continue
		}

		elementData, ok := idDataMap[id]
		if !ok {
			continue
		}

		isBase := false
		if originalNode, nodeExists := globalAllNodes[elementData.Name]; nodeExists {
			isBase = originalNode.IsBaseElement
		}
		hasRecipes := len(elementData.FromPair) > 0

		if isBase || hasRecipes {
			outputSlice = append(outputSlice, elementData)
		}
	}
	log.Printf("Jumlah elemen setelah filter untuk /graph-data: %d\n", len(outputSlice))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(outputSlice); err != nil {
		log.Printf("Error encoding ID data to JSON: %v", err)
		http.Error(w, "Gagal memformat data graf.", http.StatusInternalServerError)
	}
}
