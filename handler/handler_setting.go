package handler

import (
	"net/http"

	"TorPlayer2/ui"
)

func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
	_ = ui.Settings().Render(r.Context(), w)
}

func (h *Handler) UpdateSetting(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handleError(w, r, "Parse form", err, http.StatusBadRequest)
		return
	}

	setting := h.settingStorage.GetSetting()
	setting.DeleteAfterClosed = r.Form.Get("deleteAfterClosed") == "on"

	if err := h.settingStorage.SaveSetting(setting); err != nil {
		handleError(w, r, "Save setting", err, http.StatusInternalServerError)
		return
	}

	_ = ui.Settings().Render(r.Context(), w)
}
