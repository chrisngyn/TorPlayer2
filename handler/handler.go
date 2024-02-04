package handler

import (
	"net/http"
	"net/url"

	"github.com/gen2brain/beeep"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"

	"TorPlayer2/setting"
	"TorPlayer2/torrent"
)

type Handler struct {
	torrentManager *torrent.Manager
	settingStorage *setting.Storage
}

func New(m *torrent.Manager, settingStorage *setting.Storage) *Handler {
	if m == nil {
		panic("torrent manager is required")
	}

	return &Handler{
		torrentManager: m,
		settingStorage: settingStorage,
	}
}

func (h *Handler) Register(r chi.Router) {
	// home page
	r.Get("/", h.Home)
	r.Post("/add-torrent", h.AddTorrent)

	// torrent info
	r.Get("/torrents", h.ListTorrents)
	r.Get("/torrents/{infoHash}", func(w http.ResponseWriter, r *http.Request) {
		infoHash := chi.URLParam(r, "infoHash")
		h.TorrentInfo(w, r, infoHash)
	})

	// watch torrent
	r.Get("/torrents/{infoHash}/watch/{fileName}", func(w http.ResponseWriter, r *http.Request) {
		infoHash := chi.URLParam(r, "infoHash")
		fileName, err := url.QueryUnescape(chi.URLParam(r, "fileName"))
		if err != nil {
			handleError(w, r, "Unescape file name", err, http.StatusBadRequest)
			return
		}
		h.Watch(w, r, infoHash, fileName)
	})

	// subtitle
	r.Post("/torrents/{infoHash}/select-subtitle/{fileName}", func(w http.ResponseWriter, r *http.Request) {
		infoHash := chi.URLParam(r, "infoHash")
		fileName, err := url.QueryUnescape(chi.URLParam(r, "fileName"))
		if err != nil {
			handleError(w, r, "Unescape file name", err, http.StatusBadRequest)
			return
		}
		h.SelectSubtitle(w, r, infoHash, fileName)
	})
	r.Delete("/torrents/{infoHash}/unset-subtitle", func(w http.ResponseWriter, r *http.Request) {
		infoHash := chi.URLParam(r, "infoHash")
		h.UnsetSubtitle(w, r, infoHash)
	})
	r.Post("/torrents/{infoHash}/upload-subtitle", func(w http.ResponseWriter, r *http.Request) {
		infoHash := chi.URLParam(r, "infoHash")
		h.UploadSubtitle(w, r, infoHash)
	})
	r.Post("/torrents/{infoHash}/adjust-subtitle", func(w http.ResponseWriter, r *http.Request) {
		infoHash := chi.URLParam(r, "infoHash")
		h.AdjustSubtitle(w, r, infoHash)
	})

	// stream torrent
	r.Get("/stream/{infoHash}/{fileName}", func(w http.ResponseWriter, r *http.Request) {
		infoHash := chi.URLParam(r, "infoHash")
		fileName, err := url.QueryUnescape(chi.URLParam(r, "fileName"))
		if err != nil {
			handleError(w, r, "Unescape file name", err, http.StatusBadRequest)
			return
		}
		h.Stream(w, r, infoHash, fileName)
	})

	// open in vlc
	r.Get("/open-in-vlc/{infoHash}/{fileName}", func(w http.ResponseWriter, r *http.Request) {
		infoHash := chi.URLParam(r, "infoHash")
		fileName, err := url.QueryUnescape(chi.URLParam(r, "fileName"))
		if err != nil {
			handleError(w, r, "Unescape file name", err, http.StatusBadRequest)
			return
		}
		h.OpenInVLC(w, r, infoHash, fileName)
	})

	// setting
	r.Get("/settings", h.GetSettings)
	r.Post("/settings", h.UpdateSetting)
	r.Post("/settings/change-data-dir", h.ChangeDataDir)

}

func handleError(w http.ResponseWriter, r *http.Request, msg string, err error, status int) {
	slog.With("error", err).ErrorCtx(r.Context(), msg)
	alertError(msg, err)
	http.Error(w, msg+": "+err.Error(), status)
}

func alertError(msg string, err error) {
	_ = beeep.Alert("Error", msg+": "+err.Error(), "")
}

func redirect(w http.ResponseWriter, r *http.Request, url string) {
	w.Header().Set("HX-Redirect", url)
	http.Redirect(w, r, url, http.StatusSeeOther)
}
