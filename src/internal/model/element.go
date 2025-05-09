package model

type RecipeTreeNode struct {
	NamaElemen string
	DibuatDari [][2]*RecipeTreeNode
}

type Recipe struct {
	// ditemukan bool
	// stepByStep []int
	banyakNode int
	resep      *RecipeTreeNode
}

type RecipeList struct {
	namaElemen    string
	kumpulanResep []Recipe
	banyakResep   int
}

type Element struct {
	id       int
	tier     int
	Name     string
	FromPair [][2]int
	CanMake  []int
}
