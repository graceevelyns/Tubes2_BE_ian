package scraper

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/graceevelyns/Tubes2_BE_ian/src/internal/model"
)

const targetURL = "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

type RawRecipeEntry struct {
	ResultElement string
	Ingredient1   string
	Ingredient2   string
}

var (
	standardBaseElements = []string{"Air", "Earth", "Fire", "Water", "Time"}
)

func isStandardBaseElement(name string) bool {
	for _, baseName := range standardBaseElements {
		if name == baseName {
			return true
		}
	}
	return false
}

func FetchAndParseData() (map[string]*model.RecipeNode, []*model.RecipeNode, []string, error) {
	log.Println("Memulai proses scraping dari:", targetURL)

	res, err := http.Get(targetURL)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("gagal GET request ke %s: %w", targetURL, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, nil, nil, fmt.Errorf("status code tidak OK: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("gagal mem-parse HTML: %w", err)
	}

	log.Println("HTML berhasil di-parse. Memulai ekstraksi data...")
	var rawRecipes []RawRecipeEntry
	validResultElementsFromPage := make(map[string]bool)

	var tempOrderedResultNames []string
	isNameInTempOrderedList := make(map[string]bool)
	elementIsBaseFromRow := make(map[string]bool)

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

		isBaseCurrentRow := false
		if tds.Length() > 1 {
			descriptionText := strings.TrimSpace(tds.Eq(1).Contents().Not("ul, script, style").Text())
			descriptionTextLower := strings.ToLower(descriptionText)
			if strings.Contains(descriptionTextLower, "available from the start") || strings.Contains(descriptionTextLower, "this element does not have any recipes") {
				if isStandardBaseElement(resultElementName) {
					isBaseCurrentRow = true
				}
			}
		}

		if !isNameInTempOrderedList[resultElementName] {
			tempOrderedResultNames = append(tempOrderedResultNames, resultElementName)
			isNameInTempOrderedList[resultElementName] = true
			elementIsBaseFromRow[resultElementName] = isBaseCurrentRow
		}

		if tds.Length() > 1 {
			tds.Eq(1).Find("ul > li").Each(func(j int, recipeLi *goquery.Selection) {
				ingredientLinks := recipeLi.Find("a[href^='/wiki/']")
				if ingredientLinks.Length() >= 2 {
					ingredient1 := strings.TrimSpace(ingredientLinks.Eq(0).Text())
					ingredient2 := strings.TrimSpace(ingredientLinks.Eq(1).Text())
					if ingredient1 != "" && ingredient2 != "" {
						sortedIngredients := []string{ingredient1, ingredient2}
						sort.Strings(sortedIngredients)
						rawRecipes = append(rawRecipes, RawRecipeEntry{
							ResultElement: resultElementName,
							Ingredient1:   sortedIngredients[0],
							Ingredient2:   sortedIngredients[1],
						})
					}
				}
			})
		}
	})

	log.Printf("Scraping mentah selesai. %d elemen hasil unik potensial ditemukan berurutan. %d entri resep potensial.", len(tempOrderedResultNames), len(rawRecipes))

	var finalRawRecipes []RawRecipeEntry
	for _, rr := range rawRecipes {
		if validResultElementsFromPage[rr.ResultElement] {
			finalRawRecipes = append(finalRawRecipes, rr)
		} else {
			log.Printf("FILTER (RAW RECIPE): Membuang resep dengan hasil '%s' karena hasil tidak valid.", rr.ResultElement)
		}
	}
	rawRecipes = finalRawRecipes
	log.Printf("Setelah filter rawRecipes, tersisa %d entri resep.", len(rawRecipes))

	allNodes := make(map[string]*model.RecipeNode)
	var finalOrderedNodeKeys []string

	for _, elementName := range tempOrderedResultNames {
		isBase := elementIsBaseFromRow[elementName]
		if elementName == "Time" {
			isBase = true
		}

		node := &model.RecipeNode{NamaElemen: elementName, IsBaseElement: isBase, DibuatDari: make([][2]*model.RecipeNode, 0)}
		allNodes[elementName] = node
		finalOrderedNodeKeys = append(finalOrderedNodeKeys, elementName)
	}

	for _, baseName := range standardBaseElements {
		if node, exists := allNodes[baseName]; exists {
			if !node.IsBaseElement {
				log.Printf("Info: Mengoreksi flag IsBaseElement untuk '%s' menjadi true.", baseName)
				node.IsBaseElement = true
			}
		} else {
			log.Printf("Info: Elemen dasar standar '%s' tidak ditemukan dari scraping baris. Dibuat manual dan ditambahkan ke akhir daftar terurut.", baseName)
			node := &model.RecipeNode{NamaElemen: baseName, IsBaseElement: true, DibuatDari: make([][2]*model.RecipeNode, 0)}
			allNodes[baseName] = node
			finalOrderedNodeKeys = append(finalOrderedNodeKeys, baseName)
		}
	}

	var baseElementNodes []*model.RecipeNode
	tempBaseSet := make(map[string]bool)
	for _, key := range finalOrderedNodeKeys {
		if node, ok := allNodes[key]; ok && node.IsBaseElement {
			if !tempBaseSet[node.NamaElemen] {
				baseElementNodes = append(baseElementNodes, node)
				tempBaseSet[node.NamaElemen] = true
			}
		}
	}
	sort.Slice(baseElementNodes, func(i, j int) bool {
		return baseElementNodes[i].NamaElemen < baseElementNodes[j].NamaElemen
	})

	for _, rr := range rawRecipes {
		resultNode, rExists := allNodes[rr.ResultElement]
		if !rExists || resultNode.IsBaseElement {
			continue
		}

		if rr.Ingredient1 != "" && rr.Ingredient2 != "" {
			ing1Node, i1Exists := allNodes[rr.Ingredient1]
			ing2Node, i2Exists := allNodes[rr.Ingredient2]

			if i1Exists && i2Exists {
				combinationExists := false
				for _, existingCombo := range resultNode.DibuatDari {
					if existingCombo[0].NamaElemen == rr.Ingredient1 && existingCombo[1].NamaElemen == rr.Ingredient2 {
						combinationExists = true
						break
					}
				}
				if !combinationExists {
					resultNode.DibuatDari = append(resultNode.DibuatDari, [2]*model.RecipeNode{ing1Node, ing2Node})
				}
			}
		}
	}

	log.Printf("Pembuatan struktur node selesai. Total node unik: %d. Elemen dasar teridentifikasi: %d. Kunci terurut: %d\n", len(allNodes), len(baseElementNodes), len(finalOrderedNodeKeys))
	if len(finalOrderedNodeKeys) != len(allNodes) {
		log.Printf("PERINGATAN KRITIS: Panjang finalOrderedNodeKeys (%d) tidak sama dengan jumlah node di allNodes (%d)!", len(finalOrderedNodeKeys), len(allNodes))
	}

	return allNodes, baseElementNodes, finalOrderedNodeKeys, nil
}
