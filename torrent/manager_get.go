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
	total := len(m.addedTorrentInfoHashes)

	for i := total - 1 - offset; (0 <= i) && (i < len(m.addedTorrentInfoHashes)) && (len(infos) < limit); i-- {
		infoHash := m.addedTorrentInfoHashes[i]
		tor, ok := m.client.Torrent(infoHash)
		if !ok {
			infos = append(infos, Info{InfoHash: infoHash.String()})
			continue
		}
		if i := tor.Info(); i == nil {
			infos = append(infos, Info{InfoHash: infoHash.String()})
			continue
		}

		infos = append(infos, toInfo(tor))
	}

	return infos, total
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
