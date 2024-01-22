package handler

import (
	"net/http"
	"strconv"

	"golang.org/x/exp/slog"

	"TorPlayer2/ui"
)

func (h *Handler) ListTorrents(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	offset := max((page-1)*limit, 0)

	torrentInfos, total := h.m.ListTorrents(offset, limit)

	_ = ui.Torrents(torrentInfos, total, page, limit).Render(r.Context(), w)
}

func (h *Handler) TorrentInfo(w http.ResponseWriter, r *http.Request, infoHash string) {
	slog.With("infoHash", infoHash).InfoCtx(r.Context(), "Getting torrent info")

	tor, err := h.m.GetTorrentInfo(infoHash)
	if err != nil {
		handleError(w, r, "get torrent", err, http.StatusBadRequest)
	}

	_ = ui.Info(infoHash, tor).Render(r.Context(), w)
}
