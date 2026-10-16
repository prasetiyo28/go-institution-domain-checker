package main

import (
	"fmt"
	"log"
	"os"

	"github.com/prasetiyo28/go-institution-domain-checker/checker"
)

func main() {
	// Inisialisasi checker dengan dataset default + OTHER_INSTITUTIONS jika ada
	c, err := checker.NewDefault()
	if err != nil {
		log.Fatalf("Gagal inisialisasi checker: %v", err)
	}

	fmt.Println("==================================================")
	fmt.Printf(" Berhasil memuat %d domain terdaftar\n", c.TotalRegistered())
	if envVal := os.Getenv(checker.EnvOtherInstitutions); envVal != "" {
		fmt.Printf(" [INFO] %s aktif: %s\n", checker.EnvOtherInstitutions, envVal)
	}
	fmt.Println("==================================================")

	queries := []string{
		"ui.ac.id",
		"kampus-antah-berantah.ac.id",
		"sman1jakarta.sch.id",
		"smk-negeri-pilihan.sch.id",
		"https://www.itb.ac.id/", // normalisasi URL
		"mail.ugm.ac.id",          // subdomain checking
		"binus.edu",               // domain institusi .edu
		"raharja.info",            // domain institusi .info
		"google.com",              // ekstensi umum di luar institusi
	}

	// Jika ada argumen domain yang di-passing via CLI
	if len(os.Args) > 1 {
		queries = os.Args[1:]
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
			instName := "Institusi Terdaftar"
			province := "-"
			if res.Institution != nil {
				instName = res.Institution.Name
				if res.Institution.StateProvince != "" {
					province = res.Institution.StateProvince
				}
			}
			fmt.Printf("[TAKEN]     %-32s -> %s (%s)\n", res.Domain, instName, province)
		}
	}
}
