package validators

import (
	"errors"
	"strings"
)

func ValidateKategori(nama, deskripsi string) error {
	if strings.TrimSpace(nama) == "" {
		return errors.New("nama kategori wajib diisi")
	}
	if len(nama) < 1 {
		return errors.New("nama kategori minimal 1 karakter")
	}
	if len(deskripsi) > 100 {
		return errors.New("deskripsi terlalu panjang")
	}
	return nil
}
