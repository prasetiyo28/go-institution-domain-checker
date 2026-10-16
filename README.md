# go-institution-domain-checker

[![Go Reference](https://pkg.go.dev/badge/github.com/prasetiyo28/go-institution-domain-checker/checker.svg)](https://pkg.go.dev/github.com/prasetiyo28/go-institution-domain-checker/checker)
[![Go Report Card](https://goreportcard.com/badge/github.com/prasetiyo28/go-institution-domain-checker)](https://goreportcard.com/report/github.com/prasetiyo28/go-institution-domain-checker)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Library Go untuk memeriksa ketersediaan (*availability*) domain universitas dan sekolah di Indonesia (`.ac.id`, `.sch.id`, `.edu`, `.info`) berbasis data lokal in-memory tanpa perlu melakukan *hit* ke soket WHOIS eksternal atau server PANDI.

---

## ✨ Fitur Utama

- 📦 **Zero-Configuration (`//go:embed`):** Dataset bawaan institusi pendidikan Indonesia (245+ kampus, 350+ domain) sudah ter-bundle langsung ke dalam pustaka Go. Siap pakai tanpa perlu menyalin file JSON.
- ⚙️ **Dukungan Dataset Tambahan (`OTHER_INSTITUTIONS`):** Tambahkan institusi baru secara opsional lewat Environment Variable (file path atau raw JSON string) atau via kode.
- ⚡ **Pencarian $O(1)$ Instan:** Menggunakan pemetaan `map[string]*Institution` di memori dengan fallback pencarian subdomain bertingkat (*hierarchical subdomain matching*).
- 🚫 **Bebas Rate-Limit:** Berjalan 100% lokal di memori, tanpa ketergantungan koneksi jaringan atau batasan rate-limit WHOIS.
- 🔒 **Thread-Safe:** Aman digunakan pada konkurensi goroutine tinggi (`sync.RWMutex`).
- 🧹 **Sanitasi Otomatis:** Otomatis membersihkan protokol (`http://`, `https://`), prefiks `www.`, subpath trailing, dan normalisasi huruf kecil.

---

## 📥 Instalasi

Tambahkan modul ini ke proyek Go Anda:

```bash
go get github.com/prasetiyo28/go-institution-domain-checker/checker
```

---

## 🚀 Cara Penggunaan

### 1. Penggunaan Standar (Zero-Configuration)

Cukup panggil `checker.NewDefault()` untuk menggunakan dataset resmi bawaan:

```go
package main

import (
	"fmt"
	"log"

	"github.com/prasetiyo28/go-institution-domain-checker/checker"
)

func main() {
	// Inisialisasi checker dengan dataset bawaan
	c, err := checker.NewDefault()
	if err != nil {
		log.Fatalf("Gagal inisialisasi checker: %v", err)
	}

	// Periksa domain (mendukung root domain maupun subdomain)
	res, err := c.CheckAvailability("itb.ac.id")
	if err != nil {
		log.Fatalf("Domain tidak valid: %v", err)
	}

	if res.IsAvailable {
		fmt.Printf("✅ Domain %s belum terdaftar (Bisa digunakan)\n", res.Domain)
	} else {
		fmt.Printf("❌ Domain %s sudah terdaftar untuk %s (%s)\n",
			res.Domain,
			res.Institution.Name,
			res.Institution.StateProvince,
		)
	}
}
```

---

## 🧩 Menambahkan Institusi Tambahan (`OTHER_INSTITUTIONS`)

Anda dapat memperluas atau menambahkan institusi secara opsional tanpa perlu mengubah kode library.

### Format JSON yang Didukung

Struktur data harus berupa **JSON Array of Objects**:

#### Format Minimal (Hanya Field Wajib)
```json
[
  {
    "name": "Institut Teknologi Kustom",
    "domains": [
      "itkustom.ac.id"
    ]
  },
  {
    "name": "SMK Unggulan Baru",
    "domains": [
      "smkunggulan.sch.id"
    ]
  }
]
```

#### Format Lengkap (Dengan Metadata Opsional)
```json
[
  {
    "name": "Universitas Mandiri Nusantara",
    "country": "Indonesia",
    "domains": [
      "umn-contoh.ac.id",
      "kampus.umn-contoh.ac.id"
    ],
    "web_pages": [
      "https://www.umn-contoh.ac.id/"
    ],
    "alpha_two_code": "ID",
    "state-province": "DKI Jakarta"
  }
]
```

---

### Cara Mengaktifkan

#### Opsi A: Menggunakan Environment Variable (Rekomendasi)

Setel `OTHER_INSTITUTIONS` pada environment server / file `.env`:

```bash
# 1. Mengarahkan ke file JSON lokal
export OTHER_INSTITUTIONS="/app/config/custom_institutions.json"

# 2. Atau langsung berupa raw JSON string
export OTHER_INSTITUTIONS='[{"name":"Poltek Swasta","domains":["poltekswasta.ac.id"]}]'
```

Di dalam kode Go, `checker.NewDefault()` akan otomatis mendeteksi dan menggabungkan data tersebut:
```go
c, err := checker.NewDefault()
```

#### Opsi B: Menggunakan Functional Option via Kode Go

```go
c, err := checker.NewDefault(
    checker.WithAdditionalFile("./my_institutions.json"),
)
```

Atau menggunakan `io.Reader`:
```go
c, err := checker.NewDefault(
    checker.WithAdditionalReader(myReader),
)
```

---

## 💻 Menjalankan Contoh CLI

Repositori ini menyertakan command-line interface (CLI) bawaan untuk pengujian langsung:

```bash
# Menjalankan domain demo bawaan
go run cmd/checker/main.go

# Memeriksa domain tertentu via argumen CLI
go run cmd/checker/main.go ui.ac.id mail.ugm.ac.id kampus-baru.ac.id
```

Output:
```text
==================================================
 Berhasil memuat 354 domain terdaftar
==================================================
[TAKEN]     ui.ac.id                         -> Universitas Indonesia (Jawa Barat)
[TAKEN]     mail.ugm.ac.id                   -> Universitas Gadjah Mada (DI Yogyakarta)
[AVAILABLE] kampus-baru.ac.id                -> Belum terdaftar (Bisa digunakan)
```

---

## 🧪 Menjalankan Pengujian (Unit Tests)

```bash
go test -v ./...
```

Untuk melihat cakupan tes (*coverage*):
```bash
go test -cover ./...
```

---

## 📄 Lisensi

Proyek ini dilisensikan di bawah lisensi [MIT](LICENSE).
