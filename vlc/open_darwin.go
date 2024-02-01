//go:build !ios

package vlc

import (
	"os/exec"
)

func open(url string) error {
	cmd := exec.Command("open", url, "-a", "VLC")
	return cmd.Run()
}
