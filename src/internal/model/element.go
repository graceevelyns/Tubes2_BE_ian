package model

// nanti sesuaikan aja yak @ian
type RecipeNode struct {
	// NamaElemen adalah nama dari elemen :)
	// contoh: "Steam", "Air"
	NamaElemen string `json:"namaElemen"`

	// IsBaseElement bernilai true jika elemen ini adalah salah satu dari 4 elemen dasar -> buat endpoint
	// (Air, Earth, Fire, Water) yang tidak dibuat dari elemen lain
	IsBaseElement bool `json:"isBaseElement"`

	// DibuatDari as slice dari pasangan (array dengan 2 elemen) pointer ke RecipeNode lain
	// tiap pasangan merepresentasikan dua bahan yang dikombinasikan untuk menghasilkan NamaElemen
	DibuatDari [][2]*RecipeNode `json:"dibuatDari,omitempty"`
}
