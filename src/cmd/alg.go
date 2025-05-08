package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sort"
)


type RecipeTreeNode struct{
    namaElemen string;
    dibuatDari [][2]*RecipeTreeNode;
}

type Recipe struct{
    // ditemukan bool;
    // stepByStep []int;
    banyakNode int;
    resep *RecipeTreeNode
}

type RecipeList struct{
    namaElemen string;
    kumpulanResep []Recipe;
	banyakResep int;
}



func checkBasic(start int, needFound int, mult int) RecipeList {
	if elements[start].Name == "water" || elements[start].Name == "fire" || elements[start].Name == "earth" || elements[start].Name == "air" {
		return RecipeList{
			namaElemen: elements[start].Name,
			buatDari:   nil,
			banyakResep: 1,
		}
	}
	return dfs(el, start, needFound, mult)
}


visit := make([]bool, len(elements))
func getDFS(start int, needFound int) RecipeList {

	for i := 0; i < len(visit); i++ {
		visit[i] = false
	}
	return dfs(start, needFound, 1)
}

func dfs(start int, needFound int, mult int) RecipeList {
	q := elements[start].FromPair
	notFound := 0
	retResep := RecipeList{
		namaElemen: elements[start].Name,
		kumpulanResep: make([]Recipe, 0),
		banyakResep: 0,
	}
	if(visit[start]) {
		return retResep
	}
	visit[start] = true
	for pair := range q {
		resepLeft := checkBasic(q[pair][0], needFound, 1)
		resepRight := checkBasic(q[pair][1], needFound, resepLeft.banyakResep)

		for i := 0; i < resepLeft.banyakResep; i++ {
			for j := 0; j < resepRight.banyakResep; j++ {
				res := Recipe{
					banyakNode: resepLeft.kumpulanResep[i].banyakNode + resepRight.kumpulanResep[j].banyakNode + 1,
					resep: &RecipeTreeNode{
						namaElemen: elements[start].Name,
						dibuatDari: [][2]*RecipeTreeNode{
							{resepLeft.kumpulanResep[i].resep, resepRight.kumpulanResep[j].resep},
						},
					},
				}
				// fmt.Println("Resep ditemukan untuk", elements[start].Name, "dari", elements[q[pair][0]].Name, "dan", elements[q[pair][1]].Name)
				// fmt.Println("Banyak node:", res.banyakNode)
				retResep.kumpulanResep = append(retResep.kumpulanResep, res)
				retResep.banyakResep+= 1
				notFound += mult
				if notFound >= needFound {
					visit[start] = false
					return retResep
				}
			}
		}
	}
	visit[start] = false
	return retResep
}


func bfs(start int, needFound int) RecipeList {
	queue := make([]int, 0)
	queue = append(queue, start)
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if elements[node].Name == "water" || elements[node].Name == "fire" || elements[node].Name == "earth" || elements[node].Name == "air" {
			return RecipeList{
				namaElemen: elements[node].Name,
				buatDari:   nil,
				banyakResep: 1,
			}
		}

		for _, pair := range elements[node].FromPair {
			left := pair[0]
			right := pair[1]

			if left != -1 && right != -1 {
				queue = append(queue, left, right)
			}
		}
	}
}
