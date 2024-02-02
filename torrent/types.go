package torrent

import (
	"github.com/anacrolix/torrent"
)

type Info struct {
	InfoHash       string
	Name           string
	Length         int64
	BytesCompleted int64
	Files          []File
	Stats          Stats
}

type File struct {
	DisplayPath    string
	Length         int64
	BytesCompleted int64
}

type Stats struct {
	TotalPeers       int
	PendingPeers     int
	ActivePeers      int
	ConnectedSeeders int
	HalfOpenPeers    int
}

func toInfo(tor *torrent.Torrent) Info {
	stats := tor.Stats()
	torInfo := Info{
		InfoHash:       tor.InfoHash().String(),
		Name:           tor.Name(),
		Length:         tor.Length(),
		BytesCompleted: tor.BytesCompleted(),
		Files:          make([]File, 0),
		Stats: Stats{
			TotalPeers:       stats.TotalPeers,
			PendingPeers:     stats.PendingPeers,
			ActivePeers:      stats.ActivePeers,
			ConnectedSeeders: stats.ConnectedSeeders,
			HalfOpenPeers:    stats.HalfOpenPeers,
		},
	}

	for _, file := range tor.Files() {
		torInfo.Files = append(torInfo.Files, File{
			DisplayPath:    file.DisplayPath(),
			Length:         file.Length(),
			BytesCompleted: file.BytesCompleted(),
		})
	}
	return torInfo
}
