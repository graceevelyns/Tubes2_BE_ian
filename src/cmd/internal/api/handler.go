package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/graceevelyns/Tubes2_BE_ian/src/cmd/internal/algorithm"
	// "github.com/graceevelyns/Tubes2_BE_ian/src/internal/scraper"
	_ "github.com/graceevelyns/Tubes2_BE_ian/src/cmd/docs"
)

// RecipeSolution represents the solution containing found recipes
//
//	@swagger:api
type RecipeSolution struct {
	ElementName  string       `json:"elementName"`
	SearchParams SearchParams `json:"searchParams"`
	Found        bool         `json:"found"`
	SearchTimeMs float64      `json:"searchTimeMs"`
	NodesVisited int          `json:"nodesVisited"`
	Recipes      []RecipePath `json:"recipes"`
}

// SearchParams represents the search parameters
//
//	@swagger:api
type SearchParams struct {
	Algorithm string `json:"algorithm"`
	Count     int    `json:"count"`
	Mode      string `json:"mode"`
}

// RecipePath represents a single recipe path
//
//	@swagger:api
type RecipePath struct {
	NamaElemen    string         `json:"namaElemen"`
	IsBaseElement bool           `json:"isBaseElement"`
	DibuatDari    [][]RecipePath `json:"dibuatDari,omitempty"`
}

func transformRecipeTree(node *algorithm.RecipeTreeNode) RecipePath {
	if node == nil {
		return RecipePath{}
	}
	isBase := algorithm.IsBasicElement(*node)
	recipePath := RecipePath{
		NamaElemen:    node.NamaElemen,
		IsBaseElement: isBase,
	}
	if !isBase && len(node.DibuatDari) > 0 {
		recipePath.DibuatDari = make([][]RecipePath, 0, len(node.DibuatDari))
		for _, childNodePair := range node.DibuatDari {
			if childNodePair.LeftChild != nil && childNodePair.RightChild != nil {
				pair := []RecipePath{
					transformRecipeTree(childNodePair.LeftChild),
					transformRecipeTree(childNodePair.RightChild),
				}
				recipePath.DibuatDari = append(recipePath.DibuatDari, pair)
			}
		}
	}
	return recipePath
}

func countNodesInTree(node *algorithm.RecipeTreeNode) int {
	if node == nil {
		return 0
	}
	visited := make(map[string]bool)
	return countNodesRecursive(node, visited)
}

func countNodesRecursive(node *algorithm.RecipeTreeNode, visited map[string]bool) int {
	if node == nil {
		return 0
	}

	count := 1
	if visited[node.NamaElemen] {
		// return 0 // Aktifkan ini jika ingin jumlah elemen *unik* di pohon hasil
	}
	visited[node.NamaElemen] = true

	for _, childPair := range node.DibuatDari {
		count += countNodesRecursive(childPair.LeftChild, visited)
		count += countNodesRecursive(childPair.RightChild, visited)
	}
	return count
}

type SolveHandler struct {
	ElementNameToID map[string]int
	ElementIDToName map[int]string
}

func NewSolveHandler(nameToID map[string]int, idToName map[int]string) *SolveHandler {
	return &SolveHandler{
		ElementNameToID: nameToID,
		ElementIDToName: idToName,
	}
}

