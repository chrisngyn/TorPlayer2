package setting

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/adrg/xdg"
)

type Storage struct {
	mu      sync.Mutex
	setting Settings
}

func NewStorage() *Storage {
	dataDir := filepath.Join(xdg.UserDirs.Download, "TorPlayer")
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.Mkdir(dataDir, 0700); err != nil {
			panic(fmt.Errorf("create data directory: %w", err))
		}
	}
	return &Storage{
		setting: Settings{
			DataDir:           dataDir,
			DeleteAfterClosed: true,
		},
	}
}

func (s *Storage) GetSetting() Settings {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.setting
}

func (s *Storage) SaveSetting(setting Settings) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.setting = setting
	return nil
}

func (s *Storage) CleanUp() error {
	setting := s.GetSetting()
	if setting.DeleteAfterClosed {
		if err := os.RemoveAll(setting.DataDir); err != nil {
			return fmt.Errorf("remove data directory: %w", err)
		}
	}
	return nil
}
