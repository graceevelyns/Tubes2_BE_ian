package algorithm



type RecipeTreeNodeChild struct {
	parent *RecipeTreeNode;
	leftChild *RecipeTreeNode;
	rightChild *RecipeTreeNode;

	leftChildID int;
	rightChildID int;
}

// tambah banyakResep ke parent"
func bfsHelper (a *RecipeTreeNode) {
	kiri := (a.parent.leftChild == a)
	newBanyakResep := 1
	for a.parent != nil {
		if a.parent.leftChild.banyakResep > 0 && a.parent.rightChild.banyakResep > 0 {
			a.parent.parent.banyakResep /= a.banyakResep
			a.banyakResep = newBanyakResep
		}else if a.parent.leftChild.banyakResep == 0  && a.parent.rightChild.banyakResep == 0 {
			a.banyakResep = newBanyakResep
			newBanyakResep = 0
		}else if (kiri && a.parent.leftChild.banyakResep == 0 ) || (!kiri && a.parent.rightChild.banyakResep == 0) {
			a.banyakResep = newBanyakResep
			if kiri {
				newBanyakResep = a.banyakResep * a.parent.rightChild.banyakResep
			}else{
				newBanyakResep = a.banyakResep * a.parent.leftChild.banyakResep
			}
		}else if (kiri && a.parent.rightChild.banyakResep == 0) || (!kiri && a.parent.leftChild.banyakResep == 0) {
			a.banyakResep = newBanyakResep
			newBanyakResep = 0
		}else{ // ERROR !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			newBanyakResep = 0
		}
		a = a.parent.parent
	}
	if newBanyakResep > 0 {
		if a.banyakResep == 0 {
			a.banyakResep = newBanyakResep
		}else{
			a.banyakResep *= newBanyakResep
		}
	}
}
		
func isBasicElement (a RecipeTreeNode) bool {
	if a.namaElemen == "water" || a.namaElemen == "fire" || a.namaElemen == "earth" || a.namaElemen == "air" {
		return true
	}
	return false
}

func isBasicElement_typeElement (a Element) bool {
	if a.Name == "water" || a.Name == "fire" || a.Name == "earth" || a.Name == "air" {
		return true
	}
	return false
}

func isBFSPath (a *RecipeTreeNodeChild) bool{
	return (a.leftChild.banyakResep > 0 || isBasicElement(*a.leftChild)) && (a.rightChild.banyakResep > 0 || isBasicElement(*a.rightChild))
}

func bfs(start int, needFound int, mult int) *RecipeTreeNode {
	var Tree = RecipeTreeNode{
		namaElemen: elements[start].Name,
		dibuatDari: make([]RecipeTreeNodeChild, 0),
		banyakResep: 0,
	}


	//initialize queue
	queue := make([]RecipeTreeNodeChild, 0)
	for i := 0 ; i < len(elements[start].FromPair); i++ {
		queue = append(queue, RecipeTreeNodeChild{
			parent: &Tree,
			leftChild: &RecipeTreeNode{
				namaElemen: elements[elements[start].FromPair[i][0]].Name,
				dibuatDari: make([]RecipeTreeNodeChild, 0),
				banyakResep: 0,
			},
			rightChild: &RecipeTreeNode{
				namaElemen: elements[elements[start].FromPair[i][1]].Name,
				dibuatDari: make([]RecipeTreeNodeChild, 0),
				banyakResep: 0,
			},
			leftChildID: elements[start].FromPair[i][0],
			rightChildID: elements[start].FromPair[i][1],
		})
	}


	for len(queue) > 0 {
		var current = queue[0]
		queue = queue[1:]
		
		current.parent.dibuatDari = append(current.parent.dibuatDari, current)

		if isBasicElement(*current.leftChild) {
			bfsHelper(current.leftChild)
		} else {
			for i := 0 ; i < len(elements[current.leftChildID].FromPair); i++ {
				queue = append(queue, RecipeTreeNodeChild{
					parent: current.leftChild,
					leftChild: &RecipeTreeNode{
						namaElemen: elements[elements[current.leftChildID].FromPair[i][0]].Name,
						dibuatDari: make([]RecipeTreeNodeChild, 0),
						banyakResep: 0,
					},
					rightChild: &RecipeTreeNode{
						namaElemen: elements[elements[current.leftChildID].FromPair[i][1]].Name,
						dibuatDari: make([]RecipeTreeNodeChild, 0),
						banyakResep: 0,
					},
					leftChildID: elements[current.leftChildID].FromPair[i][0],
					rightChildID: elements[current.leftChildID].FromPair[i][1],
				})
			}
		}

		if isBasicElement(*current.rightChild) {
			bfsHelper(current.rightChild)
		} else {
			for i := 0 ; i < len(elements[current.rightChildID].FromPair); i++ {
				queue = append(queue, RecipeTreeNodeChild{
					parent: current.rightChild,
					leftChild: &RecipeTreeNode{
						namaElemen: elements[elements[current.rightChildID].FromPair[i][0]].Name,
						dibuatDari: make([]RecipeTreeNodeChild, 0),
						banyakResep: 0,
					},
					rightChild: &RecipeTreeNode{
						namaElemen: elements[elements[current.rightChildID].FromPair[i][1]].Name,
						dibuatDari: make([]RecipeTreeNodeChild, 0),
						banyakResep: 0,
					},
					leftChildID: elements[current.rightChildID].FromPair[i][0],
					rightChildID: elements[current.rightChildID].FromPair[i][1],
				})
			}
		}
	}

	var res = &Tree
	return res
}