// @Summary		Get recipes for an element
// @Description	Finds recipes to create the specified element using either DFS or BFS algorithm
// @Tags			recipes
// @Accept			json
// @Produce		json
// @Param			element		query		string				true	"Element name to find recipes for"
// @Param			algorithm	query		string				false	"Search algorithm (dfs or bfs)"	Enums(dfs, bfs)	default(dfs)
// @Param			count		query		int					false	"Number of recipes to find"		minimum(1)		default(1)
// @Param			mode		query		string				false	"Search mode"					default(shortest)
// @Success		200			{object}	api.RecipeSolution	"Successful response"
// @Router			/solve-recipe [get]
func (sh *SolveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	elementNameQuery := strings.TrimSpace(queryParams.Get("element"))
	algo := strings.ToLower(strings.TrimSpace(queryParams.Get("algorithm")))
	countStr := strings.TrimSpace(queryParams.Get("count"))
	mode := strings.ToLower(strings.TrimSpace(queryParams.Get("mode")))

	if elementNameQuery == "" {
		http.Error(w, "Query parameter 'element' is required.", http.StatusBadRequest)
		return
	}

	if algo == "" {
		algo = "dfs"
	}
	if algo != "dfs" && algo != "bfs" {
		http.Error(w, "Invalid 'algorithm' parameter. Use 'dfs' or 'bfs'.", http.StatusBadRequest)
		return
	}

	var recipeCount int = 1
	if countStr != "" {
		var err error
		recipeCount, err = strconv.Atoi(countStr)
		if err != nil || recipeCount < 1 {
			http.Error(w, "Invalid 'count' parameter. Must be a positive integer.", http.StatusBadRequest)
			recipeCount = 1
		}
	}

	if mode == "" {
		mode = "shortest"
	}

	log.Printf("Handler /solve-recipe: element='%s', algorithm='%s', count=%d, mode='%s'", elementNameQuery, algo, recipeCount, mode)

	normalizedQueryName := strings.ToLower(elementNameQuery)
	startID, ok := sh.ElementNameToID[normalizedQueryName]

	elementNameForOutput := elementNameQuery
	if !ok {
		http.Error(w, fmt.Sprintf("Element '%s' not found in our database.", elementNameQuery), http.StatusNotFound)
		return
	}
	if officialName, nameOk := sh.ElementIDToName[startID]; nameOk {
		elementNameForOutput = officialName
	}

	var resultTree *algorithm.RecipeTreeNode
	var nodesVisitedByAlgo int

	startTime := time.Now()
	multParam := 1 // placeholder

	if algo == "dfs" {
		log.Printf("Calling DFS for element ID %d (%s), needFound: %d", startID, elementNameForOutput, recipeCount)
		if recipeCount == 1 {
			resultTree = algorithm.Dfs(startID, recipeCount, multParam)
		} else {
			resultTree = algorithm.ParallelDFS(startID, recipeCount)
		}
		// nodesVisitedByAlgo = hasilDariDFS.NodesExplored // contoh
	} else { // bfs
		log.Printf("Calling BFS for element ID %d (%s), needFound: %d", startID, elementNameForOutput, recipeCount)
		if recipeCount == 1 {
			resultTree = algorithm.Bfs(startID, recipeCount, multParam)
		} else {
			resultTree = algorithm.ParallelBfs(startID, recipeCount)
		}
		// nodesVisitedByAlgo = hasilDariBFS.NodesExplored // contoh
	}

	searchTimeMs := float64(time.Since(startTime).Microseconds()) / 1000.0

	if resultTree == nil || (resultTree.BanyakResep == 0 && !algorithm.IsBasicElement(*resultTree)) {
		isActuallyBase := false
		if resultTree != nil {
			isActuallyBase = algorithm.IsBasicElement(*resultTree)
		}
		if !isActuallyBase {
			log.Printf("Algorithm did not find a valid recipe tree for '%s' (ID: %d).", elementNameForOutput, startID)
			solution := RecipeSolution{
				ElementName:  elementNameForOutput,
				SearchParams: SearchParams{Algorithm: algo, Count: recipeCount, Mode: mode},
				Found:        false,
				SearchTimeMs: searchTimeMs,
				NodesVisited: nodesVisitedByAlgo,
				Recipes:      []RecipePath{},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(solution)
			return
		}
	}

	var recipePaths []RecipePath
	nodesVisitedInResultTree := 0
	if resultTree != nil {
		recipePaths = append(recipePaths, transformRecipeTree(resultTree))
		nodesVisitedInResultTree = countNodesInTree(resultTree)
	}

	finalNodesVisited := nodesVisitedByAlgo
	if finalNodesVisited == 0 && nodesVisitedInResultTree > 0 {
		finalNodesVisited = nodesVisitedInResultTree
	}

	solution := RecipeSolution{
		ElementName:  elementNameForOutput,
		SearchParams: SearchParams{Algorithm: algo, Count: recipeCount, Mode: mode},
		Found:        true,
		SearchTimeMs: searchTimeMs,
		NodesVisited: finalNodesVisited, // Gunakan nilai yang sudah ditentukan
		Recipes:      recipePaths,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(solution); err != nil {
		log.Printf("Error encoding solution to JSON for element %s: %v", elementNameForOutput, err)
		http.Error(w, "Failed to format solution data.", http.StatusInternalServerError)
	}
}
