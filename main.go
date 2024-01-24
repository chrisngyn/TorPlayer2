package main

import (
	"embed"
	"errors"
	"fmt"
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

const (
	openPort = 19576
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
		Addr:    fmt.Sprintf(":%d", openPort),
		Handler: r,
	}
	closeFns = append(closeFns, closeFn{"http server", httpServer.Close})
	go func() {
		slog.With("port", openPort).Info("Starting server")
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

	systray.Run(setupOnReady(fmt.Sprintf("http://localhost:%d", openPort)), onExit)
}

func setupOnReady(serverAddr string) func() {
	return func() {
		systray.SetIcon(appIcon)
		systray.SetTitle("TorPlayer")
		systray.SetTooltip("TorPlayer")
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
