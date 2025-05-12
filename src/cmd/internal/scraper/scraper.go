package scraper

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/graceevelyns/Tubes2_BE_ian/src/cmd/internal/model"
)

const (
	targetURL = "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
)

type RawRecipeEntry struct {
	ResultElement string
	Ingredient1   string
	Ingredient2   string
}

type Element struct {
	ID       int     `json:"Id"`       // unique identifier for the element
	Name     string  `json:"Name"`     // name of the element
	Tier     int     `json:"Tier"`     // tier of the element, 0 for base elements
	FromPair [][]int `json:"FromPair"` // list of ingredient ID pairs that can create this element
	CanMake  []int   `json:"CanMake"`  // list of element IDs that can be made using this element as an ingredient
}

var (
	standardBaseElements     = []string{"Air", "Earth", "Fire", "Water", "Time"}
	forbiddenIngredientNames = map[string]bool{
		"Time":         true,
		"Ruins":        true,
		"Archeologist": true,
	}
)

func normalizeElementName(name string) string {
	return strings.Title(strings.ToLower(strings.TrimSpace(name)))
}

// isStandardBaseElement checks if the given element name is one of the standard base elements
func isStandardBaseElement(name string) bool {
	normalizedName := normalizeElementName(name)
	for _, baseName := range standardBaseElements {
		if normalizedName == baseName {
			return true
		}
	}
	return false
}

func isRecipeValidInternal(ing1ID int, ing2ID int, resultID int, elementMapByID map[int]*Element) bool {
	ing1Elem, ok1 := elementMapByID[ing1ID]
	ing2Elem, ok2 := elementMapByID[ing2ID]
	resultElem, okResult := elementMapByID[resultID]

	if !ok1 || !ok2 || !okResult {
		return false
	}

	if forbiddenIngredientNames[ing1Elem.Name] {
		return false
	}
	if forbiddenIngredientNames[ing2Elem.Name] {
		return false
	}

	// check if tier is valid (not -1, which indicates an unprocessed or invalid tier)
	if ing1Elem.Tier == -1 || ing2Elem.Tier == -1 || resultElem.Tier == -1 {
		return false
	}

	return true
}

func fetchDocument(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed GET request to %s: %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code not OK: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}
	return doc, nil
}

func extractElementTiers(doc *goquery.Document) map[string]int {
	elementTiers := make(map[string]int)
	for _, baseName := range standardBaseElements {
		elementTiers[normalizeElementName(baseName)] = 0 // base elements are tier 0
	}

	doc.Find("div.mw-parser-output > h2, div.mw-parser-output > h3").Each(func(_ int, headingSel *goquery.Selection) {
		headlineSpan := headingSel.Find("span.mw-headline")
		if headlineSpan.Length() == 0 {
			return
		}
		headlineID, idExists := headlineSpan.Attr("id")
		if !idExists {
			return
		}

		currentTierNum := -1
		if strings.HasPrefix(headlineID, "Tier_") && strings.HasSuffix(headlineID, "_elements") {
			tierStr := strings.TrimSuffix(strings.TrimPrefix(headlineID, "Tier_"), "_elements")
			if val, errConv := strconv.Atoi(tierStr); errConv == nil {
				currentTierNum = val
			}
		} else if headlineID == "Starting_elements" {
			currentTierNum = 0
		}

		if currentTierNum != -1 {
			node := headingSel.Next()
			foundTableForThisTier := false
			for node.Length() > 0 && !foundTableForThisTier {
				if node.Is("h2, h3") { // stop if hit the next tier heading
					break
				}
				var tableToProcess *goquery.Selection
				if node.Is("table.list-table") {
					tableToProcess = node
				} else if node.Is("div") {
					tableToProcess = node.Find("table.list-table").First()
				}

				if tableToProcess != nil && tableToProcess.Length() > 0 {
					foundTableForThisTier = true
					tableToProcess.Find("tbody tr").Each(func(i int, rowSelection *goquery.Selection) {
						tds := rowSelection.Find("td")
						if tds.Length() < 1 {
							return
						}

						// attempt to get element name from the first link
						resultElementNode := tds.Eq(0).Find("a[href^='/wiki/']").First()
						resultElementName := strings.TrimSpace(resultElementNode.Text())

						// fallback if no link or link text is empty
						if resultElementName == "" {
							clonedTd := tds.Eq(0).Clone()
							clonedTd.Find("span, script, style, sup, br, img, a.image, div.gallerytext").Remove()
							resultElementName = strings.TrimSpace(clonedTd.Text())
							if strings.Contains(resultElementName, "\n") {
								resultElementName = strings.TrimSpace(strings.Split(resultElementName, "\n")[0])
							}
						}
						resultElementName = normalizeElementName(resultElementName)

						if resultElementName != "" {
							isFilteredOut := false
							tds.Eq(0).Find("a[title*='Myths and Monsters Pack'], img[alt*='Myths and Monsters Pack'], span[class*='pack-icon']").EachWithBreak(func(_ int, s *goquery.Selection) bool {
								if !isStandardBaseElement(resultElementName) {
									isFilteredOut = true
									return false
								}
								return true
							})
							if isFilteredOut {
								return
							}
							if tds.Length() > 1 {
								if tds.Eq(1).Find("a[href='/wiki/Elements_(Myths_and_Monsters)']").Length() > 0 {
									isFilteredOut = true
								}
							}
							if isFilteredOut {
								return
							}

							// assign tier if not already assigned or if previous assignment was placeholder (-1)
							if _, exists := elementTiers[resultElementName]; !exists || elementTiers[resultElementName] == -1 {
								elementTiers[resultElementName] = currentTierNum
							}
						}
					})
				}
				node = node.Next()
			}
		}
	})
	return elementTiers
}

