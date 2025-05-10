package algorithm

import (
	"log"

	"github.com/graceevelyns/Tubes2_BE_ian/src/internal/scraper"
)

var elements []*scraper.Element

func InitializeAlgorithmElements(loadedElements []*scraper.Element) {
	log.Println("[INIT_ALGO_DEBUG] Memulai InitializeAlgorithmElements...")
	if loadedElements == nil {
		log.Println("[INIT_ALGO_DEBUG] Peringatan: loadedElements adalah nil. Menginisialisasi elements sebagai slice kosong.")
		elements = []*scraper.Element{}
		log.Printf("[INIT_ALGO_DEBUG] Selesai InitializeAlgorithmElements. Jumlah elements: %d", len(elements))
		return
	}
	elements = loadedElements
	log.Printf("[INIT_ALGO_DEBUG] Elements untuk package algorithm berhasil diinisialisasi. Jumlah aktual: %d. Pointer: %p", len(elements), elements)

	// Log yang sudah ada dari Anda, mungkin perlu disesuaikan jika ID 16 bukan indeks yang valid
	// Kita akan tambahkan pengecekan batas untuk log ini agar tidak panic
	if len(elements) > 16 { // Jika panjangnya 17 atau lebih, indeks 16 valid
		log.Printf("[INIT_ALGO_DEBUG] Info: Elemen di INDEKS 16: Nama='%s', ID=%d", elements[16].Name, elements[16].ID)
	} else {
		log.Printf("[INIT_ALGO_DEBUG] Info: Jumlah elemen %d, tidak cukup untuk mengakses INDEKS 16 secara langsung.", len(elements))
	}

	// Log tambahan untuk melihat elemen di indeks 15 (yang mungkin adalah elemen dengan ID 16)
	if len(elements) > 15 { // Jika panjangnya 16 atau lebih, indeks 15 valid
		log.Printf("[INIT_ALGO_DEBUG] Info: Elemen di INDEKS 15: Nama='%s', ID=%d", elements[15].Name, elements[15].ID)
	} else {
		log.Printf("[INIT_ALGO_DEBUG] Info: Jumlah elemen %d, tidak cukup untuk mengakses INDEKS 15.", len(elements))
	}

	if len(elements) == 0 {
		log.Println("[INIT_ALGO_DEBUG] Peringatan: Elements kosong setelah inisialisasi.")
	}
	log.Println("[INIT_ALGO_DEBUG] Selesai InitializeAlgorithmElements.")
}

type RecipeTreeNode struct {
	NamaElemen  string
	DibuatDari  []RecipeTreeNodeChild
	BanyakResep int
	Parent      *RecipeTreeNodeChild
}

