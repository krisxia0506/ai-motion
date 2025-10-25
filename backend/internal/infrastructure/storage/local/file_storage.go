package local

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrFileNotFound      = errors.New("file not found")
	ErrInvalidFileType   = errors.New("invalid file type")
	ErrFileTooLarge      = errors.New("file too large")
	ErrStoragePathNotSet = errors.New("storage path not set")
)

const (
	MaxFileSize = 100 * 1024 * 1024
)

var AllowedMimeTypes = map[string]bool{
	"text/plain":           true,
	"application/pdf":      true,
	"application/epub+zip": true,
	"image/jpeg":           true,
	"image/png":            true,
	"image/gif":            true,
	"image/webp":           true,
	"video/mp4":            true,
	"video/webm":           true,
}

type FileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) (*FileStorage, error) {
	if basePath == "" {
		return nil, ErrStoragePathNotSet
	}

	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &FileStorage{
		basePath: basePath,
	}, nil
}

func (s *FileStorage) Upload(ctx context.Context, filename string, reader io.Reader, fileSize int64) (string, error) {
	if fileSize > MaxFileSize {
		return "", ErrFileTooLarge
	}

	mimeType := mime.TypeByExtension(filepath.Ext(filename))
	if !AllowedMimeTypes[mimeType] {
		return "", ErrInvalidFileType
	}

	relPath := s.generateFilePath(filename)
	fullPath := filepath.Join(s.basePath, relPath)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	written, err := io.Copy(file, reader)
	if err != nil {
		os.Remove(fullPath)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	if written != fileSize {
		os.Remove(fullPath)
		return "", fmt.Errorf("file size mismatch: expected %d, got %d", fileSize, written)
	}

	return relPath, nil
}

func (s *FileStorage) Download(ctx context.Context, filePath string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, filePath)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

func (s *FileStorage) Delete(ctx context.Context, filePath string) error {
	fullPath := filepath.Join(s.basePath, filePath)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotFound
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (s *FileStorage) Exists(ctx context.Context, filePath string) (bool, error) {
	fullPath := filepath.Join(s.basePath, filePath)

	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *FileStorage) GetFileSize(ctx context.Context, filePath string) (int64, error) {
	fullPath := filepath.Join(s.basePath, filePath)

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, ErrFileNotFound
		}
		return 0, fmt.Errorf("failed to stat file: %w", err)
	}

	return info.Size(), nil
}

func (s *FileStorage) generateFilePath(filename string) string {
	now := time.Now()
	dateDir := now.Format("2006/01/02")

	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	safeFilename := sanitizeFilename(nameWithoutExt)

	timestamp := now.Format("150405")
	randomSuffix := randomString(6)

	finalFilename := fmt.Sprintf("%s_%s_%s%s", safeFilename, timestamp, randomSuffix, ext)

	return filepath.Join(dateDir, finalFilename)
}

func sanitizeFilename(filename string) string {
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return -1
	}, filename)

	if len(filename) > 50 {
		filename = filename[:50]
	}

	return filename
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
