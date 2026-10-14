package main

import (
	"fmt"
	"log"

	"id-domain-checker/checker"
)

func main() {
	c := checker.NewListChecker()

	// Muat dataset institusi Indonesia
	datasetPath := "indonesia_institutions.json"
	if err := c.LoadFromFile(datasetPath); err != nil {
		log.Fatalf("Gagal memuat dataset: %v", err)
	}

	fmt.Println("==================================================")
	fmt.Printf(" Berhasil memuat %d domain terdaftar\n", c.TotalRegistered())
	fmt.Println("==================================================")

	queries := []string{
		"ui.ac.id",
		"kampus-antah-berantah.ac.id",
		"sman1jakarta.sch.id",
		"smk-negeri-pilihan.sch.id",
		"https://www.itb.ac.id/", // uji normalisasi URL
		"google.com",            // ekstensi tidak diizinkan
	}

	for _, domain := range queries {
		res, err := c.CheckAvailability(domain)
		if err != nil {
			fmt.Printf("[INVALID]   %-32s -> %v\n", domain, err)
			continue
		}

		if res.IsAvailable {
			fmt.Printf("[AVAILABLE] %-32s -> Belum terdaftar (Bisa digunakan)\n", res.Domain)
		} else {
			fmt.Printf("[TAKEN]     %-32s -> %s (%s)\n",
				res.Domain,
				res.Institution.Name,
				res.Institution.StateProvince,
			)
		}
	}
}
