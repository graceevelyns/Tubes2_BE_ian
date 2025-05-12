package algorithm

import (
	"fmt" // Tambahkan untuk formatting log
	// "log"

	"github.com/graceevelyns/Tubes2_BE_ian/src/cmd/internal/scraper"
)

// Variabel global 'elements' dari file lain di package ini
// var elements []*scraper.Element
// func InitializeAlgorithmElements(loadedElements []*scraper.Element) { ... }

type RecipeTreeNodeChild struct {
	Parent       *RecipeTreeNode
	LeftChild    *RecipeTreeNode
	RightChild   *RecipeTreeNode
	LeftChildID  int
	RightChildID int
}

// tambah BanyakResep ke Parent"

func BfsHelper(a *RecipeTreeNode) {
	// log.Printf("[BFS_HELPER_DEBUG] === Memulai BfsHelper untuk node: '%s', BanyakResepAwalNodeA: %d ===", a.NamaElemen, a.BanyakResep)

	if a.Parent == nil {
		// log.Printf("[BFS_HELPER_DEBUG] FATAL: a.Parent adalah nil untuk node '%s'. Ini AKAN MENYEBABKAN PANIC jika kode berlanjut!", a.NamaElemen)
		// Seharusnya ada return di sini atau panic akan terjadi di baris berikutnya.
		// Namun, karena Anda minta tidak ubah kode, kita biarkan.
		// Panic akan terjadi di baris: kiri := (a.Parent.LeftChild == a)
	} else {
		// Log jika a.Parent TIDAK nil
		// var parentParentNodeName string = "N/A (a.Parent.Parent is nil)"
		if a.Parent.Parent != nil { // Pastikan kakek tidak nil
			// parentParentNodeName = a.Parent.Parent.NamaElemen // Ini adalah *RecipeTreeNode, jadi .NamaElemen ada
			// Untuk log yang lebih dalam, kita periksa lagi:
			if a.Parent.Parent.Parent != nil { // Parent dari kakek adalah *RecipeTreeNodeChild
				// Kita tidak bisa langsung .NamaElemen di sini.
				// Jika ingin nama dari parent *RecipeTreeNodeChild* ini, itu adalah a.Parent.Parent.Parent.Parent.NamaElemen
				// log.Printf("[BFS_HELPER_DEBUG] Node kakek ('%s') memiliki parent (RecipeTreeNodeChild). LeftChild kakek-parent: %s, RightChild kakek-parent: %s",parentParentNodeName,a.Parent.Parent.Parent.LeftChild.NamaElemen,  // Asumsi tidak nila.Parent.Parent.Parent.RightChild.NamaElemen) // Asumsi tidak nil
			} else {
				// log.Printf("[BFS_HELPER_DEBUG] Node kakek ('%s') TIDAK memiliki parent (a.Parent.Parent.Parent adalah nil).", parentParentNodeName)
			}
		} else {
			// log.Printf("[BFS_HELPER_DEBUG] a.Parent ada, tapi a.Parent.Parent (kakek) adalah nil.")
		}

		if a.Parent.LeftChild != nil {
			// log.Printf("[BFS_HELPER_DEBUG] Parent-RTC dari '%s' memiliki LeftChild '%s' (BanyakResep: %d)", a.NamaElemen, a.Parent.LeftChild.NamaElemen, a.Parent.LeftChild.BanyakResep)
		} else {
			// log.Printf("[BFS_HELPER_DEBUG] Parent-RTC dari '%s' memiliki LeftChild nil", a.NamaElemen)
		}
		if a.Parent.RightChild != nil {
			// log.Printf("[BFS_HELPER_DEBUG] Parent-RTC dari '%s' memiliki RightChild '%s' (BanyakResep: %d)", a.NamaElemen, a.Parent.RightChild.NamaElemen, a.Parent.RightChild.BanyakResep)
		} else {
			// log.Printf("[BFS_HELPER_DEBUG] Parent-RTC dari '%s' memiliki RightChild nil", a.NamaElemen)
		}
	}

	// Jika a.Parent adalah nil, baris berikut akan panic.
	// kiri := (a.Parent.LeftChild == a)
	// var parentOfParentNodeName string = "N/A (a.Parent.Parent is nil)"
	// if a.Parent.Parent != nil {
	// 	parentOfParentNodeName = a.Parent.Parent.NamaElemen
	// }
	// log.Printf("[BFS_HELPER_DEBUG] Node '%s' adalah anak kiri (kiri=%t). Parent-RTC-nya adalah anak dari node: '%s'", a.NamaElemen, kiri, parentOfParentNodeName)

	newBanyakResep := 1
	// log.Printf("[BFS_HELPER_DEBUG] newBanyakResep diinisialisasi menjadi %d", newBanyakResep)

	// originalNodeA := a

	iterCount := 0
	for a.Parent != nil { // Jika a.Parent awalnya nil, loop ini tidak akan berjalan.
		iterCount++
		// currentParentNodeName := "N/A"
		// currentGrandParentNodeName := "N/A"
		// currentGreatGrandParentNodeName := "N/A" // Untuk log yang lebih dalam

		if a.Parent.Parent != nil {
			// currentParentNodeName = a.Parent.Parent.NamaElemen
			if a.Parent.Parent.Parent != nil && a.Parent.Parent.Parent.Parent != nil { // Parent dari RecipeTreeNodeChild adalah RecipeTreeNode
				// currentGreatGrandParentNodeName = a.Parent.Parent.Parent.Parent.NamaElemen // Ini adalah Parent dari RTC buyut
			}
		}
		if a.Parent.Parent != nil && a.Parent.Parent.Parent != nil {
			// a.Parent.Parent.Parent adalah *RecipeTreeNodeChild. Parent-nya adalah *RecipeTreeNode.
			if a.Parent.Parent.Parent.Parent != nil {
				// currentGrandParentNodeName = a.Parent.Parent.Parent.Parent.NamaElemen
			}
		}

		// log.Printf("[BFS_HELPER_DEBUG] Iterasi ke-%d loop. Node 'a': '%s'. Parent-RTC-nya ('%p') adalah anak dari '%s'. Kakek-RTC-nya ('%p') adalah anak dari '%s'",iterCount, a.NamaElemen, a.Parent, currentParentNodeName,a.Parent.Parent.Parent,     // Ini adalah RecipeTreeNodeChild "buyut"currentGrandParentNodeName) // Ini nama elemen dari parent si "buyut RTC"

		// log.Printf("[BFS_HELPER_DEBUG]   Kondisi: a.Parent.LeftChild.BanyakResep=%d, a.Parent.RightChild.BanyakResep=%d", a.Parent.LeftChild.BanyakResep, a.Parent.RightChild.BanyakResep)

		// ... (sisa kondisi if/else if dengan logging yang sudah ada, tapi pastikan tidak ada akses .NamaElemen yang salah) ...
		// Contoh perbaikan untuk salah satu log di dalam if/else:
		// if a.Parent.Parent.BanyakResep (Ini benar karena a.Parent.Parent adalah RecipeTreeNode)

		// Bagian yang paling berisiko adalah saat mengupdate 'a'
		nextNodeA := a.Parent.Parent // Ini adalah RecipeTreeNode "kakek"
		// var nextNodeAName string = "N/A (nextNodeA will be nil)"
		if nextNodeA != nil {
			// nextNodeAName = nextNodeA.NamaElemen
		}
		// log.Printf("[BFS_HELPER_DEBUG]   Akan mengupdate 'a' dari '%s' menjadi a.Parent.Parent (yaitu '%s')", a.NamaElemen, nextNodeAName)

		if a.Parent.Parent == nil {
			// log.Printf("[BFS_HELPER_DEBUG]   PERINGATAN: a.Parent.Parent adalah nil! Loop akan berhenti setelah 'a' diupdate menjadi nil.")
		}
		a = nextNodeA // a sekarang menjadi "kakek" atau nil
	}
	// ... (sisa BfsHelper dengan logging yang ada) ...
	// log.Printf("[BFS_HELPER_DEBUG] Keluar dari loop 'for a.Parent != nil'. Node 'a' sekarang: '%s'. newBanyakResep: %d", a.NamaElemen, newBanyakResep)

	if newBanyakResep > 0 {
		// log.Printf("[BFS_HELPER_DEBUG] newBanyakResep (%d) > 0. Node 'a' saat ini: '%s', a.BanyakResep: %d", newBanyakResep, a.NamaElemen, a.BanyakResep)
		if a.BanyakResep == 0 {
			a.BanyakResep = newBanyakResep
			// log.Printf("[BFS_HELPER_DEBUG]   a.BanyakResep diupdate menjadi newBanyakResep: %d", a.BanyakResep)
		} else {
			a.BanyakResep *= newBanyakResep
			// log.Printf("[BFS_HELPER_DEBUG]   a.BanyakResep dikalikan dengan newBanyakResep. Hasil: %d", a.BanyakResep)
		}
	} else {
		// log.Printf("[BFS_HELPER_DEBUG] newBanyakResep (%d) tidak lebih besar dari 0. Tidak ada update terakhir untuk a.BanyakResep.", newBanyakResep)
	}
	// log.Printf("[BFS_HELPER_DEBUG] === Selesai BfsHelper untuk node AWAL: '%s'. BanyakResepAkhirNodeAWAL (originalNodeA): %d. BanyakResepNodeRoot (jika 'a' mencapai root): %d ===", originalNodeA.NamaElemen, originalNodeA.BanyakResep, a.BanyakResep)

}

