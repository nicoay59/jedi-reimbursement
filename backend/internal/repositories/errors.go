package repositories

import "errors"

var (
	ErrNotFound = errors.New("data tidak ditemukan")
	ErrConflict = errors.New("data telah berubah")
)
