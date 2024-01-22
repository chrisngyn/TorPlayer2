package handler

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gabriel-vasile/mimetype"

	"TorPlayer2/subtitle"
	"TorPlayer2/ui"
)

func (h *Handler) Watch(w http.ResponseWriter, r *http.Request, infoHash, fileName string) {
	fileName, err := url.QueryUnescape(fileName)
	if err != nil {
		handleError(w, r, "Unescape file name", err, http.StatusBadRequest)
		return
	}
	torrentInfo, err := h.m.GetTorrentInfo(infoHash)
	if err != nil {
		handleError(w, r, "Get torrent", err, http.StatusBadRequest)
		return
	}

	// reset subtitle state
	h.subtitleStateStorage.SetSubtitleState(infoHash, subtitle.State{})
	_ = ui.VideoPlayer(torrentInfo, fileName).Render(r.Context(), w)
}

func (h *Handler) Stream(w http.ResponseWriter, r *http.Request, infoHash, fileName string) {
	file, err := h.m.GetFile(infoHash, fileName)
	if err != nil {
		handleError(w, r, "get file", err, http.StatusBadRequest)
		return
	}

	// TODO: Maybe implement more effective file reader to seed up the download
	reader := file.NewReader()
	reader.SetResponsive()
	reader.SetReadahead(file.Length() / 100) // Read ahead 1% of the file

	mime, err := mimetype.DetectReader(reader)
	if err != nil {
		handleError(w, r, "detect mime type", err, http.StatusBadRequest)
		return
	} else {
		w.Header().Set("Content-Type", mime.String())
	}

	http.ServeContent(w, r, file.DisplayPath(), time.Time{}, reader)
}
