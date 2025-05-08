// INI STUB DOANG
package algorithm

import (
	"log"
	"github.com/graceevelyns/Tubes2_BE_ian/src/internal/model"
)

func FindMultipleRecipesDFS(targetName string, maxCount int, allNodes map[string]*model.RecipeNode, baseElements []*model.RecipeNode) ([]*model.RecipeNode, int, bool) {
	log.Printf("[STUB] DFS: Mencari max %d resep untuk %s\n", maxCount, targetName)

	var resultTrees []*model.RecipeNode = make([]*model.RecipeNode, 0)
	nodesVisited := 0
	found := false

	if targetName == "Brick" && maxCount > 0 {
		_, mExists := allNodes["Mud"]
		_, fExists := allNodes["Fire"]
		_, eExists := allNodes["Earth"]
		_, wExists := allNodes["Water"]

		if mExists && fExists && eExists && wExists {
			log.Println("[STUB] DFS: Membuat resep hardcoded untuk Brick (Mud + Fire, Mud = Earth + Water).")

			earthStubNode := &model.RecipeNode{NamaElemen: "Earth", IsBaseElement: true}
			waterStubNode := &model.RecipeNode{NamaElemen: "Water", IsBaseElement: true}
			fireStubNode := &model.RecipeNode{NamaElemen: "Fire", IsBaseElement: true}

			mudStubNode := &model.RecipeNode{
				NamaElemen:    "Mud",
				IsBaseElement: false,
				DibuatDari:    [][2]*model.RecipeNode{{earthStubNode, waterStubNode}},
			}

			brickTree := &model.RecipeNode{
				NamaElemen:    "Brick",
				IsBaseElement: false,
				DibuatDari:    [][2]*model.RecipeNode{{mudStubNode, fireStubNode}},
			}

			resultTrees = append(resultTrees, brickTree)
			nodesVisited = 5
			found = true
		} else {
			log.Printf("[STUB] DFS: Gagal membuat resep Brick hardcoded karena node Mud (%t) / Fire (%t) / Earth (%t) / Water (%t) tidak ditemukan di allNodes.", mExists, fExists, eExists, wExists)
		}
	}

	if !found {
		log.Printf("[STUB] DFS: Tidak ditemukan resep (stub) untuk %s.\n", targetName)
		nodesVisited = 1
	}

	return resultTrees, nodesVisited, found
}