func Dfs(start int, needFound int, mult int) *RecipeTreeNode {
	log.Printf("[DFS_DEBUG] Memulai DFS untuk start (ID/Indeks?): %d, needFound: %d, mult: %d. Panjang elements: %d", start, needFound, mult, len(elements))

	// Pengecekan batas KRUSIAL sebelum mengakses elements[start]
	if start < 0 || start >= len(elements) {
		log.Printf("[DFS_DEBUG] ERROR: 'start' (%d) di luar jangkauan slice 'elements' (panjang: %d). Mengembalikan nil.", start, len(elements))
		return nil // Jika 'start' adalah indeks, ini penting. Jika 'start' adalah ID 1-based, maka akses seharusnya elements[start-1]
	}
	// Jika 'start' adalah ID 1-based dan Anda belum mengubah logika, log ini akan menunjukkan elemen yang salah
	log.Printf("[DFS_DEBUG] Mengakses elements[%d]: Nama='%s', ID=%d, Tier=%d, Jumlah FromPair=%d",
		start, elements[start].Name, elements[start].ID, elements[start].Tier, len(elements[start].FromPair))

	if IsBasicElement_typeElement(*elements[start]) {
		log.Printf("[DFS_DEBUG] Elemen '%s' (ID: %d) adalah elemen dasar. Mengembalikan node dasar.", elements[start].Name, elements[start].ID)
		return &RecipeTreeNode{
			NamaElemen:  elements[start].Name,
			DibuatDari:  nil,
			BanyakResep: 1,
		}
	}

	log.Printf("[DFS_DEBUG] Elemen '%s' (ID: %d) BUKAN elemen dasar. Membuat Tree.", elements[start].Name, elements[start].ID)
	var Tree = RecipeTreeNode{
		NamaElemen:  elements[start].Name,
		DibuatDari:  make([]RecipeTreeNodeChild, 0),
		BanyakResep: 0,
	}
	log.Printf("[DFS_DEBUG] Tree untuk '%s' dibuat. BanyakResep awal: %d", Tree.NamaElemen, Tree.BanyakResep)

	if len(elements[start].FromPair) == 0 {
		log.Printf("[DFS_DEBUG] Elemen '%s' (ID: %d) tidak memiliki FromPair. Akan mengembalikan Tree dengan BanyakResep: %d.", Tree.NamaElemen, elements[start].ID, Tree.BanyakResep)
	}

	for i := 0; i < len(elements[start].FromPair); i++ {
		parentElementName := elements[start].Name // Untuk logging
		pair := elements[start].FromPair[i]
		leftChildID := pair[0]
		rightChildID := pair[1]
		log.Printf("[DFS_DEBUG] [%s] Iterasi FromPair ke-%d: LeftChildID=%d, RightChildID=%d", parentElementName, i, leftChildID, rightChildID)

		if !IsValid(start, leftChildID, rightChildID) {
			log.Printf("[DFS_DEBUG] [%s] Pair (LeftID: %d, RightID: %d) TIDAK VALID menurut IsValid. Melanjutkan ke pair berikutnya.", parentElementName, leftChildID, rightChildID)
			continue
		}
		log.Printf("[DFS_DEBUG] [%s] Pair (LeftID: %d, RightID: %d) VALID. Memanggil DFS rekursif.", parentElementName, leftChildID, rightChildID)

		// Logika needFound-Tree.BanyakResep bisa jadi kompleks, mari kita log nilainya
		needFoundForLeft := needFound - Tree.BanyakResep
		log.Printf("[DFS_DEBUG] [%s] Memanggil DFS untuk LeftChildID: %d, needFoundForLeft: %d", parentElementName, leftChildID, needFoundForLeft)
		var leftNode = Dfs(leftChildID, needFoundForLeft, mult) // `mult` mungkin perlu disesuaikan untuk anak kanan

		if leftNode == nil {
			log.Printf("[DFS_DEBUG] [%s] Panggilan DFS untuk LeftChildID: %d mengembalikan nil. Mungkin tidak bisa membentuk resep ini.", parentElementName, leftChildID)
			// Jika salah satu anak nil, mungkin kombinasi ini tidak valid atau tidak bisa diselesaikan
			// Tergantung logika bisnis, Anda bisa `continue` atau menangani kasus ini
			// Untuk saat ini, biarkan logika asli berjalan, tapi waspadai dampaknya
		} else {
			log.Printf("[DFS_DEBUG] [%s] Panggilan DFS untuk LeftChildID: %d selesai. leftNode.BanyakResep: %d", parentElementName, leftChildID, leftNode.BanyakResep)
		}

		// mult untuk anak kanan menggunakan leftNode.BanyakResep
		needFoundForRight := needFound - Tree.BanyakResep // Ini mungkin perlu dihitung ulang jika Tree.BanyakResep berubah oleh leftNode
		multForRight := 1                                 // Default, jika leftNode nil atau BanyakResep 0
		if leftNode != nil {
			multForRight = leftNode.BanyakResep
		}
		log.Printf("[DFS_DEBUG] [%s] Memanggil DFS untuk RightChildID: %d, needFoundForRight: %d, multForRight: %d", parentElementName, rightChildID, needFoundForRight, multForRight)
		var rightNode = Dfs(rightChildID, needFoundForRight, multForRight)

		if rightNode == nil {
			log.Printf("[DFS_DEBUG] [%s] Panggilan DFS untuk RightChildID: %d mengembalikan nil. Mungkin tidak bisa membentuk resep ini.", parentElementName, rightChildID)
			// Sama seperti leftNode, tangani jika perlu
		} else {
			log.Printf("[DFS_DEBUG] [%s] Panggilan DFS untuk RightChildID: %d selesai. rightNode.BanyakResep: %d", parentElementName, rightChildID, rightNode.BanyakResep)
		}

		// Hanya tambahkan jika kedua node valid, untuk menghindari panic jika salah satunya nil
		if leftNode != nil && rightNode != nil {
			Tree.DibuatDari = append(Tree.DibuatDari, RecipeTreeNodeChild{Parent: &Tree, LeftChild: leftNode, RightChild: rightNode, LeftChildID: leftChildID, RightChildID: rightChildID})
			newRecipesFromPair := leftNode.BanyakResep * rightNode.BanyakResep
			Tree.BanyakResep += newRecipesFromPair
			log.Printf("[DFS_DEBUG] [%s] Pair (LeftID: %d, RightID: %d) ditambahkan ke DibuatDari. Resep baru dari pair: %d. Total Tree.BanyakResep: %d",
				parentElementName, leftChildID, rightChildID, newRecipesFromPair, Tree.BanyakResep)
		} else {
			log.Printf("[DFS_DEBUG] [%s] Tidak menambahkan pair (LeftID: %d, RightID: %d) karena satu atau kedua node anak adalah nil.", parentElementName, leftChildID, rightChildID)
		}

		if Tree.BanyakResep >= needFound {
			log.Printf("[DFS_DEBUG] [%s] Tree.BanyakResep (%d) >= needFound (%d). Menghentikan loop FromPair.", parentElementName, Tree.BanyakResep, needFound)
			break
		}
	}
	var result = &Tree
	log.Printf("[DFS_DEBUG] Selesai DFS untuk '%s' (ID dari start: %d). Mengembalikan Tree.BanyakResep: %d. Jumlah DibuatDari: %d", result.NamaElemen, start, result.BanyakResep, len(result.DibuatDari))
	return result
}
