package main

type Element struct {
	id       int
	tier     int
	Name     string
	FromPair [][2]int
	CanMake  []int
}


var elements = make([]Element , 0)

type RecipeTreeNode struct{
    namaElemen string;
    dibuatDari [][2]*RecipeTreeNode;
	banyakResep int;
}


func dfs(start int, needFound int, mult int) *RecipeTreeNode {
	if(elements[start].Name == "water" || elements[start].Name == "fire" || elements[start].Name == "earth" || elements[start].Name == "air") {
		return &RecipeTreeNode{
			namaElemen: elements[start].Name,
			dibuatDari: nil,
			banyakResep: 1,
		}
	} 
	var Tree = RecipeTreeNode{
		namaElemen: elements[start].Name,
		dibuatDari: make([][2]*RecipeTreeNode, 0),
		banyakResep: 0,
	}
	
	for i := 0 ; i < len(elements[start].FromPair); i++ {
		var leftNode = dfs(elements[start].FromPair[i][0], needFound - Tree.banyakResep, mult)
		var rightNode = dfs(elements[start].FromPair[i][1], needFound - Tree.banyakResep, leftNode.banyakResep)
		Tree.dibuatDari = append(Tree.dibuatDari, [2]*RecipeTreeNode{leftNode, rightNode})
		Tree.banyakResep += leftNode.banyakResep * rightNode.banyakResep
		if Tree.banyakResep >= needFound {
			break
		}
	}
	var result = &Tree
	return result
}