// extractRawRecipesAndElements scrapes recipes and a list of valid element names from the document
func extractRawRecipesAndElements(doc *goquery.Document) ([]RawRecipeEntry, map[string]bool, []string) {
	var rawRecipes []RawRecipeEntry
	validResultElementsFromPage := make(map[string]bool)
	var tempOrderedResultNames []string
	isNameInTempOrderedList := make(map[string]bool)

	doc.Find("div.mw-parser-output table.list-table tbody tr").Each(func(i int, rowSelection *goquery.Selection) {
		tds := rowSelection.Find("td")
		if tds.Length() < 1 {
			return
		}

		// extract result element name (similar logic to tier extraction)
		resultElementNode := tds.Eq(0).Find("a[href^='/wiki/']").First()
		resultElementName := strings.TrimSpace(resultElementNode.Text())
		if resultElementName == "" {
			clonedTd := tds.Eq(0).Clone()
			clonedTd.Find("span, script, style, sup, br, img, a.image, div.gallerytext").Remove()
			resultElementName = strings.TrimSpace(clonedTd.Text())
			if strings.Contains(resultElementName, "\n") {
				resultElementName = strings.TrimSpace(strings.Split(resultElementName, "\n")[0])
			}
		}
		if resultElementName == "" {
			return
		}
		resultElementName = normalizeElementName(resultElementName)

		// filter out elements from specific packs
		if tds.Length() > 1 {
			if tds.Eq(1).Find("a[href='/wiki/Elements_(Myths_and_Monsters)']").Length() > 0 {
				return
			}
		}
		var isVisuallyMarkedAsPack = false
		tds.Eq(0).Find("a[title*='Myths and Monsters Pack'], img[alt*='Myths and Monsters Pack'], span[class*='pack-icon']").Each(func(_ int, s *goquery.Selection) {
			if !isStandardBaseElement(resultElementName) {
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
				if ingredientLinks.Length() >= 2 { // a recipe needs at least two ingredients
					ing1 := normalizeElementName(ingredientLinks.Eq(0).Text())
					ing2 := normalizeElementName(ingredientLinks.Eq(1).Text())
					if ing1 != "" && ing2 != "" {
						sortedIng := []string{ing1, ing2}
						sort.Strings(sortedIng)
						rawRecipes = append(rawRecipes, RawRecipeEntry{resultElementName, sortedIng[0], sortedIng[1]})
					}
				}
			})
		}
	})

	// filter rawRecipes to ensure their result elements are considered valid
	var finalRawRecipes []RawRecipeEntry
	for _, rr := range rawRecipes {
		if validResultElementsFromPage[rr.ResultElement] {
			finalRawRecipes = append(finalRawRecipes, rr)
		}
	}
	return finalRawRecipes, validResultElementsFromPage, tempOrderedResultNames
}

