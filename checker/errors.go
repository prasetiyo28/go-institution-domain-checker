package checker

import "errors"

var (
	// ErrInvalidDomainFormat terjadi saat format domain kosong atau tidak valid
	ErrInvalidDomainFormat = errors.New("checker: format domain tidak valid atau kosong")

	// ErrUnsupportedExtension terjadi saat ekstensi domain di luar institusi yang diizinkan
	ErrUnsupportedExtension = errors.New("checker: bukan domain institusi pendidikan yang didukung")

	// ErrDatasetEmpty terjadi saat dataset JSON yang dimuat kosong
	ErrDatasetEmpty = errors.New("checker: dataset institusi kosong")
)
