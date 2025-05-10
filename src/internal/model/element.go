package model

type RecipeTreeNodeTier struct {
	NamaElemen string
	DibuatDari [][2]*RecipeTreeNodeTier
}

// type Recipe struct {
// 	// ditemukan bool
// 	// stepByStep []int
// 	BanyakNode int
// 	Resep      *RecipeTreeNode
// }

// type RecipeList struct {
// 	NamaElemen    string
// 	KumpulanResep []Recipe
// 	BanyakResep   int
// }

type Element struct {
	Id       int
	Tier     int
	Name     string
	FromPair [][2]int
	CanMake  []int
}

type RecipeOutputNode struct {
	NamaElemen    string                `json:"namaElemen"`
	IsBaseElement bool                  `json:"isBaseElement"`
	ID            int                   `json:"id,omitempty"`
	Tier          int                   `json:"tier,omitempty"`
	DibuatDari    [][]*RecipeOutputNode `json:"dibuatDari,omitempty"`
}
