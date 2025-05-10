package algorithm

type Element struct {
	id       int
	tier     int
	Name     string
	FromPair [][2]int
	CanMake  []int
}

var elements = make([]Element, 0)

type RecipeTreeNode struct {
	namaElemen  string
	dibuatDari  []RecipeTreeNodeChild
	banyakResep int
	parent      *RecipeTreeNodeChild
}

func dfs(start int, needFound int, mult int) *RecipeTreeNode {
	if isBasicElement_typeElement(elements[start]) {
		return &RecipeTreeNode{
			namaElemen:  elements[start].Name,
			dibuatDari:  nil,
			banyakResep: 1,
		}
	}
	var Tree = RecipeTreeNode{
		namaElemen:  elements[start].Name,
		dibuatDari:  make([]RecipeTreeNodeChild, 0),
		banyakResep: 0,
	}

	for i := 0; i < len(elements[start].FromPair); i++ {
		if !isValid(start, elements[start].FromPair[i][0], elements[start].FromPair[i][1]) {
			continue
		}
		var leftNode = dfs(elements[start].FromPair[i][0], needFound-Tree.banyakResep, mult)
		var rightNode = dfs(elements[start].FromPair[i][1], needFound-Tree.banyakResep, leftNode.banyakResep)
		Tree.dibuatDari = append(Tree.dibuatDari, RecipeTreeNodeChild{parent: &Tree, leftChild: leftNode, rightChild: rightNode, leftChildID: elements[start].FromPair[i][0], rightChildID: elements[start].FromPair[i][1]})
		Tree.banyakResep += leftNode.banyakResep * rightNode.banyakResep
		if Tree.banyakResep >= needFound {
			break
		}
	}
	var result = &Tree
	return result
}
