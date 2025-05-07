package scraper

// silahkan dibaca semoga tidak bingung :D

import (
	"fmt"      // buat cetak teks ke konsol
	"log"      // buat logging
	"net/http" // buat mengirim HTTP request dan menerima response
	"sort"     // buat mengurutkan slice string
	"strings"  // buat manipulasi string

	"github.com/PuerkitoBio/goquery"                           // buat parsing HTML, ini dari spek
	"github.com/graceevelyns/Tubes2_BE_ian/src/internal/model" // model elemen dan resep
)

const targetURL = "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)" // web yang di-scrape

// RawRecipeEntry -> struct sementara untuk menampung data mentah hasil scraping per resep
// nanti diolah lebih lanjut
type RawRecipeEntry struct {
	ResultElement string
	Ingredient1   string
	Ingredient2   string
}

// ScrapedData -> struct untuk hasil akhir dari scraping, termasuk semua node dan elemen dasar
type ScrapedData struct {
	// ini pointer ke RecipeNode yang merupakan elemen dasar, liat di model.go
	AllNodes     map[string]*model.RecipeNode // so kalo mau cari elemen Lava -> AllNodes["Lava"] -> dapet RecipeNode-nya
	BaseElements []*model.RecipeNode          // kalo ini pointer ke RecipeNode yang merupakan elemen dasar
}

var (
	// daftar elemen dasar yang diketahui
	standardBaseElements = []string{"Air", "Earth", "Fire", "Water"}
)

