package handler

import (
	"fmt"
	"net/http"
	"path"

	"github.com/a-h/templ"
	"golang.org/x/exp/slog"

	"TorPlayer2/ui"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	templ.Handler(ui.Index()).ServeHTTP(w, r)
}

func (h *Handler) AddTorrent(w http.ResponseWriter, r *http.Request) {
	var (
		infoHash string
		err      error
	)

	textInput := r.FormValue("textInput")
	if textInput != "" {
		slog.With("textInput", textInput).InfoCtx(r.Context(), "Adding torrent from text input")
		infoHash, err = h.addTorrentFromTextInput(textInput)

	} else {
		slog.InfoCtx(r.Context(), "Adding torrent from file input")
		infoHash, err = h.addTorrentFromFile(r)
	}

	if err != nil {
		handleError(w, r, "add torrent", err, http.StatusInternalServerError)
		return
	}

	infoURL := path.Join("/torrents/", infoHash)
	http.Redirect(w, r, infoURL, http.StatusSeeOther)
}

func (h *Handler) addTorrentFromTextInput(textInput string) (string, error) {
	return h.m.AddTorrentFromString(textInput)
}

func (h *Handler) addTorrentFromFile(r *http.Request) (string, error) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return "", fmt.Errorf("parse multipart form: %w", err)
	}
	file, _, err := r.FormFile("fileInput")
	if err != nil {
		return "", fmt.Errorf("get torrent file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.With("error", err).ErrorCtx(r.Context(), "close file")
		}
	}()

	infoHash, err := h.m.AddTorrentFromFileContent(file)
	if err != nil {
		return "", fmt.Errorf("add torrent: %w", err)
	}

	return infoHash, nil
}
