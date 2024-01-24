package main

import (
	"embed"
	"errors"
	"log"
	"log/slog"
	"net/http"

	"github.com/getlantern/systray"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/browser"

	"TorPlayer2/handler"
	"TorPlayer2/request"
	"TorPlayer2/setting"
	"TorPlayer2/torrent"
)

//go:embed static
var fs embed.FS

//go:embed static/appicon.ico
var appIcon []byte

type closeFn struct {
	name string
	fn   func() error
}

func main() {
	var closeFns []closeFn

	settingStorage := setting.NewStorage()
	closeFns = append(closeFns, closeFn{"setting storage", settingStorage.CleanUp})

	settings := settingStorage.GetSetting()

	m := torrent.NewManager(settings.DataDir)
	closeFns = append(closeFns, closeFn{"torrent manager", m.Close})

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(setting.Middleware(settingStorage))
	r.Use(request.Middleware)

	r.Handle("/static/*", http.FileServer(http.FS(fs)))

	httpHandler := handler.New(m, settingStorage)
	httpHandler.Register(r)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	closeFns = append(closeFns, closeFn{"http server", httpServer.Close})
	go func() {
		slog.Info("Starting server on port 8080\n")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	onExit := func() {
		for _, closeFn := range closeFns {
			if err := closeFn.fn(); err != nil {
				slog.With("error", err, "name", closeFn.name).Error("Failed to close")
			}
		}
	}

	systray.Run(setupOnReady("http://localhost:8080"), onExit)

}

func setupOnReady(serverAddr string) func() {
	return func() {
		systray.SetIcon(appIcon)
		systray.SetTitle("TorPlayer2")
		systray.SetTooltip("TorPlayer2")
		mOpenBrowser := systray.AddMenuItem("Open Torplay web", "Open Torplay in browser")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

		go func() {
			for {
				select {
				case <-mOpenBrowser.ClickedCh:
					err := browser.OpenURL(serverAddr)
					if err != nil {
						slog.With("error", err).Error("Failed to open browser")
					}
				case <-mQuit.ClickedCh:
					systray.Quit()
					return
				}
			}
		}()

		_ = browser.OpenURL(serverAddr)
	}
}
