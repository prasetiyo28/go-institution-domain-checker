package checker

import (
	"strings"
	"testing"
)

const sampleJSON = `[
  {
    "name": "Universitas Indonesia",
    "country": "Indonesia",
    "domains": ["ui.ac.id"],
    "web_pages": ["https://www.ui.ac.id/"],
    "alpha_two_code": "ID",
    "state-province": "Jawa Barat"
  },
  {
    "name": "Universitas Syiah Kuala",
    "country": "Indonesia",
    "domains": ["usk.ac.id", "unsyiah.ac.id"],
    "web_pages": ["https://www.usk.ac.id/"],
    "alpha_two_code": "ID",
    "state-province": "Aceh"
  },
  {
    "name": "SMA Negeri 1 Jakarta",
    "country": "Indonesia",
    "domains": ["sman1jakarta.sch.id"],
    "web_pages": ["http://sman1jakarta.sch.id/"],
    "alpha_two_code": "ID",
    "state-province": "DKI Jakarta"
  }
]`

func setupTestChecker(t *testing.T) *ListChecker {
	c := NewListChecker()
	reader := strings.NewReader(sampleJSON)
	if err := c.LoadFromReader(reader); err != nil {
		t.Fatalf("gagal inisialisasi dataset pengujian: %v", err)
	}
	return c
}

func TestLoadFromReader(t *testing.T) {
	c := setupTestChecker(t)

	// Total domain: ui.ac.id (1) + usk.ac.id, unsyiah.ac.id (2) + sman1jakarta.sch.id (1) = 4
	expectedCount := 4
	if got := c.TotalRegistered(); got != expectedCount {
		t.Errorf("TotalRegistered() = %d; ingin %d", got, expectedCount)
	}
}

func TestCheckAvailability_TakenDomains(t *testing.T) {
	c := setupTestChecker(t)

	testCases := []struct {
		domain       string
		expectedName string
	}{
		{"ui.ac.id", "Universitas Indonesia"},
		{"https://UI.AC.ID/", "Universitas Indonesia"}, // Memastikan normalisasi URL berjalan
		{"usk.ac.id", "Universitas Syiah Kuala"},
		{"unsyiah.ac.id", "Universitas Syiah Kuala"}, // Domain alias
		{"sman1jakarta.sch.id", "SMA Negeri 1 Jakarta"},
	}

	for _, tc := range testCases {
		res, err := c.CheckAvailability(tc.domain)
		if err != nil {
			t.Errorf("CheckAvailability(%s) error tidak terduga: %v", tc.domain, err)
			continue
		}

		if res.IsAvailable {
			t.Errorf("CheckAvailability(%s) = true; ingin false (harus sudah terdaftar)", tc.domain)
		}

		if res.Institution == nil || res.Institution.Name != tc.expectedName {
			t.Errorf("CheckAvailability(%s) nama institusi = %v; ingin %s", tc.domain, res.Institution, tc.expectedName)
		}
	}
}

func TestCheckAvailability_AvailableDomains(t *testing.T) {
	c := setupTestChecker(t)

	availableDomains := []string{
		"kampus-fiktif-1234.ac.id",
		"smk-negeri-baru.sch.id",
		"universitas-nusantara-indah.ac.id",
	}

	for _, domain := range availableDomains {
		res, err := c.CheckAvailability(domain)
		if err != nil {
			t.Errorf("CheckAvailability(%s) error tidak terduga: %v", domain, err)
			continue
		}

		if !res.IsAvailable {
			t.Errorf("CheckAvailability(%s) = false; ingin true (tersedia)", domain)
		}

		if res.Institution != nil {
			t.Errorf("CheckAvailability(%s) institusi harus nil untuk domain yang tersedia", domain)
		}
	}
}

func TestCheckAvailability_InvalidExtension(t *testing.T) {
	c := setupTestChecker(t)

	invalidDomains := []string{
		"google.com",
		"kemdikbud.go.id",
		"kompas.id",
		"kampus.net",
	}

	for _, domain := range invalidDomains {
		_, err := c.CheckAvailability(domain)
		if err == nil {
			t.Errorf("CheckAvailability(%s) diharapkan gagal karena ekstensi bukan .ac.id atau .sch.id", domain)
		}
	}
}

func TestNormalizeDomain(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"https://www.UI.AC.ID/", "ui.ac.id"},
		{"http://sman1jakarta.sch.id", "sman1jakarta.sch.id"},
		{"  UNAIR.AC.ID  ", "unair.ac.id"},
		{"www.itb.ac.id", "itb.ac.id"},
	}

	for _, tc := range cases {
		got := NormalizeDomain(tc.input)
		if got != tc.expected {
			t.Errorf("NormalizeDomain(%q) = %q; ingin %q", tc.input, got, tc.expected)
		}
	}
}
