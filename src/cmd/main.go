// src/cmd/main.go
package main

import (
	"fmt"
	"log"
	"net/http" // package untuk fungsionalitas HTTP

	// "github.com/graceevelyns/Tubes2_BE_ian/internal/api" // kalo sudah punya router
)

func main() {
	fmt.Println("Memulai Backend Server Little Alchemy 2 Solver...")

	// router := api.NewRouter() // kalo menggunakan router kustom dari package api

	// for now ini handler default atau handler sederhana
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Selamat datang di Backend Little Alchemy 2 Solver!")
	})

	// port didefinisikan dan server mulai mendengarkan:
	port := ":8080" // angka port diawali dengan titik dua ":"
	                // server akan mendengarkan di semua antarmuka jaringan yang tersedia
	                // pada mesin ini di port 8080.

	fmt.Printf("Server akan berjalan di alamat http://localhost%s\n", port)

	// http.ListenAndServe memulai server HTTP dengan alamat dan handler yang diberikan.
	// parameter kedua adalah handler; 'nil' berarti menggunakan DefaultServeMux dari package http,
	// yang akan menggunakan handler yang telah kita daftarkan dengan http.HandleFunc.
	// kalo router kustom (seperti router dari package api), masukkin di sini.
	// contoh: log.Fatal(http.ListenAndServe(port, router))
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Gagal memulai server: %s\n", err)
	}
}