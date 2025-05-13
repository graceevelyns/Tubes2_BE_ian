package algorithm

import (
	"fmt" // Tambahkan untuk formatting log
	"log"

	// "log"
	"sync"
	"sync/atomic"

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

func BfsHelperORI(a *RecipeTreeNode) {
	kiri := (a.Parent.LeftChild == a)
	newBanyakResep := 1
	for a.Parent != nil {
		if a.Parent.LeftChild != nil && a.Parent.RightChild != nil {
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
			var childRelation = RecipeTreeNodeChild{
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

		if current.Parent == nil {
			continue
		}
		current.Parent.DibuatDari = append(current.Parent.DibuatDari, current)

		// Proses LeftChild dari 'current'
		if IsBasicElement(*current.LeftChild) {

			BfsHelperORI(current.LeftChild)
		} else {
			expandingParentNode := current.LeftChild   // Ini adalah *RecipeTreeNode yang akan menjadi .Parent untuk anak-anaknya
			expandingID_asIndex := current.LeftChildID // INI ADALAH ID, AKAN DIGUNAKAN SEBAGAI INDEKS (POTENSI MASALAH)

			if expandingID_asIndex < 0 || expandingID_asIndex >= len(elements) {
			} else {

				for i := 0; i < len(elements[expandingID_asIndex].FromPair); i++ {
					pair := elements[expandingID_asIndex].FromPair[i]
					grandLeftID_fromPair := pair[0]
					grandRightID_fromPair := pair[1]

					if IsValid(current.LeftChildID, grandLeftID_fromPair, grandRightID_fromPair) { // IsValid pakai ID (yg akan jadi indeks)

						if grandLeftID_fromPair < 0 || grandLeftID_fromPair >= len(elements) || grandRightID_fromPair < 0 || grandRightID_fromPair >= len(elements) {
							continue
						}

						tempGrandLeftNode := &RecipeTreeNode{NamaElemen: elements[grandLeftID_fromPair].Name, DibuatDari: make([]RecipeTreeNodeChild, 0)}
						tempGrandRightNode := &RecipeTreeNode{NamaElemen: elements[grandRightID_fromPair].Name, DibuatDari: make([]RecipeTreeNodeChild, 0)}

						var newChildRelation = RecipeTreeNodeChild{
							Parent:       expandingParentNode, // Parent-nya adalah current.LeftChild
							LeftChild:    tempGrandLeftNode,
							RightChild:   tempGrandRightNode,
							LeftChildID:  grandLeftID_fromPair,
							RightChildID: grandRightID_fromPair,
						}
						tempGrandLeftNode.Parent = &newChildRelation
						tempGrandRightNode.Parent = &newChildRelation

						queue = append(queue, newChildRelation)
					}
				}
			}
		}

		// Proses RightChild dari 'current' (logika serupa dengan LeftChild)
		if IsBasicElement(*current.RightChild) {
			BfsHelperORI(current.RightChild)
		} else {
			expandingParentNode := current.RightChild
			expandingID_asIndex := current.RightChildID // ID DIGUNAKAN SEBAGAI INDEKS

			if expandingID_asIndex < 0 || expandingID_asIndex >= len(elements) {
			} else {

				for i := 0; i < len(elements[expandingID_asIndex].FromPair); i++ {
					pair := elements[expandingID_asIndex].FromPair[i]
					grandLeftID_fromPair := pair[0]
					grandRightID_fromPair := pair[1]

					if IsValid(current.RightChildID, grandLeftID_fromPair, grandRightID_fromPair) { // IsValid pakai ID (yg akan jadi indeks)

						if grandLeftID_fromPair < 0 || grandLeftID_fromPair >= len(elements) || grandRightID_fromPair < 0 || grandRightID_fromPair >= len(elements) {
							continue
						}

						tempGrandLeftNode := &RecipeTreeNode{NamaElemen: elements[grandLeftID_fromPair].Name, DibuatDari: make([]RecipeTreeNodeChild, 0)}
						tempGrandRightNode := &RecipeTreeNode{NamaElemen: elements[grandRightID_fromPair].Name, DibuatDari: make([]RecipeTreeNodeChild, 0)}

						var newChildRelation = RecipeTreeNodeChild{
							Parent:       expandingParentNode,
							LeftChild:    tempGrandLeftNode,
							RightChild:   tempGrandRightNode,
							LeftChildID:  grandLeftID_fromPair,
							RightChildID: grandRightID_fromPair,
						}
						tempGrandLeftNode.Parent = &newChildRelation
						tempGrandRightNode.Parent = &newChildRelation

						queue = append(queue, newChildRelation)
					}
				}
			}
		}

		if Tree.BanyakResep >= needFound {
			break
		}
	}

	var res = &Tree
	return BFSCleaner(res)
}

func BFSCleaner(tree *RecipeTreeNode) *RecipeTreeNode {
	log.Println("[BFSCleaner] Memulai pembersihan tree untuk: %s", tree.NamaElemen)
	log.Println("[BFSCleaner] BanyakResep: %d", tree.BanyakResep)
	if tree == nil || tree.BanyakResep == 0 {
		return nil
	}

	newTree := &RecipeTreeNode{
		NamaElemen:  tree.NamaElemen,
		DibuatDari:  make([]RecipeTreeNodeChild, 0),
		BanyakResep: tree.BanyakResep,
	}

	for _, child := range tree.DibuatDari {
		log.Println("[BFSCleaner] Memproses child dari tree: %s", child.Parent.NamaElemen)
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

func ParallelBfs(targetID, needCount int) *RecipeTreeNode {
	// log.Print("[PARALLEL_BFS_DEBUG] === Memulai Parallel BFS untuk targetID: %d, needCount: %d. Panjang elements: %d ===", targetID, needCount, len(elements))
	queue := make([]RecipeTreeNodeChild, 0)

	if targetID < 0 || targetID >= len(elements) {
		return nil
	}

	root := RecipeTreeNode{
		Parent:      nil,
		NamaElemen:  elements[targetID].Name,
		DibuatDari:  []RecipeTreeNodeChild{},
		BanyakResep: 0,
	}

	var queueMu sync.Mutex
	var treeMu sync.Mutex // Protects all tree modifications
	var resepCount int32 = 0
	var wg sync.WaitGroup

	// Initial queue population with validity checks
	for _, pair := range elements[targetID].FromPair {
		if !IsValid(targetID, pair[0], pair[1]) {
			continue
		}

		left := &RecipeTreeNode{
			NamaElemen:  elements[pair[0]].Name,
			DibuatDari:  make([]RecipeTreeNodeChild, 0),
			BanyakResep: 0,
		}
		right := &RecipeTreeNode{
			NamaElemen:  elements[pair[1]].Name,
			DibuatDari:  make([]RecipeTreeNodeChild, 0),
			BanyakResep: 0,
		}

		child := RecipeTreeNodeChild{
			Parent:       &root,
			LeftChild:    left,
			RightChild:   right,
			LeftChildID:  pair[0],
			RightChildID: pair[1],
		}
		left.Parent = &child
		right.Parent = &child

		queue = append(queue, child)
	}

	for len(queue) > 0 && atomic.LoadInt32(&resepCount) < int32(needCount) {
		// log.Println("[PARALLEL_BFS_DEBUG] Queue size:", len(queue))
		queueMu.Lock()
		if len(queue) == 0 {
			queueMu.Unlock()
			break
		}
		item := queue[0]
		queue = queue[1:]

		queueMu.Unlock()
		wg.Add(1)
		go func(node RecipeTreeNodeChild) {
			defer wg.Done()
			// log.Println("[PARALLEL_BFS_DEBUG] Memproses node:", node.LeftChild.NamaElemen, "dan", node.RightChild.NamaElemen)

			// Add relation to parent with lock
			treeMu.Lock()
			node.Parent.DibuatDari = append(node.Parent.DibuatDari, node)
			treeMu.Unlock()

			processChild := func(childID int, childNode *RecipeTreeNode) {
				if childID < 0 || childID >= len(elements) {
					return
				}

				if IsBasicElement_typeElement(*elements[childID]) {
					// log.Println("[PARALLEL_BFS_DEBUG] Memproses elemen dasar:", childNode.NamaElemen)
					// tempChildNode := childNode
					treeMu.Lock()
					BfsHelperORI(childNode)
					atomic.StoreInt32(&resepCount, int32(root.BanyakResep))
					treeMu.Unlock()
				} else {
					// log.Println("[PARALLEL_BFS_DEBUG] Memproses anak dari node:", childNode.NamaElemen)
					for _, pair := range elements[childID].FromPair {
						if !IsValid(childID, pair[0], pair[1]) {
							// log.Println("[PARALLEL_BFS_DEBUG] Pair tidak valid:", pair, "dari ID:", childID)
							continue
						}

						left := &RecipeTreeNode{
							NamaElemen:  elements[pair[0]].Name,
							DibuatDari:  make([]RecipeTreeNodeChild, 0),
							BanyakResep: 0,
						}
						right := &RecipeTreeNode{
							NamaElemen:  elements[pair[1]].Name,
							DibuatDari:  make([]RecipeTreeNodeChild, 0),
							BanyakResep: 0,
						}

						newChild := RecipeTreeNodeChild{
							Parent:       childNode,
							LeftChild:    left,
							RightChild:   right,
							LeftChildID:  pair[0],
							RightChildID: pair[1],
						}
						left.Parent = &newChild
						right.Parent = &newChild

						// Add to queue safely
						queueMu.Lock()
						queue = append(queue, newChild)
						queueMu.Unlock()
						// log.Println("[PARALLEL_BFS_DEBUG] Menambahkan ke queue:", newChild.LeftChild.NamaElemen, newChild.RightChild.NamaElemen, "Parent:", newChild.Parent.NamaElemen)
					}
				}
			}
			// log.Println("[PARALLEL_BFS_DEBUG-1] Queue size:", len(queue))
			processChild(node.LeftChildID, node.LeftChild)
			// log.Println("[PARALLEL_BFS_DEBUG0] Queue size:", len(queue))
			processChild(node.RightChildID, node.RightChild)
			// log.Println("[PARALLEL_BFS_DEBUG1] Queue size:", len(queue))
		}(item)
		wg.Wait()
		// log.Println("[PARALLEL_BFS_DEBUG2] Queue size:", len(queue))
	}
	// log.Println("[PARALLEL_BFS_DEBUG3] Queue size:", len(queue))

	// root.BanyakResep = int(atomic.LoadInt32(&resepCount))
	// root.BanyakResep = 3
	return BFSCleaner(&root)
}
