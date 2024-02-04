package setting

import (
	"fmt"
)

// persistedSettings struct is for storing the Settings of the application
type persistedSettings struct {
	Locale            string `yaml:"locale"`
	DataDir           string `yaml:"data_dir"`
	DeleteAfterClosed bool   `yaml:"delete_after_closed"`
}

func newPersistedSettings(settings Settings) persistedSettings {
	return persistedSettings{
		Locale:            settings.GetLocale(),
		DataDir:           settings.GetLatestDataDir(),
		DeleteAfterClosed: settings.GetDeleteAfterClosed(),
	}
}

type Settings struct {
	locale            string
	dataDir           string
	newDataDir        string
	deleteAfterClosed bool
}

func NewSettings(locale, dataDir string, deleteAfterClosed bool) Settings {
	return Settings{
		locale:            locale,
		dataDir:           dataDir,
		deleteAfterClosed: deleteAfterClosed,
	}
}

func (s *Settings) GetLocale() string {
	return s.locale
}

func (s *Settings) GetCurrentDataDir() string {
	return s.dataDir
}

func (s *Settings) GetLatestDataDir() string {
	if s.newDataDir != "" {
		return s.newDataDir
	}
	return s.dataDir
}

func (s *Settings) IsDataDirChanged() bool {
	return s.newDataDir != "" && s.newDataDir != s.dataDir
}

func (s *Settings) GetDeleteAfterClosed() bool {
	return s.deleteAfterClosed
}

func (s *Settings) SetLocale(locale string) error {
	if locale != "vi" && locale != "en" {
		return fmt.Errorf("invalid locale: %s", locale)
	}
	s.locale = locale
	return nil
}

func (s *Settings) SetDataDir(dataDir string) error {
	if dataDir == "" {
		return fmt.Errorf("data dir is empty")
	}
	s.newDataDir = dataDir
	return nil
}

func (s *Settings) SetDeleteAfterClosed(deleteAfterClosed bool) {
	s.deleteAfterClosed = deleteAfterClosed
}
