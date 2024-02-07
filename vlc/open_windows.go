//go:build windows

package vlc

import (
	"os"
	"os/exec"
)

func open(url string) error {
	vlcPath, err := getVlcPath()
	if err != nil {
		return err
	}
	cmd := exec.Command(vlcPath, url)
	return cmd.Run()
}

const (
	vlcPath64 = `C:\Program Files\VideoLAN\VLC\vlc.exe`
	vlcPath32 = `C:\Program Files (x86)\VideoLAN\VLC\vlc.exe`
)

func getVlcPath() (string, error) {
	if _, err := os.Stat(vlcPath64); err == nil {
		return vlcPath64, nil
	}
	if _, err := os.Stat(vlcPath32); err == nil {
		return vlcPath32, nil
	}
	return "", os.ErrNotExist
}
