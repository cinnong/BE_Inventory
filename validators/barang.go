// validators/barang.go
package validators

import (
	"errors"
	"strings"
)

func ValidateBarang(nama string, stok int) error {
	if strings.TrimSpace(nama) == "" {
		return errors.New("nama barang wajib diisi")
	}
	if stok < 0 {
		return errors.New("stok tidak boleh negatif")
	}
	return nil
}
