package service

import (
	"context"
	"cyolo-efs/conf"
	"cyolo-efs/storage"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

type (
	FileHandler interface {
		DownloadFile(ctx context.Context, key string) ([]byte, error)
		UploadFile(ctx context.Context, filename string, content []byte, ttl time.Duration) string
	}

	service struct {
		cfg conf.Config
		fs  storage.MemoryFS
	}
)

func New(cfg conf.Config, fs storage.MemoryFS) FileHandler {
	return &service{
		cfg: cfg,
		fs:  fs,
	}
}

func (s service) DownloadFile(ctx context.Context, path string) ([]byte, error) {
	f, err := s.fs.DownloadFile(ctx, path)

	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s service) UploadFile(ctx context.Context, filename string, content []byte, ttl time.Duration) string {
	// generate a unique filename
	path := generatePath(filename)

	s.fs.UploadFile(ctx, content, path, ttl)

	return fmt.Sprintf(s.cfg.BaseUrl+"%s", path)
}

func generatePath(filename string) string {
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}
