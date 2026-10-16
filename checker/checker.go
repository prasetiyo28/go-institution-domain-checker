package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// AllowedExtensions mendefinisikan ekstensi institusi pendidikan yang didukung secara default
var AllowedExtensions = map[string]bool{
	"ac.id":  true,
	"sch.id": true,
	"edu":    true,
	"info":   true,
	"ac.is":  true,
	"ac.ic":  true,
}

// ListChecker adalah struktur utama checker berbasis memori
type ListChecker struct {
	mu         sync.RWMutex
	registered map[string]*Institution
}

// New membuat instance kosong baru dari ListChecker
func New() *ListChecker {
	return &ListChecker{
		registered: make(map[string]*Institution),
	}
}

// NewListChecker adalah alias untuk New guna kompatibilitas mundur
func NewListChecker() *ListChecker {
	return New()
}

// Reset membersihkan semua data domain yang terdaftar dalam memori
func (c *ListChecker) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.registered = make(map[string]*Institution)
}

// LoadFromFile memuat dataset dari file JSON lokal
func (c *ListChecker) LoadFromFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("gagal membuka file dataset: %w", err)
	}
	defer file.Close()

	return c.LoadFromReader(file)
}

// LoadFromReader membaca dataset JSON dari io.Reader dan menggabungkannya ke dalam memori
func (c *ListChecker) LoadFromReader(r io.Reader) error {
	var items []Institution
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&items); err != nil {
		return fmt.Errorf("format JSON tidak cocok: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.registered == nil {
		c.registered = make(map[string]*Institution, len(items))
	}

	for i := range items {
		inst := &items[i]
		for _, rawDomain := range inst.Domains {
			clean := NormalizeDomain(rawDomain)
			if clean != "" {
				c.registered[clean] = inst
			}
		}
	}

	return nil
}

// CheckAvailability memeriksa ketersediaan domain institusi
func (c *ListChecker) CheckAvailability(domain string) (*CheckResult, error) {
	cleanDomain := NormalizeDomain(domain)
	if cleanDomain == "" {
		return nil, ErrInvalidDomainFormat
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	// 1. Cek kecocokan eksak domain
	if inst, exists := c.registered[cleanDomain]; exists {
		ext, _ := GetValidExtension(cleanDomain)
		return &CheckResult{
			Domain:      cleanDomain,
			Extension:   ext,
			IsAvailable: false,
			Institution: inst,
		}, nil
	}

	// 2. Cek apakah domain merupakan subdomain dari institusi yang sudah terdaftar
	// Contoh: 'mail.ugm.ac.id' -> periksa root 'ugm.ac.id'
	parts := strings.Split(cleanDomain, ".")
	for i := 1; i < len(parts)-1; i++ {
		candidate := strings.Join(parts[i:], ".")
		if inst, exists := c.registered[candidate]; exists {
			ext, _ := GetValidExtension(cleanDomain)
			return &CheckResult{
				Domain:      cleanDomain,
				Extension:   ext,
				IsAvailable: false,
				Institution: inst,
			}, nil
		}
	}

	// 3. Validasi akhiran ekstensi domain institusi
	ext, ok := GetValidExtension(cleanDomain)
	if !ok {
		return nil, fmt.Errorf("%w: domain '%s' bukan domain institusi pendidikan yang valid (.ac.id / .sch.id / .edu)", ErrUnsupportedExtension, cleanDomain)
	}

	// Domain dengan ekstensi valid dan belum terdaftar
	return &CheckResult{
		Domain:      cleanDomain,
		Extension:   ext,
		IsAvailable: true,
		Institution: nil,
	}, nil
}

// TotalRegistered mengembalikan jumlah domain unik yang terindeks dalam memori
func (c *ListChecker) TotalRegistered() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.registered)
}

// NormalizeDomain membersihkan protokol, subpath, prefix www, dan spasi
func NormalizeDomain(domain string) string {
	d := strings.ToLower(strings.TrimSpace(domain))
	d = strings.TrimPrefix(d, "https://")
	d = strings.TrimPrefix(d, "http://")
	d = strings.TrimPrefix(d, "www.")
	d = strings.TrimRight(d, "/")
	return d
}

// GetValidExtension memvalidasi akhiran domain dan mengembalikan ekstensi yang cocok
func GetValidExtension(domain string) (string, bool) {
	// Urutkan pengecekan ekstensi yang lebih panjang terlebih dahulu (misal 'ac.id' sebelum 'id')
	extensions := []string{"ac.id", "sch.id", "ac.is", "ac.ic", "edu", "info"}
	for _, ext := range extensions {
		if strings.HasSuffix(domain, "."+ext) || domain == ext {
			return ext, true
		}
	}

	for ext := range AllowedExtensions {
		if strings.HasSuffix(domain, "."+ext) || domain == ext {
			return ext, true
		}
	}
	return "", false
}
