package validators

import (
	"errors"
	"regexp"
	"strings"
)

func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !re.MatchString(email) {
		return errors.New("format email tidak valid")
	}
	return nil
}

func ValidateTelepon(telp string) error {
	re := regexp.MustCompile(`^[0-9]+$`)
	if !re.MatchString(telp) {
		return errors.New("telepon hanya boleh berisi angka")
	}
	return nil
}

func ValidatePeminjaman(nama string, email string, telp string, jumlah int, status string) error {
	if strings.TrimSpace(nama) == "" {
		return errors.New("nama peminjam wajib diisi")
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if err := ValidateTelepon(telp); err != nil {
		return err
	}
	if jumlah <= 0 {
		return errors.New("jumlah pinjam harus lebih dari 0")
	}
	if status != "dipinjam" && status != "dikembalikan" {
		return errors.New("status harus 'dipinjam' atau 'dikembalikan'")
	}
	return nil
}
