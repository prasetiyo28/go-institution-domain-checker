package main

import (
	"fmt"
	"log"

	"github.com/prasetiyo28/go-institution-domain-checker/checker"
)

func main() {
	// Inisialisasi checker dengan dataset bawaan (zero-configuration)
	c, err := checker.NewDefault()
	if err != nil {
		log.Fatalf("Gagal inisialisasi checker: %v", err)
	}

	domainsToCheck := []string{
		"itb.ac.id",
		"kampus-baru-contoh.ac.id",
		"binus.edu",
	}

	for _, domain := range domainsToCheck {
		res, err := c.CheckAvailability(domain)
		if err != nil {
			log.Printf("Domain %s error: %v", domain, err)
			continue
		}

		if res.IsAvailable {
			fmt.Printf("✅ [TERSEDIA] %s dapat didaftarkan!\n", res.Domain)
		} else {
			fmt.Printf("❌ [TERDAFTAR] %s sudah digunakan oleh %s\n", res.Domain, res.Institution.Name)
		}
	}
}