// buildRecipeTreeNodes constructs an initial graph-like structure of elements and their recipes
func buildRecipeTreeNodes(rawRecipes []RawRecipeEntry, allElementNamesOnPage []string) map[string]*model.RecipeTreeNodeTier {
	allNodes := make(map[string]*model.RecipeTreeNodeTier)

	// ensure all unique element names from page (and base elements) have a node
	elementNameSet := make(map[string]bool)
	for _, name := range allElementNamesOnPage {
		if _, exists := allNodes[name]; !exists {
			allNodes[name] = &model.RecipeTreeNodeTier{NamaElemen: name}
		}
		elementNameSet[name] = true
	}
	for _, baseName := range standardBaseElements {
		normBaseName := normalizeElementName(baseName)
		if _, exists := allNodes[normBaseName]; !exists {
			allNodes[normBaseName] = &model.RecipeTreeNodeTier{NamaElemen: normBaseName}
		}
		elementNameSet[normBaseName] = true
	}

	// populate the DibuatDari field for each node
	for _, rr := range rawRecipes {
		resNode, rExists := allNodes[rr.ResultElement]
		ing1Node, i1Exists := allNodes[rr.Ingredient1]
		ing2Node, i2Exists := allNodes[rr.Ingredient2]

		if rExists && i1Exists && i2Exists {
			if isStandardBaseElement(resNode.NamaElemen) { // base elements are not made from anything
				continue
			}
			// check if this combination already exists to avoid duplicates
			comboExists := false
			for _, combo := range resNode.DibuatDari {
				if (combo[0].NamaElemen == rr.Ingredient1 && combo[1].NamaElemen == rr.Ingredient2) ||
					(combo[0].NamaElemen == rr.Ingredient2 && combo[1].NamaElemen == rr.Ingredient1) {
					comboExists = true
					break
				}
			}
			if !comboExists {
				resNode.DibuatDari = append(resNode.DibuatDari, [2]*model.RecipeTreeNodeTier{ing1Node, ing2Node})
			}
		}
	}
	return allNodes
}

// assignIDsAndCreateElements sorts elements, assigns them unique IDs, and creates the final Element structs
func assignIDsAndCreateElements(allNodes map[string]*model.RecipeTreeNodeTier, elementTiersFromWeb map[string]int, allElementNamesFromExtraction []string) ([]*Element, map[string]int, map[int]*Element) {
	type elementInfoForIDSorter struct {
		Name string
		Tier int
	}

	var elementsForIDSorting []elementInfoForIDSorter

	uniqueElementNames := make(map[string]bool)
	for _, name := range allElementNamesFromExtraction {
		uniqueElementNames[name] = true
	}

	for _, baseName := range standardBaseElements {
		uniqueElementNames[normalizeElementName(baseName)] = true
	}

	for name := range uniqueElementNames {
		currentTier := -1 // default to unassigned tier
		if isStandardBaseElement(name) {
			currentTier = 0
		} else if tier, found := elementTiersFromWeb[name]; found {
			currentTier = tier
		}
		elementsForIDSorting = append(elementsForIDSorting, elementInfoForIDSorter{Name: name, Tier: currentTier})
	}

	sort.Slice(elementsForIDSorting, func(i, j int) bool {
		elI := elementsForIDSorting[i]
		elJ := elementsForIDSorting[j]
		if elI.Tier != elJ.Tier {
			if elI.Tier == -1 {
				return false
			}
			if elJ.Tier == -1 {
				return true
			}
			return elI.Tier < elJ.Tier
		}
		return elI.Name < elJ.Name
	})

	nameToID := make(map[string]int)
	elementMapByID := make(map[int]*Element)
	var orderedElements []*Element
	currentIDCounter := 0

	for _, sortedElInfo := range elementsForIDSorting {
		el := &Element{
			ID:       currentIDCounter,
			Name:     sortedElInfo.Name,
			Tier:     sortedElInfo.Tier,
			FromPair: make([][]int, 0),
			CanMake:  make([]int, 0),
		}
		nameToID[el.Name] = el.ID
		elementMapByID[el.ID] = el
		orderedElements = append(orderedElements, el)
		currentIDCounter++
	}
	return orderedElements, nameToID, elementMapByID
}

