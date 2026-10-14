package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// AllowedExtensions mendefinisikan ekstensi institusi pendidikan yang didukung
var AllowedExtensions = map[string]bool{
	"ac.id":  true,
	"sch.id": true,
}

// Institution mencerminkan skema data institusi pendidikan
type Institution struct {
	Name          string   `json:"name"`
	Country       string   `json:"country"`
	Domains       []string `json:"domains"`
	WebPages      []string `json:"web_pages"`
	AlphaTwoCode  string   `json:"alpha_two_code"`
	StateProvince string   `json:"state-province"`
}

// CheckResult berisi hasil pengecekan domain
type CheckResult struct {
	Domain      string       `json:"domain"`
	Extension   string       `json:"extension"`
	IsAvailable bool         `json:"is_available"`
	Institution *Institution `json:"institution,omitempty"`
}

// ListChecker adalah struktur utama checker berbasis memori
type ListChecker struct {
	mu         sync.RWMutex
	registered map[string]*Institution
}

// NewListChecker membuat instance baru dari ListChecker
func NewListChecker() *ListChecker {
	return &ListChecker{
		registered: make(map[string]*Institution),
	}
}

// LoadFromFile memuat dataset dari file JSON
func (c *ListChecker) LoadFromFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("gagal membuka file dataset: %w", err)
	}
	defer file.Close()

	return c.LoadFromReader(file)
}

// LoadFromReader membaca dataset JSON dari io.Reader
func (c *ListChecker) LoadFromReader(r io.Reader) error {
	var items []Institution
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&items); err != nil {
		return fmt.Errorf("format JSON tidak cocok: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.registered = make(map[string]*Institution, len(items))

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

	ext, ok := GetValidExtension(cleanDomain)
	if !ok {
		return nil, fmt.Errorf("domain '%s' bukan domain institusi pendidikan Indonesia yang valid (.ac.id / .sch.id)", cleanDomain)
	}

	c.mu.RLock()
	inst, exists := c.registered[cleanDomain]
	c.mu.RUnlock()

	if exists {
		return &CheckResult{
			Domain:      cleanDomain,
			Extension:   ext,
			IsAvailable: false,
			Institution: inst,
		}, nil
	}

	return &CheckResult{
		Domain:      cleanDomain,
		Extension:   ext,
		IsAvailable: true,
		Institution: nil,
	}, nil
}

// TotalRegistered mengembalikan jumlah domain yang terindeks
func (c *ListChecker) TotalRegistered() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.registered)
}

// NormalizeDomain membersihkan prefix URL / www dan whitespace
func NormalizeDomain(domain string) string {
	d := strings.ToLower(strings.TrimSpace(domain))
	d = strings.TrimPrefix(d, "https://")
	d = strings.TrimPrefix(d, "http://")
	d = strings.TrimPrefix(d, "www.")
	d = strings.TrimRight(d, "/")
	return d
}

// GetValidExtension memvalidasi akhiran domain
func GetValidExtension(domain string) (string, bool) {
	for ext := range AllowedExtensions {
		if strings.HasSuffix(domain, "."+ext) {
			return ext, true
		}
	}
	return "", false
}
