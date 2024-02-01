// Package uri provides URI generation helpers for the web interface.
package uri

import (
	"fmt"
	"net/url"
)

func StreamURI(infoHash, fileName string) string {
	return fmt.Sprintf("/stream/%s/%s", infoHash, url.QueryEscape(fileName))
}

func OpenInVLCURI(infoHash, fileName string) string {
	return fmt.Sprintf("/open-in-vlc/%s/%s", infoHash, url.QueryEscape(fileName))
}

func WatchURI(infoHash, fileName string) string {
	return fmt.Sprintf("/torrents/%s/watch/%s", infoHash, url.QueryEscape(fileName))
}

func InfoURI(infoHash string) string {
	return fmt.Sprintf("/torrents/%s", infoHash)
}
