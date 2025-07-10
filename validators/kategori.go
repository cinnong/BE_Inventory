package validators

import (
	"errors"
	"strings"
)

func ValidateKategori(nama, deskripsi string) error {
	if strings.TrimSpace(nama) == "" {
		return errors.New("nama kategori wajib diisi")
	}
	if len(nama) < 3 {
		return errors.New("nama kategori minimal 3 karakter")
	}
	if len(deskripsi) > 100 {
		return errors.New("deskripsi terlalu panjang")
	}
	return nil
}
