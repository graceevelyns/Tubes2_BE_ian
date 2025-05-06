package model

type Element struct {
    Name string
    IsBaseElement bool
    // mungkin perlu field lain seperti 'id' jika diperlukan untuk graf?
}

type Recipe struct {
    Result string // nama elemen hasil
    Ingredients []string // daftar nama elemen input (harusnya 2)
}