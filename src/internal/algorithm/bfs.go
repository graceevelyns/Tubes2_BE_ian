// INI STUB DOANG
package algorithm

import (
	"log"

	"github.com/graceevelyns/Tubes2_BE_ian/src/internal/model"
)

func FindShortestRecipeBFS(targetName string, allNodes map[string]*model.RecipeNode, baseElements []*model.RecipeNode) (*model.RecipeNode, int, bool) {
	log.Printf("[STUB] BFS: Mencari resep terpendek untuk %s\n", targetName)

	if targetName == "Steam" {
		waterNode, wExists := allNodes["Water"]
		fireNode, fExists := allNodes["Fire"]

		if wExists && fExists {
			log.Println("[STUB] BFS: Menemukan resep hardcoded untuk Steam (Water + Fire).")
			resultTree := &model.RecipeNode{
				NamaElemen:    "Steam",
				IsBaseElement: false,
				DibuatDari: [][2]*model.RecipeNode{
					{
						waterNode,
						fireNode,
					},
				},
			}
			return resultTree, 3, true
		} else {
			log.Printf("[STUB] BFS: Gagal membuat resep Steam hardcoded karena node Water (%t) atau Fire (%t) tidak ditemukan di allNodes.", wExists, fExists)
		}
	}

	log.Printf("[STUB] BFS: Tidak ditemukan resep (stub) untuk %s.\n", targetName)
	return nil, 1, false
}
