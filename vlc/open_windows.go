//go:build windows

package vlc

import (
	"os/exec"
)

func open(url string) error {
	cmd := exec.Command(`C:\Program Files\VideoLAN\VLC\vlc.exe`, url)
	return cmd.Run()
}
