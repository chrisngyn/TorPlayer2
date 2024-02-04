// Package uri provides URI generation helpers for the web interface.
package uri

import (
	"fmt"
	"net/url"
)

func Stream(infoHash, fileName string) string {
	return fmt.Sprintf("/stream/%s/%s", infoHash, url.QueryEscape(fileName))
}

func OpenInVLC(infoHash, fileName string) string {
	return fmt.Sprintf("/open-in-vlc/%s/%s", infoHash, url.QueryEscape(fileName))
}

func Watch(infoHash, fileName string) string {
	return fmt.Sprintf("/torrents/%s/watch/%s", infoHash, url.QueryEscape(fileName))
}

func Info(infoHash string) string {
	return fmt.Sprintf("/torrents/%s", infoHash)
}

func UnsetSubtitle(infoHash string) string {
	return fmt.Sprintf("/torrents/%s/unset-subtitle", infoHash)
}

func SelectSubtitle(infoHash, fileName string) string {
	return fmt.Sprintf("/torrents/%s/select-subtitle/%s", infoHash, url.QueryEscape(fileName))
}

func GetSettings() string {
	return "/settings"
}

func UpdateSettings() string {
	return "/settings"
}

func ChangeDataDir() string {
	return "/settings/change-data-dir"
}