func BfsHelperORI(a *RecipeTreeNode) {
	kiri := (a.Parent.LeftChild == a)
	newBanyakResep := 1
	for a.Parent != nil {
		if a.Parent.LeftChild.BanyakResep > 0 && a.Parent.RightChild.BanyakResep > 0 {
			a.Parent.Parent.BanyakResep -= a.Parent.LeftChild.BanyakResep * a.Parent.RightChild.BanyakResep
			a.BanyakResep = newBanyakResep
		} else if a.Parent.LeftChild.BanyakResep == 0 && a.Parent.RightChild.BanyakResep == 0 {
			a.BanyakResep = newBanyakResep
			newBanyakResep = 0
		} else if (kiri && a.Parent.LeftChild.BanyakResep == 0) || (!kiri && a.Parent.RightChild.BanyakResep == 0) {
			a.BanyakResep = newBanyakResep
			if kiri {
				newBanyakResep = a.BanyakResep * a.Parent.RightChild.BanyakResep
			} else {
				newBanyakResep = a.BanyakResep * a.Parent.LeftChild.BanyakResep
			}
		} else if (kiri && a.Parent.RightChild.BanyakResep == 0) || (!kiri && a.Parent.LeftChild.BanyakResep == 0) {
			a.BanyakResep = newBanyakResep
			newBanyakResep = 0
		} else { // ERROR !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			newBanyakResep = 0
		}
		a = a.Parent.Parent
	}
	a.BanyakResep += newBanyakResep

}

