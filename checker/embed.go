package checker

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"strings"
)

//go:embed data/indonesia_institutions.json
var defaultDatasetJSON []byte

// EnvOtherInstitutions adalah nama environment variable untuk dataset kustom tambahan (opsional)
const EnvOtherInstitutions = "OTHER_INSTITUTIONS"

// Option mendefinisikan functional option untuk mengonfigurasi ListChecker
type Option func(*ListChecker) error

// WithAdditionalFile menambahkan dataset tambahan dari path file lokal
func WithAdditionalFile(filepath string) Option {
	return func(c *ListChecker) error {
		return c.LoadFromFile(filepath)
	}
}

// WithAdditionalReader menambahkan dataset tambahan dari io.Reader
func WithAdditionalReader(r io.Reader) Option {
	return func(c *ListChecker) error {
		return c.LoadFromReader(r)
	}
}

// NewDefault membuat instance ListChecker yang sudah otomatis terisi dataset bawaan institusi Indonesia,
// serta secara otomatis memuat dataset tambahan jika environment variable OTHER_INSTITUTIONS diset.
func NewDefault(opts ...Option) (*ListChecker, error) {
	c := New()

	// 1. Muat dataset default yang di-embed ke dalam binary
	if len(defaultDatasetJSON) > 0 {
		if err := c.LoadFromReader(strings.NewReader(string(defaultDatasetJSON))); err != nil {
			return nil, fmt.Errorf("gagal memuat default dataset: %w", err)
		}
	}

	// 2. Muat dataset tambahan dari environment variable OTHER_INSTITUTIONS jika ada
	if customEnv := strings.TrimSpace(os.Getenv(EnvOtherInstitutions)); customEnv != "" {
		if err := c.loadFromEnvValue(customEnv); err != nil {
			return nil, fmt.Errorf("gagal memuat data dari env %s: %w", EnvOtherInstitutions, err)
		}
	}

	// 3. Terapkan functional options tambahan (jika ada dipanggil secara programmatis)
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// loadFromEnvValue memeriksa apakah nilai env berupa path file atau raw JSON string
func (c *ListChecker) loadFromEnvValue(val string) error {
	trimmed := strings.TrimSpace(val)

	// Jika berupa raw JSON string (diawali [ atau {)
	if strings.HasPrefix(trimmed, "[") || strings.HasPrefix(trimmed, "{") {
		return c.LoadFromReader(strings.NewReader(trimmed))
	}

	// Jika berupa path file lokal
	if _, err := os.Stat(trimmed); err == nil {
		return c.LoadFromFile(trimmed)
	}

	return fmt.Errorf("nilai %s bukan file yang valid atau string JSON yang sah: %s", EnvOtherInstitutions, val)
}
