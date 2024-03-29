package main

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"fyne.io/systray"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/browser"

	"TorPlayer2/handler"
	"TorPlayer2/i18n"
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
	closeFns = append(closeFns, closeFn{"storage", func() error {
		return cleanUpStorage(settingStorage.GetSettings())
	}})

	settings := settingStorage.GetSettings()

	m := torrent.NewManager(settings.GetCurrentDataDir())
	closeFns = append(closeFns, closeFn{"torrent manager", m.Close})

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(setting.Middleware(settingStorage))
	r.Use(request.Middleware)
	r.Use(i18n.Middleware(i18n.NewBundle(), settingStorage))

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
		for i := len(closeFns) - 1; i >= 0; i-- {
			if err := closeFns[i].fn(); err != nil {
				slog.With("error", err, "name", closeFns[i].name).Error("Failed to close")
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
		mOpenBrowser := systray.AddMenuItem("Open TorPlayer web", "Open TorPlayer in browser")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

		go func() {
			for {
				select {
				case <-mOpenBrowser.ClickedCh:
					openURL(serverAddr)
				case <-mQuit.ClickedCh:
					systray.Quit()
					return
				}
			}
		}()

		openURL(serverAddr)
	}
}

func openURL(url string) {
	err := browser.OpenURL(url)
	if err != nil {
		slog.With("error", err).Error("Failed to open browser")
	}
}

func cleanUpStorage(setting setting.Settings) error {
	if setting.GetDeleteAfterClosed() {
		dif, err := os.ReadDir(setting.GetCurrentDataDir())
		if err != nil {
			return fmt.Errorf("read data directory: %w", err)
		}
		var errs []error
		for _, f := range dif {
			if err := os.RemoveAll(filepath.Join(setting.GetCurrentDataDir(), f.Name())); err != nil {
				errs = append(errs, fmt.Errorf("remove %s: %w", f.Name(), err))
			}
		}
		return errors.Join(errs...)
	}
	return nil
}
