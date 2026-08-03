package storage

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrReceiptTooLarge    = errors.New("ukuran bukti terlalu besar")
	ErrReceiptType        = errors.New("jenis bukti tidak didukung")
	ErrReceiptUnavailable = errors.New("bukti tidak tersedia")
)

type SavedReceipt struct {
	RelativePath string
	OriginalName string
	MIMEType     string
	Size         int64
}

type ReceiptStorage struct {
	baseDir   string
	targetDir string
	maxBytes  int64
}

func NewReceiptStorage(
	baseDir string,
	maxBytes int64,
) (*ReceiptStorage, error) {
	absoluteBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, fmt.Errorf("menentukan folder upload: %w", err)
	}

	targetDir := filepath.Join(
		absoluteBaseDir,
		"parking-receipts",
	)

	if err := os.MkdirAll(targetDir, 0o750); err != nil {
		return nil, fmt.Errorf("membuat folder bukti parkir: %w", err)
	}

	return &ReceiptStorage{
		baseDir:   absoluteBaseDir,
		targetDir: targetDir,
		maxBytes:  maxBytes,
	}, nil
}

func (s *ReceiptStorage) Save(
	file multipart.File,
	header *multipart.FileHeader,
) (SavedReceipt, error) {
	if header == nil {
		return SavedReceipt{}, ErrReceiptUnavailable
	}

	if header.Size > s.maxBytes {
		return SavedReceipt{}, ErrReceiptTooLarge
	}

	buffer := make([]byte, 512)
	readBytes, readErr := file.Read(buffer)
	if readErr != nil && !errors.Is(readErr, io.EOF) {
		return SavedReceipt{}, fmt.Errorf("membaca bukti: %w", readErr)
	}

	buffer = buffer[:readBytes]
	mimeType := http.DetectContentType(buffer)

	extension, ok := receiptExtension(mimeType)
	if !ok {
		return SavedReceipt{}, ErrReceiptType
	}

	randomName, err := randomFileName(extension)
	if err != nil {
		return SavedReceipt{}, err
	}

	fullPath := filepath.Join(s.targetDir, randomName)
	output, err := os.OpenFile(
		fullPath,
		os.O_CREATE|os.O_EXCL|os.O_WRONLY,
		0o640,
	)
	if err != nil {
		return SavedReceipt{}, fmt.Errorf("membuat file bukti: %w", err)
	}

	success := false
	defer func() {
		_ = output.Close()
		if !success {
			_ = os.Remove(fullPath)
		}
	}()

	totalWritten, err := output.Write(buffer)
	if err != nil {
		return SavedReceipt{}, fmt.Errorf("menulis bukti: %w", err)
	}

	remainingLimit := s.maxBytes - int64(totalWritten) + 1
	if remainingLimit < 1 {
		remainingLimit = 1
	}

	written, err := io.Copy(
		output,
		io.LimitReader(file, remainingLimit),
	)
	if err != nil {
		return SavedReceipt{}, fmt.Errorf("menyimpan bukti: %w", err)
	}

	totalSize := int64(totalWritten) + written
	if totalSize > s.maxBytes {
		return SavedReceipt{}, ErrReceiptTooLarge
	}

	if err := output.Sync(); err != nil {
		return SavedReceipt{}, fmt.Errorf("menyelesaikan file bukti: %w", err)
	}

	success = true

	originalName := filepath.Base(strings.TrimSpace(header.Filename))
	if originalName == "." || originalName == "" {
		originalName = "bukti" + extension
	}
	if len(originalName) > 255 {
		originalName = originalName[:255]
	}

	return SavedReceipt{
		RelativePath: filepath.ToSlash(
			filepath.Join("parking-receipts", randomName),
		),
		OriginalName: originalName,
		MIMEType:     mimeType,
		Size:         totalSize,
	}, nil
}

func (s *ReceiptStorage) Open(
	relativePath string,
) (*os.File, error) {
	fullPath, err := s.safePath(relativePath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(fullPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrReceiptUnavailable
	}
	if err != nil {
		return nil, fmt.Errorf("membuka bukti: %w", err)
	}

	return file, nil
}

func (s *ReceiptStorage) Delete(
	relativePath string,
) error {
	fullPath, err := s.safePath(relativePath)
	if err != nil {
		return err
	}

	if err := os.Remove(fullPath); err != nil &&
		!errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("menghapus bukti: %w", err)
	}

	return nil
}

func (s *ReceiptStorage) safePath(
	relativePath string,
) (string, error) {
	cleanPath := filepath.Clean(filepath.FromSlash(relativePath))
	if cleanPath == "." ||
		filepath.IsAbs(cleanPath) ||
		cleanPath == ".." ||
		strings.HasPrefix(cleanPath, ".."+string(filepath.Separator)) {
		return "", ErrReceiptUnavailable
	}

	fullPath := filepath.Join(s.baseDir, cleanPath)
	relativeToBase, err := filepath.Rel(s.baseDir, fullPath)
	if err != nil ||
		relativeToBase == ".." ||
		strings.HasPrefix(
			relativeToBase,
			".."+string(filepath.Separator),
		) {
		return "", ErrReceiptUnavailable
	}

	return fullPath, nil
}

func receiptExtension(mimeType string) (string, bool) {
	switch mimeType {
	case "image/jpeg":
		return ".jpg", true
	case "image/png":
		return ".png", true
	case "application/pdf":
		return ".pdf", true
	default:
		return "", false
	}
}

func randomFileName(extension string) (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("membuat nama file acak: %w", err)
	}

	return hex.EncodeToString(bytes) + extension, nil
}