func populateElementRelationships(orderedElements []*Element, allNodes map[string]*model.RecipeTreeNodeTier, nameToID map[string]int, elementMapByID map[int]*Element) {
	canMakeTemp := make(map[int]map[int]bool)

	for _, currentElement := range orderedElements {
		resultNode, nodeExists := allNodes[currentElement.Name]
		if !nodeExists || isStandardBaseElement(currentElement.Name) || resultNode.DibuatDari == nil || len(resultNode.DibuatDari) == 0 {
			continue
		}

		tempValidFromPairs := [][]int{}
		processedIngredientPairKeys := make(map[string]bool)

		for _, recipePairNodes := range resultNode.DibuatDari {
			ing1Node := recipePairNodes[0]
			ing2Node := recipePairNodes[1]
			if ing1Node == nil || ing2Node == nil {
				continue
			}

			ing1ID, ok1 := nameToID[ing1Node.NamaElemen]
			ing2ID, ok2 := nameToID[ing2Node.NamaElemen]

			if ok1 && ok2 {
				currentPairIDs := []int{ing1ID, ing2ID}
				if currentPairIDs[0] > currentPairIDs[1] {
					currentPairIDs[0], currentPairIDs[1] = currentPairIDs[1], currentPairIDs[0]
				}
				pairKey := fmt.Sprintf("%d-%d", currentPairIDs[0], currentPairIDs[1])

				if !processedIngredientPairKeys[pairKey] {
					processedIngredientPairKeys[pairKey] = true
					if isRecipeValidInternal(ing1ID, ing2ID, currentElement.ID, elementMapByID) {
						tempValidFromPairs = append(tempValidFromPairs, currentPairIDs)

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
		}
		currentElement.FromPair = tempValidFromPairs
	}

	for ingID, producesMap := range canMakeTemp {
		if ingElement, ok := elementMapByID[ingID]; ok {
			for prodID := range producesMap {
				ingElement.CanMake = append(ingElement.CanMake, prodID)
			}
			sort.Ints(ingElement.CanMake)
		}
	}
}

func FetchAndProcessData() ([]*Element, error) {
	log.Println("starting scraping process from:", targetURL)

	doc, err := fetchDocument(targetURL)
	if err != nil {
		return nil, err
	}

	elementTiersFromWeb := extractElementTiers(doc)

	rawRecipes, _, tempOrderedResultNames := extractRawRecipesAndElements(doc)

	allNodes := buildRecipeTreeNodes(rawRecipes, tempOrderedResultNames)

	orderedElements, nameToID, elementMapByID := assignIDsAndCreateElements(allNodes, elementTiersFromWeb, tempOrderedResultNames)

	populateElementRelationships(orderedElements, allNodes, nameToID, elementMapByID)

	var finalFilteredElements []*Element
	for _, el := range orderedElements {
		isBase := isStandardBaseElement(el.Name)
		hasValidRecipes := len(el.FromPair) > 0


		if isBase {
			finalFilteredElements = append(finalFilteredElements, el)
		} else if hasValidRecipes && el.Tier != -1 {
			finalFilteredElements = append(finalFilteredElements, el)
		}
	}
	log.Println("all data processing in scraper complete. data will be returned directly.")
	return finalFilteredElements, nil
}

func GetProcessedElements() []*Element {
	log.Println("GetProcessedElements called: will perform a new fetch and process data.")
	elements, err := FetchAndProcessData()
	if err != nil {
		log.Printf("error during FetchAndProcessData called from GetProcessedElements: %v. returning empty slice.", err)
		return []*Element{}
	}
	return elements
}