// FetchAndParseData adalah fungsi utama untuk melakukan scraping dan membangun struktur data awal
// fungsi ini mengembalikan peta semua node elemen dan slice dari elemen dasar
func FetchAndParseData() (map[string]*model.RecipeNode, []*model.RecipeNode, error) {
	log.Println("Memulai proses scraping dari:", targetURL)

	res, err := http.Get(targetURL) // mengirim GET request ke URL

	// jika ada error saat mengirim request, log error dan kembalikan nil
	if err != nil {
		return nil, nil, fmt.Errorf("gagal melakukan GET request ke %s: %w", targetURL, err)
	}
	defer res.Body.Close()

	// jika status code bukan 200 (OK), log error dan kembalikan nil
	if res.StatusCode != 200 {
		return nil, nil, fmt.Errorf("status code tidak OK: %d %s", res.StatusCode, res.Status)
	}

	// parse HTML dari response body menggunakan goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal mem-parse HTML: %w", err)
	}

	log.Println("HTML berhasil di-parse. Memulai ekstraksi data...")
	var rawRecipes []RawRecipeEntry

	// scraping data dari halaman wiki
	// dari inspeksi, dapetnya -> elemen dan resepnya ada dalam <tr> di dalam <table> dengan kelas "list-table"
	// elemen hasil ada di <td> pertama, resep ada di <td> kedua dalam <ul><li>

	doc.Find("div.mw-parser-output table.list-table tbody tr").Each(func(i int, rowSelection *goquery.Selection) {
		tds := rowSelection.Find("td")

		if tds.Length() < 1 { // harus ada setidaknya 1 td untuk nama elemen hasil
			return // kalo gaada skip aja
		}

		// ambil nama elemen hasil dari <td> pertama
		resultElementNode := tds.Eq(0).Find("a[href^='/wiki/']").First()
		resultElementName := strings.TrimSpace(resultElementNode.Text())

		if resultElementName == "" { // fallback kalo ga ada <a> atau formatnya berbeda
			resultElementName = strings.TrimSpace(tds.Eq(0).Contents().Not("span, script, style").Text())
			if strings.Contains(resultElementName, "\n") {
				resultElementName = strings.TrimSpace(strings.Split(resultElementName, "\n")[0])
			}
		}

		if resultElementName == "" {
			return // ga bisa proses kalo tidak ada nama elemen hasil
		}

		// cek jika ini elemen dasar dari teks di kolom kedua (jika ada)
		isBase := false
		if tds.Length() > 1 {
			secondTdText := strings.TrimSpace(tds.Eq(1).Contents().Not("ul, script, style").Text())
			if strings.Contains(strings.ToLower(secondTdText), "available from the start") || strings.Contains(strings.ToLower(secondTdText), "this element does not have any recipes") {
				for _, baseName := range standardBaseElements {
					if resultElementName == baseName { // cek apakah nama elemen hasil adalah salah satu dari elemen dasar
						isBase = true
						break
					}
				}
			}

			// ekstrak resep dari <ul><li> di <td> kedua
			tds.Eq(1).Find("ul > li").Each(func(j int, recipeLi *goquery.Selection) {
				ingredientLinks := recipeLi.Find("a[href^='/wiki/']") // untuk setiap <li> ambil semua <a> yang ada di dalamnya

				if ingredientLinks.Length() >= 2 {
					ingredient1 := strings.TrimSpace(ingredientLinks.Eq(0).Text()) // ambil nama bahan dari <a> pertama
					ingredient2 := strings.TrimSpace(ingredientLinks.Eq(1).Text()) // ambil nama bahan dari <a> kedua

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

		if isBase {
			alreadyHasRecipe := false
			for _, rr := range rawRecipes {
				if rr.ResultElement == resultElementName && rr.Ingredient1 != "" {
					alreadyHasRecipe = true
					break
				}
			}
			if !alreadyHasRecipe {
				rawRecipes = append(rawRecipes, RawRecipeEntry{ResultElement: resultElementName})
			}
		}
	})

	log.Printf("Scraping mentah selesai. Ditemukan %d entri resep potensial.\n", len(rawRecipes))

	// saatnya bikin graf dari raw recipenya :)
	allNodes := make(map[string]*model.RecipeNode)
	var baseElementNodes []*model.RecipeNode

	// inisialisasi elemen dasar yang diketahui
	isKnownBase := make(map[string]bool)
	for _, name := range standardBaseElements {
		isKnownBase[name] = true
		if _, exists := allNodes[name]; !exists { // kalo belum ada, buat node baru
			node := &model.RecipeNode{NamaElemen: name, IsBaseElement: true, DibuatDari: make([][2]*model.RecipeNode, 0)}
			allNodes[name] = node
			baseElementNodes = append(baseElementNodes, node)
		}
	}

	// buat node untuk semua elemen yang terlibat (hasil dan bahan)
	for _, rr := range rawRecipes {
		if _, exists := allNodes[rr.ResultElement]; !exists { // kalo belum ada, buat node baru
			allNodes[rr.ResultElement] = &model.RecipeNode{
				NamaElemen:    rr.ResultElement,
				IsBaseElement: isKnownBase[rr.ResultElement],
				DibuatDari:    make([][2]*model.RecipeNode, 0),
			}
		} else if isKnownBase[rr.ResultElement] { // kalo udah ada dan merupakan base, pastikan IsBaseElement true
			allNodes[rr.ResultElement].IsBaseElement = true
		}

		if rr.Ingredient1 != "" && rr.Ingredient2 != "" {
			if _, exists := allNodes[rr.Ingredient1]; !exists {
				allNodes[rr.Ingredient1] = &model.RecipeNode{
					NamaElemen:    rr.Ingredient1,
					IsBaseElement: isKnownBase[rr.Ingredient1],
					DibuatDari:    make([][2]*model.RecipeNode, 0),
				}
			}
			if _, exists := allNodes[rr.Ingredient2]; !exists {
				allNodes[rr.Ingredient2] = &model.RecipeNode{
					NamaElemen:    rr.Ingredient2,
					IsBaseElement: isKnownBase[rr.Ingredient2],
					DibuatDari:    make([][2]*model.RecipeNode, 0),
				}
			}
		}
	}

	// nah kalo udah ada semua node, kita bisa bikin graf dari recipenya ig ...
	// pastiin ga berulang
	for _, rr := range rawRecipes {
		if rr.Ingredient1 != "" && rr.Ingredient2 != "" && !isKnownBase[rr.ResultElement] {
			resultNode, rExists := allNodes[rr.ResultElement]
			ing1Node, i1Exists := allNodes[rr.Ingredient1]
			ing2Node, i2Exists := allNodes[rr.Ingredient2]

			if rExists && i1Exists && i2Exists {
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

	log.Printf("Pembuatan struktur node selesai. Total node unik: %d. Elemen dasar teridentifikasi: %d\n", len(allNodes), len(baseElementNodes))
	return allNodes, baseElementNodes, nil
}
