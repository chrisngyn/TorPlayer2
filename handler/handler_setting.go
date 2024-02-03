package handler

import (
	"net/http"

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

	settings := h.settingStorage.GetSettings()
	settings.DeleteAfterClosed = r.Form.Get("deleteAfterClosed") == "on"

	if err := h.settingStorage.SaveSetting(settings); err != nil {
		handleError(w, r, "Save setting", err, http.StatusInternalServerError)
		return
	}

	_ = ui.Settings(settings).Render(r.Context(), w)
}
