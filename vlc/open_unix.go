//go:build (linux || freebsd || openbsd || netbsd) && !android

package vlc

import (
	"os/exec"
)

func open(url string) error {
	cmd := exec.Command("vlc", url)
	return cmd.Run()
}
