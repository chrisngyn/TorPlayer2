package setting

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
)

type Storage struct {
	configFilePath string
	mu             sync.RWMutex
	setting        Settings
	newDataDir     string
}

func NewStorage() *Storage {
	configFilePath, err := xdg.ConfigFile("TorPlayer/user_settings.yaml")
	if err != nil {
		panic(fmt.Errorf("get config file path: %w", err))
	}

	if err := intSettingFileIfNotExist(configFilePath); err != nil {
		panic(fmt.Errorf("init setting file: %w", err))
	}

	settings, err := loadSettingsFromFile(configFilePath)
	if err != nil {
		panic(fmt.Errorf("load setting file: %w", err))
	}

	return &Storage{
		configFilePath: configFilePath,
		setting:        settings,
	}
}

func (s *Storage) GetSettings() Settings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.setting
}

func (s *Storage) GetLanguage() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.setting.GetLocale()
}

func (s *Storage) UpdateSettings(updateFn func(settings Settings) (Settings, error)) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	newSetting, err := updateFn(s.setting)
	if err != nil {
		return fmt.Errorf("update setting: %w", err)
	}

	if err := saveSettingsToFile(s.configFilePath, newSetting); err != nil {
		return fmt.Errorf("save setting: %w", err)
	}

	s.setting = newSetting
	return nil
}

func intSettingFileIfNotExist(configFilePath string) error {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		setting := defaultSetting()
		if err := saveSettingsToFile(configFilePath, setting); err != nil {
			return fmt.Errorf("save default setting: %w", err)
		}
	}

	return nil
}

func defaultSetting() Settings {
	dataDir := filepath.Join(xdg.UserDirs.Download, "TorPlayer")
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.Mkdir(dataDir, 0700); err != nil {
			panic(fmt.Errorf("create data directory: %w", err))
		}
	}
	return NewSettings("vi", dataDir, true)
}

func loadSettingsFromFile(configFilePath string) (Settings, error) {
	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		return Settings{}, fmt.Errorf("read yaml file: %w", err)
	}

	var setting persistedSettings
	if err := yaml.Unmarshal(yamlFile, &setting); err != nil {
		return Settings{}, fmt.Errorf("load yaml file: %w", err)
	}

	if setting.Locale == "" {
		setting.Locale = "vi"
	}

	return NewSettings(setting.Locale, setting.DataDir, setting.DeleteAfterClosed), nil
}

func saveSettingsToFile(configFilePath string, settings Settings) error {
	yamlFile, err := yaml.Marshal(newPersistedSettings(settings))
	if err != nil {
		return fmt.Errorf("marshal yaml file: %w", err)
	}

	if err := os.WriteFile(configFilePath, yamlFile, 0644); err != nil {
		return fmt.Errorf("write yaml file: %w", err)
	}
	return nil
}
