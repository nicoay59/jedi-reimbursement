package storage

import (
	"bytes"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

type memoryMultipartFile struct {
	*bytes.Reader
}

func (m memoryMultipartFile) Close() error {
	return nil
}

func TestReceiptStorageSavePNG(t *testing.T) {
	baseDir := t.TempDir()

	storage, err := NewReceiptStorage(baseDir, 1024*1024)
	if err != nil {
		t.Fatalf("NewReceiptStorage() error = %v", err)
	}

	content := append(
		[]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'},
		bytes.Repeat([]byte{0}, 128)...,
	)

	saved, err := storage.Save(
		memoryMultipartFile{Reader: bytes.NewReader(content)},
		&multipart.FileHeader{
			Filename: "bukti.png",
			Size:     int64(len(content)),
		},
	)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if saved.MIMEType != "image/png" {
		t.Fatalf("MIMEType = %q", saved.MIMEType)
	}

	fullPath := filepath.Join(
		baseDir,
		filepath.FromSlash(saved.RelativePath),
	)

	if _, err := os.Stat(fullPath); err != nil {
		t.Fatalf("file tidak tersimpan: %v", err)
	}
}

func TestReceiptStorageRejectsText(t *testing.T) {
	storage, err := NewReceiptStorage(t.TempDir(), 1024)
	if err != nil {
		t.Fatalf("NewReceiptStorage() error = %v", err)
	}

	_, err = storage.Save(
		memoryMultipartFile{
			Reader: bytes.NewReader([]byte("plain text")),
		},
		&multipart.FileHeader{
			Filename: "bukti.txt",
			Size:     10,
		},
	)
	if err != ErrReceiptType {
		t.Fatalf("Save() error = %v", err)
	}
}

func TestReceiptStorageRejectsTraversal(t *testing.T) {
	storage, err := NewReceiptStorage(t.TempDir(), 1024)
	if err != nil {
		t.Fatalf("NewReceiptStorage() error = %v", err)
	}

	if _, err := storage.Open("../secret.txt"); err != ErrReceiptUnavailable {
		t.Fatalf("Open() error = %v", err)
	}
}
