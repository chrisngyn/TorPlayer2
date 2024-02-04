package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ncruces/zenity"

	"TorPlayer2/handler/uri"
	"TorPlayer2/setting"
	"TorPlayer2/ui"
)

func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings := h.settingStorage.GetSettings()
	err := ui.Settings(settings).Render(r.Context(), w)
	if err != nil {
		handleError(w, r, "Render settings", err, http.StatusInternalServerError)
	}
}

func (h *Handler) UpdateSetting(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handleError(w, r, "Parse form", err, http.StatusBadRequest)
		return
	}

	err := h.settingStorage.UpdateSettings(func(s setting.Settings) (setting.Settings, error) {
		if r.Form.Has("locale") {
			if err := s.SetLocale(r.Form.Get("locale")); err != nil {
				return s, fmt.Errorf("set locale: %w", err)
			}
		}
		if r.Form.Has("deleteAfterClosed") {
			s.SetDeleteAfterClosed(r.Form.Get("deleteAfterClosed") == "on")
		}
		return s, nil
	})
	if err != nil {
		handleError(w, r, "Update setting", err, http.StatusInternalServerError)
		return
	}

	redirect(w, r, uri.GetSettings())
}

func (h *Handler) ChangeDataDir(w http.ResponseWriter, r *http.Request) {
	err := h.settingStorage.UpdateSettings(func(s setting.Settings) (setting.Settings, error) {
		newDataDir, err := zenity.SelectFile(
			zenity.Directory(),
			zenity.Title("Select data directory"),
			zenity.Filename(s.GetCurrentDataDir()),
		)
		if err != nil {
			return s, fmt.Errorf("select data directory: %w", err)
		}

		if err := s.SetDataDir(newDataDir); err != nil {
			return s, fmt.Errorf("set data directory: %w", err)
		}

		return s, nil
	})
	if errors.Is(err, zenity.ErrCanceled) {
		redirect(w, r, uri.GetSettings())
		return
	}
	if err != nil {
		handleError(w, r, "Update setting", err, http.StatusInternalServerError)
		return
	}

	redirect(w, r, uri.GetSettings())
}
