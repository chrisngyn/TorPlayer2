package handler

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"golang.org/x/exp/slog"

	"TorPlayer2/subtitle"
	"TorPlayer2/ui"
)

func (h *Handler) SelectSubtitle(w http.ResponseWriter, r *http.Request, infoHash, fileName string) {
	fileName, err := url.QueryUnescape(fileName)
	if err != nil {
		handleError(w, r, "Unescape file name", err, http.StatusBadRequest)
		return
	}

	torrentInfo, err := h.torrentManager.GetTorrentInfo(infoHash)
	if err != nil {
		handleError(w, r, "Get torrent", err, http.StatusBadRequest)
		return
	}

	file, err := h.torrentManager.GetFile(infoHash, fileName)
	if err != nil {
		handleError(w, r, "Get file", err, http.StatusBadRequest)
		return
	}

	originalContent, err := io.ReadAll(file.NewReader())
	if err != nil {
		handleError(w, r, "Read file", err, http.StatusBadRequest)
		return
	}
	// because the reader read more than the file length, we need to trim it. Otherwise, the subtitle
	// may contain some data from the next file.
	originalContent = originalContent[:file.Length()]

	content, err := subtitle.Normalize(originalContent, path.Ext(file.DisplayPath()), 0)
	if err != nil {
		handleError(w, r, "Normalize subtitle", err, http.StatusBadRequest)
		return
	}

	state := subtitle.State{
		Name:            file.DisplayPath(),
		Content:         content,
		OriginalContent: originalContent,
	}

	_ = ui.Subtitle(torrentInfo, state).Render(r.Context(), w)
}

func (h *Handler) UnsetSubtitle(w http.ResponseWriter, r *http.Request, infoHash string) {
	torrentInfo, err := h.torrentManager.GetTorrentInfo(infoHash)
	if err != nil {
		handleError(w, r, "Get torrent", err, http.StatusBadRequest)
		return
	}

	_ = ui.Subtitle(torrentInfo, subtitle.State{}).Render(r.Context(), w)
}

func (h *Handler) UploadSubtitle(w http.ResponseWriter, r *http.Request, infoHash string) {
	torrentInfo, err := h.torrentManager.GetTorrentInfo(infoHash)
	if err != nil {
		handleError(w, r, "Get torrent", err, http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		handleError(w, r, "parse multipart form", err, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("fileInput")
	if err != nil {
		handleError(w, r, "get subtitle file", err, http.StatusBadRequest)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.With("error", err).ErrorCtx(r.Context(), "close file")
		}
	}()

	originalContent, err := io.ReadAll(file)
	if err != nil {
		handleError(w, r, "Read file", err, http.StatusBadRequest)
		return
	}

	content, err := subtitle.Normalize(originalContent, path.Ext(header.Filename), 0)
	if err != nil {
		handleError(w, r, "normalize subtitle", err, http.StatusBadRequest)
		return
	}

	state := subtitle.State{
		Name:            "File: " + header.Filename,
		Content:         content,
		OriginalContent: originalContent,
	}

	_ = ui.Subtitle(torrentInfo, state).Render(r.Context(), w)
}

func (h *Handler) AdjustSubtitle(w http.ResponseWriter, r *http.Request, infoHash string) {
	torrentInfo, err := h.torrentManager.GetTorrentInfo(infoHash)
	if err != nil {
		handleError(w, r, "Get torrent", err, http.StatusBadRequest)
		return
	}

	adjustmentMilliseconds, err := strconv.ParseInt(r.URL.Query().Get("adjustmentMilliseconds"), 10, 64)
	if err != nil {
		handleError(w, r, "Parse adjustment milliseconds", err, http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		handleError(w, r, "Parse form", err, http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	originalContent, err := base64.StdEncoding.DecodeString(r.FormValue("originalContent"))
	if err != nil {
		handleError(w, r, "Decode original content", err, http.StatusBadRequest)
		return
	}

	state := subtitle.State{
		Name:            name,
		OriginalContent: originalContent,
	}
	if err := state.Adjust(adjustmentMilliseconds); err != nil {
		handleError(w, r, "Adjust subtitle", err, http.StatusBadRequest)
		return
	}

	_ = ui.Subtitle(torrentInfo, state).Render(r.Context(), w)
}
