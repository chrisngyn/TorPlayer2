package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"TorPlayer2/handler"
	"TorPlayer2/setting"
	"TorPlayer2/torrent"
)

//go:embed static
var fs embed.FS

func main() {
	settingStorage := setting.NewStorage()
	defer executeCleanupFunc("setting storage", settingStorage.CleanUp)

	settings := settingStorage.GetSetting()

	m := torrent.NewManager(settings.DataDir)
	defer executeCleanupFunc("torrent manager", m.Close)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(setting.Middleware(settingStorage))
	r.Use(middleware.Recoverer)
	r.Handle("/static/*", http.FileServer(http.FS(fs)))

	httpHandler := handler.New(m, settingStorage)
	httpHandler.Register(r)

	slog.Info("Starting server on port 8080\n")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

func executeCleanupFunc(entityName string, cleanupFunc func() error) {
	if err := cleanupFunc(); err != nil {
		slog.With("error", err).Error(fmt.Sprintf("Failed to clean up %s", entityName))
	}
}
