package storage

import (
	"context"
	"cyolo-efs/conf"
	logger "cyolo-efs/logging"
	"fmt"
	"time"
)

type (
	MemoryFS interface {
		DownloadFile(ctx context.Context, path string) ([]byte, error)
		UploadFile(ctx context.Context, data []byte, path string, ttl time.Duration)
		Title() string
		Start()
		Stop()
	}

	fileData struct {
		data     []byte
		addedAt  time.Time
		removeAt time.Time
	}

	addFileMsg struct {
		path string
		data []byte
		ttl  time.Duration
	}

	getFileMsg struct {
		path   string
		respCh chan getFileResp
	}

	getFileResp struct {
		data []byte
		err  error
	}

	storage struct {
		cfg      conf.Config
		files    map[string]*fileData
		fileChan chan string

		uploadCh   chan addFileMsg
		downloadCh chan getFileMsg
		stopCh     chan struct{}
	}
)

func (s *storage) Title() string {
	return "Storage"
}

func New(cfg conf.Config) MemoryFS {
	s := &storage{
		cfg:   cfg,
		files: make(map[string]*fileData),

		uploadCh:   make(chan addFileMsg),
		downloadCh: make(chan getFileMsg),
		stopCh:     make(chan struct{}),
	}

	return s
}

func (s *storage) DownloadFile(_ context.Context, path string) ([]byte, error) {
	respCh := make(chan getFileResp)
	s.downloadCh <- getFileMsg{
		path:   path,
		respCh: respCh,
	}
	resp := <-respCh
	return resp.data, resp.err
}

func (s *storage) UploadFile(_ context.Context, data []byte, path string, ttl time.Duration) {
	s.uploadCh <- addFileMsg{
		path: path,
		data: data,
		ttl:  ttl,
	}
}

func (s *storage) Start() {
	removeExpiredFilesCh := time.Tick(s.cfg.FilesCleanIntervalMil)

	for {
		select {
		case msg := <-s.uploadCh:
			s.files[msg.path] = &fileData{
				data:     msg.data,
				addedAt:  time.Now(),
				removeAt: time.Now().Add(msg.ttl),
			}
		case msg := <-s.downloadCh:
			file, ok := s.files[msg.path]
			if !ok {
				msg.respCh <- getFileResp{
					data: nil,
					err:  fmt.Errorf("file not found"),
				}
			} else if now := time.Now(); now.After(file.removeAt) {
				logger.Log.Infof("Deleted file: %s", msg.path)
				delete(s.files, msg.path)
				msg.respCh <- getFileResp{
					data: nil,
					err:  fmt.Errorf("file not found"),
				}
			} else {
				msg.respCh <- getFileResp{
					data: file.data,
					err:  nil,
				}
			}
		case <-removeExpiredFilesCh:
			now := time.Now()
			for path, file := range s.files {
				if now.After(file.removeAt) {
					logger.Log.Infof("Deleted file: %s", path)
					delete(s.files, path)
				}
			}
		case <-s.stopCh:
			return
		}
	}
}

func (s *storage) Stop() {
	close(s.stopCh)
}
