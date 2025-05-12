package scraper

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/graceevelyns/Tubes2_BE_ian/src/cmd/internal/model"
)

const targetURL = "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

type RawRecipeEntry struct {
	ResultElement string
	Ingredient1   string
	Ingredient2   string
}

type Element struct {
	ID       int     `json:"Id"`
	Name     string  `json:"Name"`
	Tier     int     `json:"Tier"`
	FromPair [][]int `json:"FromPair"`
	CanMake  []int   `json:"CanMake"`
}

var (
	standardBaseElements   = []string{"Air", "Earth", "Fire", "Water", "Time"}
	processedElementsCache []*Element
	cacheMutex             sync.RWMutex
)

func isStandardBaseElement(name string) bool {
	normalizedName := strings.Title(strings.ToLower(name))
	for _, baseName := range standardBaseElements {
		if normalizedName == baseName {
			return true
		}
	}
	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func FetchAndProcessData() ([]*Element, error) {
	log.Println("Memulai proses scraping dari:", targetURL)

	res, err := http.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("gagal GET request ke %s: %w", targetURL, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code tidak OK: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal mem-parse HTML: %w", err)
	}

	log.Println("HTML berhasil di-parse. Memulai ekstraksi data...")
	var rawRecipes []RawRecipeEntry
	validResultElementsFromPage := make(map[string]bool)
	var tempOrderedResultNames []string
	isNameInTempOrderedList := make(map[string]bool)

	doc.Find("div.mw-parser-output table.list-table tbody tr").Each(func(i int, rowSelection *goquery.Selection) {
		tds := rowSelection.Find("td")
		if tds.Length() < 1 {
			return
		}

		resultElementNode := tds.Eq(0).Find("a[href^='/wiki/']").First()
		resultElementName := strings.TrimSpace(resultElementNode.Text())
		if resultElementName == "" {
			resultElementName = strings.TrimSpace(tds.Eq(0).Contents().Not("span, script, style").Text())
			if strings.Contains(resultElementName, "\n") {
				resultElementName = strings.TrimSpace(strings.Split(resultElementName, "\n")[0])
			}
		}
		if resultElementName == "" {
			return
		}
		resultElementName = strings.Title(strings.ToLower(resultElementName))

		if tds.Length() > 1 {
			mythLink := tds.Eq(1).Find("a[href='/wiki/Elements_(Myths_and_Monsters)']")
			if mythLink.Length() > 0 {
				log.Printf("FILTER (MYTHS LINK): Melewati elemen '%s'.", resultElementName)
				return
			}
		}
		var isVisuallyMarkedAsPack = false
		tds.Eq(0).Find("a[title*='Myths and Monsters Pack'], img[alt*='Myths and Monsters Pack'], span[class*='pack-icon']").Each(func(_ int, s *goquery.Selection) {
			if !isStandardBaseElement(resultElementName) {
				log.Printf("FILTER (VISUAL PACK): Melewati elemen '%s'.", resultElementName)
				isVisuallyMarkedAsPack = true
			}
		})
		if isVisuallyMarkedAsPack {
			return
		}

		validResultElementsFromPage[resultElementName] = true
		if !isNameInTempOrderedList[resultElementName] {
			tempOrderedResultNames = append(tempOrderedResultNames, resultElementName)
			isNameInTempOrderedList[resultElementName] = true
		}

		if tds.Length() > 1 {
			tds.Eq(1).Find("ul > li").Each(func(j int, recipeLi *goquery.Selection) {
				ingredientLinks := recipeLi.Find("a[href^='/wiki/']")
				if ingredientLinks.Length() >= 2 {
					ing1 := strings.Title(strings.ToLower(strings.TrimSpace(ingredientLinks.Eq(0).Text()))) // Normalisasi
					ing2 := strings.Title(strings.ToLower(strings.TrimSpace(ingredientLinks.Eq(1).Text()))) // Normalisasi
					if ing1 != "" && ing2 != "" {
						sortedIng := []string{ing1, ing2}
						sort.Strings(sortedIng)
						rawRecipes = append(rawRecipes, RawRecipeEntry{resultElementName, sortedIng[0], sortedIng[1]})
					}
				}
			})
		}
	})
	log.Printf("Scraping mentah selesai. %d elemen hasil unik potensial. %d entri resep potensial.", len(tempOrderedResultNames), len(rawRecipes))

	var finalRawRecipes []RawRecipeEntry
	for _, rr := range rawRecipes {
		if validResultElementsFromPage[rr.ResultElement] {
			finalRawRecipes = append(finalRawRecipes, rr)
		} else {
			log.Printf("FILTER (RAW RECIPE): Membuang resep dengan hasil '%s'.", rr.ResultElement)
		}
	}
	rawRecipes = finalRawRecipes
	log.Printf("Setelah filter rawRecipes, tersisa %d entri resep.", len(rawRecipes))

	allNodes := make(map[string]*model.RecipeTreeNodeTier)
	var finalOrderedNodeKeys []string
	for _, name := range tempOrderedResultNames {
		allNodes[name] = &model.RecipeTreeNodeTier{NamaElemen: name, DibuatDari: make([][2]*model.RecipeTreeNodeTier, 0)}
		finalOrderedNodeKeys = append(finalOrderedNodeKeys, name)
	}
	tempIsNameInFinalOrderedList := make(map[string]bool)
	for _, name := range finalOrderedNodeKeys {
		tempIsNameInFinalOrderedList[name] = true
	}
	for _, baseName := range standardBaseElements {
		normBaseName := strings.Title(strings.ToLower(baseName))
		if _, exists := allNodes[normBaseName]; !exists {
			log.Printf("Info: Elemen dasar '%s' dibuat manual.", normBaseName)
			allNodes[normBaseName] = &model.RecipeTreeNodeTier{NamaElemen: normBaseName, DibuatDari: make([][2]*model.RecipeTreeNodeTier, 0)}
			if !tempIsNameInFinalOrderedList[normBaseName] {
				finalOrderedNodeKeys = append(finalOrderedNodeKeys, normBaseName)
			}
		}
	}
	uniqueKeysMap := make(map[string]bool)
	var cleanedKeys []string
	for _, k := range finalOrderedNodeKeys {
		if !uniqueKeysMap[k] {
			uniqueKeysMap[k] = true
			cleanedKeys = append(cleanedKeys, k)
		}
	}
	finalOrderedNodeKeys = cleanedKeys
	log.Printf("Node awal dibuat. %d node unik. %d kunci terurut final.", len(allNodes), len(finalOrderedNodeKeys))

	for _, rr := range rawRecipes {
		resNode, rExists := allNodes[rr.ResultElement]
		if !rExists || isStandardBaseElement(resNode.NamaElemen) {
			continue
		}
		ing1Node, i1Exists := allNodes[rr.Ingredient1]
		ing2Node, i2Exists := allNodes[rr.Ingredient2]
		if i1Exists && i2Exists {
			comboExists := false
			for _, combo := range resNode.DibuatDari {
				if combo[0].NamaElemen == rr.Ingredient1 && combo[1].NamaElemen == rr.Ingredient2 {
					comboExists = true
					break
				}
			}
			if !comboExists {
				resNode.DibuatDari = append(resNode.DibuatDari, [2]*model.RecipeTreeNodeTier{ing1Node, ing2Node})
			}
		} else {
			if !i1Exists {
				log.Printf("Peringatan (Node Build): Bahan '%s' untuk '%s' tidak ada.", rr.Ingredient1, rr.ResultElement)
			}
			if !i2Exists {
				log.Printf("Peringatan (Node Build): Bahan '%s' untuk '%s' tidak ada.", rr.Ingredient2, rr.ResultElement)
			}
		}
	}
	log.Printf("Struktur RecipeTreeNodeTier selesai. Total node: %d.", len(allNodes))

	nameToID := make(map[string]int)
	elementMapByID := make(map[int]*Element)
	var orderedElements []*Element

	nextID := 0
	for _, name := range finalOrderedNodeKeys {
		if _, nodeExists := allNodes[name]; nodeExists {
			currentID := nextID
			nameToID[name] = currentID
			el := &Element{
				ID: currentID, Name: name, Tier: -1,
				FromPair: make([][]int, 0), CanMake: make([]int, 0),
			}
			elementMapByID[currentID] = el
			orderedElements = append(orderedElements, el)
			nextID++
		}
	}
	log.Printf("Penetapan ID selesai. %d elemen mendapatkan ID.", len(orderedElements))
	if len(orderedElements) != len(allNodes) {
		log.Printf("PERINGATAN (ID Assign): Jumlah elemen dengan ID (%d) tidak sama dengan jumlah node di allNodes (%d).", len(orderedElements), len(allNodes))
	}

	// kalkulasi tier
	elementTiersByName := make(map[string]int)
	for name := range allNodes {
		elementTiersByName[name] = -1
	}
	for name := range allNodes {
		if isStandardBaseElement(name) {
			elementTiersByName[name] = 0
		}
	}
	maxIter := len(allNodes) + 5
	for i := 0; i < maxIter; i++ {
		changed := false
		for elName, node := range allNodes {
			if tier, _ := elementTiersByName[elName]; tier == 0 {
				continue
			}
			minTier := -1
			if (node.DibuatDari == nil || len(node.DibuatDari) == 0) && !isStandardBaseElement(elName) {
				continue
			}
			for _, pair := range node.DibuatDari {
				if pair[0] == nil || pair[1] == nil {
					continue
				}
				t1, ok1 := elementTiersByName[pair[0].NamaElemen]
				t2, ok2 := elementTiersByName[pair[1].NamaElemen]
				if ok1 && t1 != -1 && ok2 && t2 != -1 {
					calcTier := 1 + max(t1, t2)
					if minTier == -1 || calcTier < minTier {
						minTier = calcTier
					}
				}
			}
			if minTier != -1 {
				curTier, _ := elementTiersByName[elName]
				if curTier == -1 || minTier < curTier {
					elementTiersByName[elName] = minTier
					changed = true
				}
			}
		}
		if !changed {
			log.Printf("Kalkulasi tier konvergen iterasi %d.", i+1)
			break
		}
		if i == maxIter-1 {
			log.Println("Peringatan: Kalkulasi tier mencapai max iterasi.")
		}
	}
	for _, el := range orderedElements {
		if tier, ok := elementTiersByName[el.Name]; ok {
			el.Tier = tier
		} else {
			el.Tier = -1
			log.Printf("Peringatan (Tier Assign): Tier untuk '%s' tidak ditemukan.", el.Name)
		}
	}
	log.Println("Kalkulasi tier selesai.")

	canMakeTemp := make(map[int]map[int]bool)
	for _, currentElement := range orderedElements {
		resultNode, nodeExists := allNodes[currentElement.Name]
		if !nodeExists || isStandardBaseElement(currentElement.Name) || resultNode.DibuatDari == nil {
			continue
		}
		processedPairs := make(map[string]bool)
		for _, recipe := range resultNode.DibuatDari {
			ing1Name := recipe[0].NamaElemen
			ing2Name := recipe[1].NamaElemen
			ing1ID, ok1 := nameToID[ing1Name]
			ing2ID, ok2 := nameToID[ing2Name]
			if ok1 && ok2 {
				pairIDs := []int{ing1ID, ing2ID}
				if pairIDs[0] > pairIDs[1] {
					pairIDs[0], pairIDs[1] = pairIDs[1], pairIDs[0]
				}
				pairKey := fmt.Sprintf("%d-%d", pairIDs[0], pairIDs[1])
				if !processedPairs[pairKey] {
					currentElement.FromPair = append(currentElement.FromPair, pairIDs)
					processedPairs[pairKey] = true
				}
				if _, exists := canMakeTemp[ing1ID]; !exists {
					canMakeTemp[ing1ID] = make(map[int]bool)
				}
				canMakeTemp[ing1ID][currentElement.ID] = true
				if ing1ID != ing2ID {
					if _, exists := canMakeTemp[ing2ID]; !exists {
						canMakeTemp[ing2ID] = make(map[int]bool)
					}
					canMakeTemp[ing2ID][currentElement.ID] = true
				}
			}
		}
	}
	for ingID, producesMap := range canMakeTemp {
		if ingElement, ok := elementMapByID[ingID]; ok {
			for prodID := range producesMap {
				ingElement.CanMake = append(ingElement.CanMake, prodID)
			}
			sort.Ints(ingElement.CanMake)
		}
	}
	log.Println("Populasi FromPair dan CanMake selesai.")

	var finalFilteredElements []*Element
	for _, el := range orderedElements {
		isBase := isStandardBaseElement(el.Name)
		hasValidRecipes := len(el.FromPair) > 0
		if isBase || (hasValidRecipes && el.Tier != -1) {
			finalFilteredElements = append(finalFilteredElements, el)
		}
	}
	log.Printf("Filter akhir: %d elemen akan dikembalikan.", len(finalFilteredElements))

	cacheMutex.Lock()
	processedElementsCache = finalFilteredElements
	cacheMutex.Unlock()

	log.Println("Semua proses data dalam scraper selesai. Data disimpan di cache internal.")
	return finalFilteredElements, nil
}

func GetProcessedElements() []*Element {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	if processedElementsCache == nil {
		log.Println("Peringatan: GetProcessedElements dipanggil sebelum cache diisi. Kembalikan slice kosong.")
		return []*Element{}
	}

	elementsCopy := make([]*Element, len(processedElementsCache))
	for i, el := range processedElementsCache {

		tempEl := *el

		tempEl.FromPair = make([][]int, len(el.FromPair))
		idx := 0
		for j, pair := range el.FromPair {
			if isElementValid(processedElementsCache, el.ID, pair[0], pair[1]) { 
				tempEl.FromPair[j] = make([]int, len(pair))
				copy(tempEl.FromPair[j], pair)
				idx++
			}
		}

		tempEl.CanMake = make([]int, len(el.CanMake))
		copy(tempEl.CanMake, el.CanMake)

		elementsCopy[i] = &tempEl
	}
	return elementsCopy
}


func isElementValid(el []*Element, idParent int, idChild1 int, idChild2 int) bool{
	if el[idParent].Tier <= el[idChild1].Tier || el[idParent].Tier <= el[idChild2].Tier {
		return false
	}
	if el[idChild1].Name == "Time" || el[idChild2].Name == "Time" {
		return false
	}
	if el[idChild1].Name == "Ruins" || el[idChild2].Name == "Ruins" {
		return false
	}
	if el[idChild1].Name == "Archeologist" || el[idChild2].Name == "Archeologist" {
		return false
	}
	return true
}
