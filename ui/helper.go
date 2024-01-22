package ui

import (
	"fmt"
	"strings"
)

func byteCounter(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func percent(value, total int64) string {
	return fmt.Sprintf("%.1f%%", float64(value)/float64(total)*100)
}

func isVideoFile(filename string) bool {
	return strings.HasSuffix(filename, ".mp4")
}

var (
	subtitleFileExtensions = []string{".srt", ".vtt"}
)

func isSubtitleFile(filename string) bool {
	for _, ext := range subtitleFileExtensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}
