# id-domain-checker

Library Go untuk memeriksa ketersediaan (*availability*) domain universitas (`.ac.id`) dan sekolah (`.sch.id`) di Indonesia berdasarkan dataset lokal JSON tanpa perlu *hit* soket WHOIS eksternal.

## ✨ Fitur

- ⚡ **Pencarian $O(1)$ Instan:** Menggunakan pemetaan `map[string]*Institution` di memori.
- 🚫 **Bebas Rate-Limit:** Tanpa melakukan koneksi jaringan atau *hit* server WHOIS PANDI.
- 🔒 **Thread-Safe:** Dilengkapi `sync.RWMutex` sehingga aman digunakan dalam konkurensi goroutine / web server.
- 🧹 **Pembersihan Otomatis (Sanitization):** Otomatis membersihkan protokol (`http://`, `https://`), `www.`, path trailing, dan huruf besar/kecil.
- 🔍 **Detail Institusi:** Jika domain sudah terdaftar, mengembalikan informasi lengkap seperti nama institusi, provinsi, negara, dan situs resmi.
- 🧪 **Unit Test Komprehensif:** Dilengkapi dengan test suite bawaan Go (`go test`).

---

## 📁 Struktur Direktori

```text
id-domain-checker/
├── checker/
│   ├── checker.go           # Logika core library
│   └── checker_test.go      # Unit test
├── indonesia_institutions.json # Dataset JSON universitas & sekolah Indonesia
├── go.mod                   # Definisi Go module
├── main.go                  # Contoh implementasi CLI
└── README.md                # Dokumentasi
```

---

## 🚀 Memulai Cepat

### 1. Menjalankan Contoh Program

```bash
go run main.go
```

Output:
```text
==================================================
 Berhasil memuat 58 domain terdaftar
==================================================
[TAKEN]     ui.ac.id                         -> Universitas Indonesia (Jawa Barat)
[AVAILABLE] kampus-antah-berantah.ac.id      -> Belum terdaftar (Bisa digunakan)
[TAKEN]     sman1jakarta.sch.id              -> SMA Negeri 1 Jakarta (DKI Jakarta)
[AVAILABLE] smk-negeri-pilihan.sch.id        -> Belum terdaftar (Bisa digunakan)
[TAKEN]     itb.ac.id                        -> Institut Teknologi Bandung (Jawa Barat)
[INVALID]   google.com                       -> domain 'google.com' bukan domain institusi pendidikan Indonesia yang valid (.ac.id / .sch.id)
```

### 2. Menjalankan Unit Test

Jalankan pengujian dari folder project:

```bash
go test -v ./checker/...
```

Untuk melihat cakupan kode (*test coverage*):

```bash
go test -cover ./checker/...
```

---

## 🛠️ Contoh Kode Penggunaan

```go
package main

import (
	"fmt"
	"log"

	"id-domain-checker/checker"
)

func main() {
	c := checker.NewListChecker()

	// Muat dataset dari file JSON
	if err := c.LoadFromFile("indonesia_institutions.json"); err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Periksa ketersediaan domain
	res, err := c.CheckAvailability("ui.ac.id")
	if err != nil {
		log.Println("Domain tidak valid:", err)
		return
	}

	if res.IsAvailable {
		fmt.Printf("Domain %s tersedia!\n", res.Domain)
	} else {
		fmt.Printf("Domain %s sudah dipakai oleh %s\n", res.Domain, res.Institution.Name)
	}
}
```

---

## 📊 Skema JSON Dataset

Dataset menggunakan skema *array of objects*:

```json
[
  {
    "name": "Universitas Gadjah Mada",
    "country": "Indonesia",
    "domains": ["ugm.ac.id"],
    "web_pages": ["https://www.ugm.ac.id/"],
    "alpha_two_code": "ID",
    "state-province": "DI Yogyakarta"
  }
]
```

Lisensi: MIT