func IsBasicElement(a RecipeTreeNode) bool {
	// Normalisasi string nama elemen dasar agar konsisten (misal semua TitleCase)
	// Logika asli Anda menggunakan TitleCase, jadi kita asumsikan NamaElemen juga TitleCase.
	isBase := a.NamaElemen == "Water" || a.NamaElemen == "Fire" || a.NamaElemen == "Earth" || a.NamaElemen == "Air"
	// log.Printf("[IS_BASIC_DEBUG] IsBasicElement untuk '%s': %t", a.NamaElemen, isBase)
	return isBase
}

func IsBasicElement_typeElement(a scraper.Element) bool {
	// Sama seperti di atas, perhatikan kapitalisasi
	isBase := a.Name == "Water" || a.Name == "Fire" || a.Name == "Earth" || a.Name == "Air"
	// log.Printf("[IS_BASIC_TYPE_DEBUG] IsBasicElement_typeElement untuk '%s' (ID %d): %t", a.Name, a.ID, isBase)
	return isBase
}

func IsBFSPath(a *RecipeTreeNodeChild) bool {
	// Tambahkan pengecekan nil sebelum dereference
	if a == nil || a.LeftChild == nil || a.RightChild == nil {
		// log.Printf("[IS_BFS_PATH_DEBUG] ERROR: a atau anak-anaknya adalah nil. Mengembalikan false.")
		return false
	}
	result := (a.LeftChild.BanyakResep > 0 || IsBasicElement(*a.LeftChild)) && (a.RightChild.BanyakResep > 0 || IsBasicElement(*a.RightChild))
	// log.Printf("[IS_BFS_PATH_DEBUG] Memeriksa path untuk Parent '%s', Left '%s'(Resep:%d, Basic:%t), Right '%s'(Resep:%d, Basic:%t). Hasil: %t",a.Parent.NamaElemen,a.LeftChild.NamaElemen, a.LeftChild.BanyakResep, IsBasicElement(*a.LeftChild),a.RightChild.NamaElemen, a.RightChild.BanyakResep, IsBasicElement(*a.RightChild),result)
	return result
}

