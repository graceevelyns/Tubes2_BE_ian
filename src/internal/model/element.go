package model

type RecipeTreeNode struct {
	NamaElemen string
	DibuatDari [][2]*RecipeTreeNode
}

type Recipe struct {
	// ditemukan bool
	// stepByStep []int
	BanyakNode int
	Resep      *RecipeTreeNode
}

type RecipeList struct {
	NamaElemen    string
	KumpulanResep []Recipe
	BanyakResep   int
}

type Element struct {
	Id       int
	Tier     int
	Name     string
	FromPair [][2]int
	CanMake  []int
}
