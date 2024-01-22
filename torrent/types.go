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
}

type File struct {
	DisplayPath    string
	Length         int64
	BytesCompleted int64
}

func toInfo(tor *torrent.Torrent) Info {
	torInfo := Info{
		InfoHash:       tor.InfoHash().String(),
		Name:           tor.Name(),
		Length:         tor.Length(),
		BytesCompleted: tor.BytesCompleted(),
		Files:          make([]File, 0),
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