func IsValid(idParent int, idLeftChild int, idRightChild int) bool {
	// log.Printf("[IS_VALID_DEBUG] Memeriksa IsValid: idParent=%d, idLeftChild=%d, idRightChild=%d. Panjang elements: %d", idParent, idLeftChild, idRightChild, len(elements))

	// PENTING: Kode ini mengasumsikan ID adalah INDEKS. Jika ID adalah 1-based, ini akan salah.
	// Log akan membantu melihat elemen apa yang diakses.

	// Pengecekan batas untuk semua ID yang digunakan sebagai indeks
	if idParent < 0 || idParent >= len(elements) {
		// log.Printf("[IS_VALID_DEBUG] ERROR: idParent (%d) di luar jangkauan 'elements' (panjang %d). Mengembalikan false.", idParent, len(elements))
		return false
	}
	elParent := elements[idParent]

	if idLeftChild < 0 || idLeftChild >= len(elements) {
		// log.Printf("[IS_VALID_DEBUG] ERROR: idLeftChild (%d) di luar jangkauan 'elements' (panjang %d). Mengembalikan false.", idLeftChild, len(elements))
		return false
	}
	elLeft := elements[idLeftChild]

	if idRightChild < 0 || idRightChild >= len(elements) {
		// log.Printf("[IS_VALID_DEBUG] ERROR: idRightChild (%d) di luar jangkauan 'elements' (panjang %d). Mengembalikan false.", idRightChild, len(elements))
		return false
	}
	elRight := elements[idRightChild]

	// log.Printf("[IS_VALID_DEBUG] Data untuk IsValid: Parent (diakses pada indeks %d): '%s' (ID %d, Tier %d), Left (diakses pada indeks %d): '%s' (ID %d, Tier %d), Right (diakses pada indeks %d): '%s' (ID %d, Tier %d)",idParent, elParent.Name, elParent.ID, elParent.Tier,idLeftChild, elLeft.Name, elLeft.ID, elLeft.Tier,idRightChild, elRight.Name, elRight.ID, elRight.Tier)

	result := elParent.Tier > elLeft.Tier && elParent.Tier > elRight.Tier
	// log.Printf("[IS_VALID_DEBUG] Hasil perbandingan Tier: (%d > %d && %d > %d) -> %t",elParent.Tier, elLeft.Tier, elParent.Tier, elRight.Tier, result)
	return result
}

