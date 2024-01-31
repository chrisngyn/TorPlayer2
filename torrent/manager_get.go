package torrent

import (
	"fmt"

	"github.com/anacrolix/torrent"
)

func (m *Manager) GetTorrentInfo(infoHashHex string) (Info, error) {
	infoHash, err := infoHashFromHexString(infoHashHex)
	if err != nil {
		return Info{}, fmt.Errorf("parse infohash: %w", err)
	}
	tor, ok := m.client.Torrent(infoHash)
	if !ok {
		return Info{}, fmt.Errorf("torrent not found")
	}
	<-tor.GotInfo()

	return toInfo(tor), nil
}

func (m *Manager) ListTorrents(offset, limit int) ([]Info, int) {
	infos := make([]Info, 0, limit)
	torrents := m.client.Torrents()
	lower := min(offset, len(torrents))
	upper := min(offset+limit, len(torrents))
	for _, tor := range torrents[lower:upper] {
		if i := tor.Info(); i == nil {
			continue
		}
		infos = append(infos, toInfo(tor))
	}
	return infos, m.totalTorrent()
}

func (m *Manager) totalTorrent() int {
	total := 0
	for _, tor := range m.client.Torrents() {
		if tor.Info() != nil {
			total++
		}
	}
	return total
}

func (m *Manager) GetFile(infoHashHex, path string) (*torrent.File, error) {
	infoHash, err := infoHashFromHexString(infoHashHex)
	if err != nil {
		return nil, fmt.Errorf("parse infohash: %w", err)
	}

	tor, ok := m.client.Torrent(infoHash)
	if !ok {
		return nil, fmt.Errorf("torrent not found")
	}

	var file *torrent.File
	for _, f := range tor.Files() {
		if f.DisplayPath() == path {
			file = f
			break
		}
	}

	if file == nil {
		return nil, fmt.Errorf("file not found")
	}

	return file, nil
}
