package checker

import (
	"errors"
	"os"
	"path/filepath"
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
	c := New()
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

func TestCheckAvailability_SubdomainMatching(t *testing.T) {
	c := setupTestChecker(t)

	// ui.ac.id terdaftar, periksa apakah subdomain 'alumni.ui.ac.id' atau 'sci.ui.ac.id' dikenali
	subdomains := []string{
		"alumni.ui.ac.id",
		"sci.ui.ac.id",
		"mhs.usk.ac.id",
	}

	for _, sub := range subdomains {
		res, err := c.CheckAvailability(sub)
		if err != nil {
			t.Errorf("CheckAvailability(%s) gagal: %v", sub, err)
			continue
		}
		if res.IsAvailable {
			t.Errorf("CheckAvailability(%s) harusnya terdaftar via parent domain", sub)
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
			t.Errorf("CheckAvailability(%s) diharapkan gagal karena ekstensi bukan institusi yang didukung", domain)
		}
		if !errors.Is(err, ErrUnsupportedExtension) {
			t.Errorf("error untuk %s harus membungkus ErrUnsupportedExtension, dapat: %v", domain, err)
		}
	}
}

func TestCheckAvailability_EmptyDomain(t *testing.T) {
	c := setupTestChecker(t)
	_, err := c.CheckAvailability("   ")
	if !errors.Is(err, ErrInvalidDomainFormat) {
		t.Errorf("expected ErrInvalidDomainFormat, got: %v", err)
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

func TestNewDefault(t *testing.T) {
	// Pastikan env dibersihkan selama tes ini
	origEnv := os.Getenv(EnvOtherInstitutions)
	os.Unsetenv(EnvOtherInstitutions)
	defer func() {
		if origEnv != "" {
			os.Setenv(EnvOtherInstitutions, origEnv)
		}
	}()

	c, err := NewDefault()
	if err != nil {
		t.Fatalf("NewDefault() gagal: %v", err)
	}

	// Dataset default harus memuat lebih dari 300 domain
	if c.TotalRegistered() < 300 {
		t.Errorf("TotalRegistered() = %d; ingin >= 300", c.TotalRegistered())
	}

	// Cek beberapa kampus utama dari dataset default
	samples := []string{"itb.ac.id", "ugm.ac.id", "binus.edu", "raharja.info"}
	for _, domain := range samples {
		res, err := c.CheckAvailability(domain)
		if err != nil {
			t.Errorf("CheckAvailability(%s) error: %v", domain, err)
			continue
		}
		if res.IsAvailable {
			t.Errorf("domain default %s harusnya sudah terdaftar", domain)
		}
	}
}

func TestNewDefault_WithEnvRawJSON(t *testing.T) {
	customJSON := `[
		{
			"name": "Universitas Kustom Testing Env",
			"domains": ["kampus-env-test.ac.id"]
		}
	]`
	os.Setenv(EnvOtherInstitutions, customJSON)
	defer os.Unsetenv(EnvOtherInstitutions)

	c, err := NewDefault()
	if err != nil {
		t.Fatalf("NewDefault() dengan env raw JSON gagal: %v", err)
	}

	res, err := c.CheckAvailability("kampus-env-test.ac.id")
	if err != nil {
		t.Fatalf("CheckAvailability error: %v", err)
	}
	if res.IsAvailable {
		t.Errorf("kampus-env-test.ac.id seharusnya terdaftar dari env raw JSON")
	}
	if res.Institution.Name != "Universitas Kustom Testing Env" {
		t.Errorf("Nama institusi = %s; ingin 'Universitas Kustom Testing Env'", res.Institution.Name)
	}
}

func TestNewDefault_WithEnvFilePath(t *testing.T) {
	customJSON := `[
		{
			"name": "Institut Kustom Dari File",
			"domains": ["institut-file-test.ac.id"]
		}
	]`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "custom.json")
	if err := os.WriteFile(filePath, []byte(customJSON), 0644); err != nil {
		t.Fatalf("gagal membuat file temp: %v", err)
	}

	os.Setenv(EnvOtherInstitutions, filePath)
	defer os.Unsetenv(EnvOtherInstitutions)

	c, err := NewDefault()
	if err != nil {
		t.Fatalf("NewDefault() dengan env file path gagal: %v", err)
	}

	res, err := c.CheckAvailability("institut-file-test.ac.id")
	if err != nil {
		t.Fatalf("CheckAvailability error: %v", err)
	}
	if res.IsAvailable {
		t.Errorf("institut-file-test.ac.id seharusnya terdaftar dari file")
	}
	if res.Institution.Name != "Institut Kustom Dari File" {
		t.Errorf("Nama institusi = %s; ingin 'Institut Kustom Dari File'", res.Institution.Name)
	}
}

func TestWithAdditionalReader(t *testing.T) {
	c, err := NewDefault(
		WithAdditionalReader(strings.NewReader(`[{"name":"Programmatic Univ","domains":["prog-univ.ac.id"]}]`)),
	)
	if err != nil {
		t.Fatalf("NewDefault with Option gagal: %v", err)
	}

	res, err := c.CheckAvailability("prog-univ.ac.id")
	if err != nil {
		t.Fatalf("CheckAvailability error: %v", err)
	}
	if res.IsAvailable {
		t.Errorf("prog-univ.ac.id seharusnya terdaftar")
	}
}
