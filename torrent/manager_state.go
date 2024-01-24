package torrent

import (
	"fmt"
)

// CancelOthers cancels all pieces of all torrents except the one with the given info hash.
func (m *Manager) CancelOthers(infoHash string) error {
	infoHashHex, err := infoHashFromHexString(infoHash)
	if err != nil {
		return fmt.Errorf("parse infohash: %w", err)
	}

	for _, tor := range m.client.Torrents() {
		if tor.InfoHash() != infoHashHex {
			tor.CancelPieces(0, tor.NumPieces())
		}
	}

	return nil
}
