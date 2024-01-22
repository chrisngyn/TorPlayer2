package subtitle

import (
	"bytes"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/asticode/go-astisub"
)

func Normalize(content []byte, fileExt string, addDuration time.Duration) ([]byte, error) {
	var (
		sub *astisub.Subtitles
		err error
	)

	reader := bytes.NewReader(content)

	switch strings.Trim(fileExt, ".") {
	case "srt":
		sub, err = astisub.ReadFromSRT(reader)
	case "vtt":
		sub, err = astisub.ReadFromWebVTT(reader)
	case "ssa", "ass":
		sub, err = astisub.ReadFromSSA(reader)
	default:
		return nil, fmt.Errorf("not supported file extension %s", fileExt)
	}

	if err != nil {
		return nil, fmt.Errorf("read subtitle: %w", err)
	}

	if addDuration != 0 {
		sub.Add(addDuration)
	}

	buf := &bytes.Buffer{}
	if err := sub.WriteToWebVTT(buf); err != nil {
		return nil, fmt.Errorf("write subtitle: %w", err)
	}

	return []byte(html.UnescapeString(buf.String())), nil
}