func Bfs(start int, needFound int, mult int) *RecipeTreeNode {
	// log.Printf("[BFS_DEBUG] === Memulai BFS untuk start (digunakan sbg INDEKS): %d, needFound: %d, mult: %d. Panjang elements: %d ===", start, needFound, mult, len(elements))

	// PENTING: Logika asli menggunakan 'start' langsung sebagai indeks.
	// Jika 'start' adalah ID 1-based dari handler, ini akan salah.
	// Log akan menunjukkan elemen apa yang diakses.

	if start < 0 || start >= len(elements) {
		// log.Printf("[BFS_DEBUG] ERROR: 'start' (%d) sebagai INDEKS di luar jangkauan 'elements' (panjang %d). Mengembalikan Tree kosong.", start, len(elements))
		return &RecipeTreeNode{
			NamaElemen:  fmt.Sprintf("Error_Indeks_Root_%d_Invalid", start),
			DibuatDari:  make([]RecipeTreeNodeChild, 0),
			BanyakResep: 0,
		}
	}
	// log.Printf("[BFS_DEBUG] Root Tree akan dibuat dari elements[%d]: Nama='%s', ID=%d, Tier=%d", start, elements[start].Name, elements[start].ID, elements[start].Tier)

	var Tree = RecipeTreeNode{
		NamaElemen:  elements[start].Name,
		DibuatDari:  make([]RecipeTreeNodeChild, 0),
		BanyakResep: 0,
		// Parent untuk root Tree adalah nil, ini OK.
	}
	// log.Printf("[BFS_DEBUG] Tree root '%s' dibuat. BanyakResep awal: %d", Tree.NamaElemen, Tree.BanyakResep)

	queue := make([]RecipeTreeNodeChild, 0)

	// --- Inisialisasi Queue Awal ---
	if len(elements[start].FromPair) == 0 {
		// log.Printf("[BFS_DEBUG] [InitQueue] Elemen root '%s' (dari INDEKS %d) tidak memiliki FromPair.", Tree.NamaElemen, start)
	}
	for i := 0; i < len(elements[start].FromPair); i++ {
		pair := elements[start].FromPair[i]
		// ID dari FromPair adalah ID elemen, BUKAN indeks. Ini perlu diperhatikan saat mengakses 'elements' lagi.
		leftChildID_fromPair := pair[0]
		rightChildID_fromPair := pair[1]
		// log.Printf("[BFS_DEBUG] [InitQueue] Memeriksa pair untuk root '%s': LeftChildID=%d, RightChildID=%d", Tree.NamaElemen, leftChildID_fromPair, rightChildID_fromPair)

		// IsValid menerima ID yang akan digunakannya sebagai INDEKS jika mengikuti logika asli IsValid.
		if IsValid(start, leftChildID_fromPair, rightChildID_fromPair) {
			// log.Printf("[BFS_DEBUG] [InitQueue] Pair (LeftID: %d, RightID: %d) VALID.", leftChildID_fromPair, rightChildID_fromPair)

			// Mengakses elemen anak menggunakan ID dari pair SEBAGAI INDEKS (sesuai logika asli)
			// Ini berpotensi besar menjadi sumber masalah jika ID != Indeks
			if leftChildID_fromPair < 0 || leftChildID_fromPair >= len(elements) || rightChildID_fromPair < 0 || rightChildID_fromPair >= len(elements) {
				// log.Printf("[BFS_DEBUG] [InitQueue] ERROR: ID anak dari FromPair (L:%d, R:%d) di luar jangkauan jika digunakan sebagai INDEKS. Melewati pair.", leftChildID_fromPair, rightChildID_fromPair)
				continue
			}
			// log.Printf("[BFS_DEBUG] [InitQueue] Membuat LeftChild dari elements[%d]: Nama='%s'", leftChildID_fromPair, elements[leftChildID_fromPair].Name)
			// log.Printf("[BFS_DEBUG] [InitQueue] Membuat RightChild dari elements[%d]: Nama='%s'", rightChildID_fromPair, elements[rightChildID_fromPair].Name)

			// 1. Buat node RecipeTreeNode (anak) - Parent belum diset
			tempLeftChildNode := &RecipeTreeNode{
				NamaElemen:  elements[leftChildID_fromPair].Name, // Menggunakan ID sebagai INDEKS
				DibuatDari:  make([]RecipeTreeNodeChild, 0),
				BanyakResep: 0,
			}
			tempRightChildNode := &RecipeTreeNode{
				NamaElemen:  elements[rightChildID_fromPair].Name, // Menggunakan ID sebagai INDEKS
				DibuatDari:  make([]RecipeTreeNodeChild, 0),
				BanyakResep: 0,
			}

			// 2. Buat RecipeTreeNodeChild
			var childRelation RecipeTreeNodeChild
			childRelation = RecipeTreeNodeChild{
				Parent:       &Tree,
				LeftChild:    tempLeftChildNode,
				RightChild:   tempRightChildNode,
				LeftChildID:  leftChildID_fromPair,  // Simpan ID asli
				RightChildID: rightChildID_fromPair, // Simpan ID asli
			}

			// 3. ATUR Parent dari node anak untuk menunjuk ke childRelation
			// Ini adalah perbaikan yang disarankan sebelumnya untuk mengatasi panic di BfsHelper
			tempLeftChildNode.Parent = &childRelation
			tempRightChildNode.Parent = &childRelation
			// log.Printf("[BFS_DEBUG] [InitQueue] Pointer Parent untuk '%s' dan '%s' telah diatur ke RecipeTreeNodeChild.", tempLeftChildNode.NamaElemen, tempRightChildNode.NamaElemen)

			queue = append(queue, childRelation)
			// log.Printf("[BFS_DEBUG] [InitQueue] Menambahkan ke queue: Parent='%s', Left='%s'(ID:%d), Right='%s'(ID:%d). Queue size: %d",childRelation.Parent.NamaElemen, tempLeftChildNode.NamaElemen, leftChildID_fromPair, tempRightChildNode.NamaElemen, rightChildID_fromPair, len(queue))
		} else {
			// log.Printf("[BFS_DEBUG] [InitQueue] Pair (LeftID: %d, RightID: %d) TIDAK VALID.", leftChildID_fromPair, rightChildID_fromPair)
		}
	}
	// log.Printf("[BFS_DEBUG] Inisialisasi queue selesai. Ukuran queue akhir: %d", len(queue))

	// --- Loop Utama BFS ---
	processedInQueue := 0
	for len(queue) > 0 {
		processedInQueue++
		current := queue[0]
		queue = queue[1:]
		// log.Printf("[BFS_DEBUG] [Loop %d] Mengambil dari queue: Parent='%s', LeftChild='%s'(ID:%d), RightChild='%s'(ID:%d). Sisa queue: %d",processedInQueue, current.Parent.NamaElemen, current.LeftChild.NamaElemen, current.LeftChildID, current.RightChild.NamaElemen, current.RightChildID, len(queue))

		if current.Parent == nil {
			// log.Printf("[BFS_DEBUG] [Loop %d] PERINGATAN: current.Parent adalah nil! Ini seharusnya tidak terjadi jika Parent di set dengan benar.", processedInQueue)
			// Jika ini terjadi, baris current.Parent.DibuatDari akan panic.
			// Untuk debug, kita bisa skip jika nil, tapi ini indikasi masalah besar.
			continue
		}
		current.Parent.DibuatDari = append(current.Parent.DibuatDari, current)
		// log.Printf("[BFS_DEBUG] [Loop %d] Menambahkan relasi ke '%s'.DibuatDari. Jumlah DibuatDari sekarang: %d", processedInQueue, current.Parent.NamaElemen, len(current.Parent.DibuatDari))

		// Proses LeftChild dari 'current'
		// log.Printf("[BFS_DEBUG] [Loop %d] Memproses LeftChild: '%s' (ID:%d)", processedInQueue, current.LeftChild.NamaElemen, current.LeftChildID)
		if IsBasicElement(*current.LeftChild) {
			// log.Printf("[BFS_DEBUG] [Loop %d] LeftChild '%s' adalah elemen dasar. Memanggil BfsHelper.", processedInQueue, current.LeftChild.NamaElemen)
			if current.LeftChild.Parent == nil {
				// log.Printf("[BFS_DEBUG] [Loop %d] PERINGATAN SEBELUM BFSHELPER (KIRI): current.LeftChild.Parent adalah nil! Seharusnya ini sudah diatur.", processedInQueue)
			} else if current.LeftChild.Parent != &current {
				// log.Printf("[BFS_DEBUG] [Loop %d] PERINGATAN SEBELUM BFSHELPER (KIRI): current.LeftChild.Parent TIDAK menunjuk ke &current!", processedInQueue)
			}
			BfsHelperORI(current.LeftChild)
		} else {
			// log.Printf("[BFS_DEBUG] [Loop %d] LeftChild '%s' (ID:%d) bukan elemen dasar. Mengekspansi...", processedInQueue, current.LeftChild.NamaElemen, current.LeftChildID)
			expandingParentNode := current.LeftChild   // Ini adalah *RecipeTreeNode yang akan menjadi .Parent untuk anak-anaknya
			expandingID_asIndex := current.LeftChildID // INI ADALAH ID, AKAN DIGUNAKAN SEBAGAI INDEKS (POTENSI MASALAH)

			if expandingID_asIndex < 0 || expandingID_asIndex >= len(elements) {
				// log.Printf("[BFS_DEBUG] [Loop %d ExpLeft] ERROR: expandingID_asIndex (%d) untuk '%s' di luar jangkauan. Tidak bisa ekspansi.", processedInQueue, expandingID_asIndex, current.LeftChild.NamaElemen)
			} else {
				// log.Printf("[BFS_DEBUG] [Loop %d ExpLeft] Mengekspansi '%s' (ID:%d) dari elements[%d]", processedInQueue, elements[expandingID_asIndex].Name, current.LeftChildID, expandingID_asIndex)
				if len(elements[expandingID_asIndex].FromPair) == 0 {
					// log.Printf("[BFS_DEBUG] [Loop %d ExpLeft] Node '%s' tidak punya FromPair untuk diekspansi.", processedInQueue, elements[expandingID_asIndex].Name)
				}
				for i := 0; i < len(elements[expandingID_asIndex].FromPair); i++ {
					pair := elements[expandingID_asIndex].FromPair[i]
					grandLeftID_fromPair := pair[0]
					grandRightID_fromPair := pair[1]
					// log.Printf("[BFS_DEBUG] [Loop %d ExpLeft] Memeriksa pair untuk '%s': GrandLeftID=%d, GrandRightID=%d", processedInQueue, elements[expandingID_asIndex].Name, grandLeftID_fromPair, grandRightID_fromPair)

					if IsValid(current.LeftChildID, grandLeftID_fromPair, grandRightID_fromPair) { // IsValid pakai ID (yg akan jadi indeks)
						// log.Printf("[BFS_DEBUG] [Loop %d ExpLeft] Pair (GLID:%d, GRID:%d) VALID.", processedInQueue, grandLeftID_fromPair, grandRightID_fromPair)

						// Mengakses elemen cucu menggunakan ID dari pair SEBAGAI INDEKS
						if grandLeftID_fromPair < 0 || grandLeftID_fromPair >= len(elements) || grandRightID_fromPair < 0 || grandRightID_fromPair >= len(elements) {
							// log.Printf("[BFS_DEBUG] [Loop %d ExpLeft] ERROR: ID cucu (GL:%d, GR:%d) dari FromPair di luar jangkauan jika sebagai INDEKS. Lewati.", processedInQueue, grandLeftID_fromPair, grandRightID_fromPair)
							continue
						}
						// log.Printf("[BFS_DEBUG] [Loop %d ExpLeft] Membuat GrandLeft dari elements[%d]: Nama='%s'", processedInQueue, grandLeftID_fromPair, elements[grandLeftID_fromPair].Name)
						// log.Printf("[BFS_DEBUG] [Loop %d ExpLeft] Membuat GrandRight dari elements[%d]: Nama='%s'", processedInQueue, grandRightID_fromPair, elements[grandRightID_fromPair].Name)

						tempGrandLeftNode := &RecipeTreeNode{NamaElemen: elements[grandLeftID_fromPair].Name, DibuatDari: make([]RecipeTreeNodeChild, 0)}
						tempGrandRightNode := &RecipeTreeNode{NamaElemen: elements[grandRightID_fromPair].Name, DibuatDari: make([]RecipeTreeNodeChild, 0)}

						var newChildRelation RecipeTreeNodeChild
						newChildRelation = RecipeTreeNodeChild{
							Parent:       expandingParentNode, // Parent-nya adalah current.LeftChild
							LeftChild:    tempGrandLeftNode,
							RightChild:   tempGrandRightNode,
							LeftChildID:  grandLeftID_fromPair,
							RightChildID: grandRightID_fromPair,
						}
						tempGrandLeftNode.Parent = &newChildRelation
						tempGrandRightNode.Parent = &newChildRelation
						// log.Printf("[BFS_DEBUG] [Loop %d ExpLeft] Pointer Parent untuk grand '%s' dan '%s' telah diatur.", processedInQueue, tempGrandLeftNode.NamaElemen, tempGrandRightNode.NamaElemen)

						queue = append(queue, newChildRelation)
						// log.Printf("[BFS_DEBUG] [Loop %d ExpLeft] Menambahkan ke queue dari ekspansi '%s': Left='%s', Right='%s'. Queue size: %d", processedInQueue, elements[expandingID_asIndex].Name, tempGrandLeftNode.NamaElemen, tempGrandRightNode.NamaElemen, len(queue))
					} else {
						// log.Printf("[BFS_DEBUG] [Loop %d ExpLeft] Pair (GLID:%d, GRID:%d) TIDAK VALID.", processedInQueue, grandLeftID_fromPair, grandRightID_fromPair)
					}
				}
			}
		}

		// Proses RightChild dari 'current' (logika serupa dengan LeftChild)
		// log.Printf("[BFS_DEBUG] [Loop %d] Memproses RightChild: '%s' (ID:%d)", processedInQueue, current.RightChild.NamaElemen, current.RightChildID)
		if IsBasicElement(*current.RightChild) {
			// log.Printf("[BFS_DEBUG] [Loop %d] RightChild '%s' adalah elemen dasar. Memanggil BfsHelper.", processedInQueue, current.RightChild.NamaElemen)
			if current.RightChild.Parent == nil {
				// log.Printf("[BFS_DEBUG] [Loop %d] PERINGATAN SEBELUM BFSHELPER (KANAN): current.RightChild.Parent adalah nil!", processedInQueue)
			} else if current.RightChild.Parent != &current {
				// log.Printf("[BFS_DEBUG] [Loop %d] PERINGATAN SEBELUM BFSHELPER (KANAN): current.RightChild.Parent TIDAK menunjuk ke &current!", processedInQueue)
			}
			BfsHelperORI(current.RightChild)
		} else {
			// log.Printf("[BFS_DEBUG] [Loop %d] RightChild '%s' (ID:%d) bukan elemen dasar. Mengekspansi...", processedInQueue, current.RightChild.NamaElemen, current.RightChildID)
			expandingParentNode := current.RightChild
			expandingID_asIndex := current.RightChildID // ID DIGUNAKAN SEBAGAI INDEKS

			if expandingID_asIndex < 0 || expandingID_asIndex >= len(elements) {
				// log.Printf("[BFS_DEBUG] [Loop %d ExpRight] ERROR: expandingID_asIndex (%d) untuk '%s' di luar jangkauan. Tidak bisa ekspansi.", processedInQueue, expandingID_asIndex, current.RightChild.NamaElemen)
			} else {
				// log.Printf("[BFS_DEBUG] [Loop %d ExpRight] Mengekspansi '%s' (ID:%d) dari elements[%d]", processedInQueue, elements[expandingID_asIndex].Name, current.RightChildID, expandingID_asIndex)
				if len(elements[expandingID_asIndex].FromPair) == 0 {
					// log.Printf("[BFS_DEBUG] [Loop %d ExpRight] Node '%s' tidak punya FromPair untuk diekspansi.", processedInQueue, elements[expandingID_asIndex].Name)
				}
				for i := 0; i < len(elements[expandingID_asIndex].FromPair); i++ {
					pair := elements[expandingID_asIndex].FromPair[i]
					grandLeftID_fromPair := pair[0]
					grandRightID_fromPair := pair[1]
					// log.Printf("[BFS_DEBUG] [Loop %d ExpRight] Memeriksa pair untuk '%s': GrandLeftID=%d, GrandRightID=%d", processedInQueue, elements[expandingID_asIndex].Name, grandLeftID_fromPair, grandRightID_fromPair)

					if IsValid(current.RightChildID, grandLeftID_fromPair, grandRightID_fromPair) { // IsValid pakai ID (yg akan jadi indeks)
						// log.Printf("[BFS_DEBUG] [Loop %d ExpRight] Pair (GLID:%d, GRID:%d) VALID.", processedInQueue, grandLeftID_fromPair, grandRightID_fromPair)

						if grandLeftID_fromPair < 0 || grandLeftID_fromPair >= len(elements) || grandRightID_fromPair < 0 || grandRightID_fromPair >= len(elements) {
							// log.Printf("[BFS_DEBUG] [Loop %d ExpRight] ERROR: ID cucu (GL:%d, GR:%d) dari FromPair di luar jangkauan jika sebagai INDEKS. Lewati.", processedInQueue, grandLeftID_fromPair, grandRightID_fromPair)
							continue
						}
						// log.Printf("[BFS_DEBUG] [Loop %d ExpRight] Membuat GrandLeft dari elements[%d]: Nama='%s'", processedInQueue, grandLeftID_fromPair, elements[grandLeftID_fromPair].Name)
						// log.Printf("[BFS_DEBUG] [Loop %d ExpRight] Membuat GrandRight dari elements[%d]: Nama='%s'", processedInQueue, grandRightID_fromPair, elements[grandRightID_fromPair].Name)

						tempGrandLeftNode := &RecipeTreeNode{NamaElemen: elements[grandLeftID_fromPair].Name, DibuatDari: make([]RecipeTreeNodeChild, 0)}
						tempGrandRightNode := &RecipeTreeNode{NamaElemen: elements[grandRightID_fromPair].Name, DibuatDari: make([]RecipeTreeNodeChild, 0)}

						var newChildRelation RecipeTreeNodeChild
						newChildRelation = RecipeTreeNodeChild{
							Parent:       expandingParentNode,
							LeftChild:    tempGrandLeftNode,
							RightChild:   tempGrandRightNode,
							LeftChildID:  grandLeftID_fromPair,
							RightChildID: grandRightID_fromPair,
						}
						tempGrandLeftNode.Parent = &newChildRelation
						tempGrandRightNode.Parent = &newChildRelation
						// log.Printf("[BFS_DEBUG] [Loop %d ExpRight] Pointer Parent untuk grand '%s' dan '%s' telah diatur.", processedInQueue, tempGrandLeftNode.NamaElemen, tempGrandRightNode.NamaElemen)

						queue = append(queue, newChildRelation)
						// log.Printf("[BFS_DEBUG] [Loop %d ExpRight] Menambahkan ke queue dari ekspansi '%s': Left='%s', Right='%s'. Queue size: %d", processedInQueue, elements[expandingID_asIndex].Name, tempGrandLeftNode.NamaElemen, tempGrandRightNode.NamaElemen, len(queue))
					} else {
						// log.Printf("[BFS_DEBUG] [Loop %d ExpRight] Pair (GLID:%d, GRID:%d) TIDAK VALID.", processedInQueue, grandLeftID_fromPair, grandRightID_fromPair)
					}
				}
			}
		}

		// log.Printf("[BFS_DEBUG] [Loop %d] Selesai memproses current. Tree.BanyakResep saat ini: %d", processedInQueue, Tree.BanyakResep)
		if Tree.BanyakResep >= needFound {
			// log.Printf("[BFS_DEBUG] [Loop %d] Tree.BanyakResep (%d) >= needFound (%d). Menghentikan BFS.", processedInQueue, Tree.BanyakResep, needFound)
			break
		}
	}

	var res = &Tree
	// log.Printf("[BFS_DEBUG] === Selesai BFS untuk elemen awal (di INDEKS %d) '%s'. Mengembalikan Tree.BanyakResep: %d. Jumlah DibuatDari: %d ===", start, res.NamaElemen, res.BanyakResep, len(res.DibuatDari))
	return BFSCleaner(res)
}


func BFSCleaner(tree *RecipeTreeNode) *RecipeTreeNode {
	if tree == nil || tree.BanyakResep == 0 {
		return nil
	}

	newTree := &RecipeTreeNode{
		NamaElemen:  tree.NamaElemen,
		DibuatDari:  make([]RecipeTreeNodeChild, 0),
		BanyakResep: tree.BanyakResep,
	}

	for _, child := range tree.DibuatDari {
		leftCleaned := BFSCleaner(child.LeftChild)
		rightCleaned := BFSCleaner(child.RightChild)

		if leftCleaned != nil && rightCleaned != nil {
			newChild := RecipeTreeNodeChild{
				Parent:       newTree,
				LeftChild:    leftCleaned,
				RightChild:   rightCleaned,
				LeftChildID:  child.LeftChildID,
				RightChildID: child.RightChildID,
			}
			newTree.DibuatDari = append(newTree.DibuatDari, newChild)
		}
	}

	if len(newTree.DibuatDari) == 0 && len(tree.DibuatDari) > 0 {
		return nil
	}

	return newTree
}
