package handler

import (
	"net/http"
	"time"

	"github.com/gabriel-vasile/mimetype"

	"TorPlayer2/handler/uri"
	"TorPlayer2/ui"
	"TorPlayer2/vlc"
)

func (h *Handler) Watch(w http.ResponseWriter, r *http.Request, infoHash, fileName string) {
	torrentInfo, err := h.torrentManager.GetTorrentInfo(infoHash)
	if err != nil {
		handleError(w, r, "Get torrent", err, http.StatusBadRequest)
		return
	}

	// This step is for speed up the download!
	if err := h.torrentManager.CancelOthers(infoHash); err != nil {
		handleError(w, r, "Cancel others", err, http.StatusBadRequest)
		return
	}

	_ = ui.VideoPlayer(torrentInfo, fileName).Render(r.Context(), w)
}

func (h *Handler) Stream(w http.ResponseWriter, r *http.Request, infoHash, fileName string) {
	file, err := h.torrentManager.GetFile(infoHash, fileName)
	if err != nil {
		handleError(w, r, "get file", err, http.StatusBadRequest)
		return
	}

	// TODO: Maybe implement more effective file reader to seed up the download
	file.Download()
	reader := file.NewReader()
	reader.SetResponsive()
	//reader.SetReadahead(file.Length() / 100) // Read ahead 1% of the file

	mime, err := mimetype.DetectReader(reader)
	if err != nil {
		handleError(w, r, "detect mime type", err, http.StatusBadRequest)
		return
	} else {
		w.Header().Set("Content-Type", mime.String())
	}

	http.ServeContent(w, r, file.DisplayPath(), time.Time{}, reader)
}

func (h *Handler) OpenInVLC(w http.ResponseWriter, r *http.Request, infoHash, fileName string) {
	protocol := "http"
	if r.TLS != nil {
		protocol = "https"
	}

	streamURL := protocol + "://" + r.Host + uri.Stream(infoHash, fileName)

	if err := vlc.OpenURL(streamURL); err != nil {
		handleError(w, r, "Open in VLC", err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
