package subtitle

import (
	"path"
	"time"
)

type State struct {
	Name                   string
	Content                []byte
	OriginalContent        []byte
	AdjustmentMilliseconds int64
}

func (s *State) IsZero() bool {
	return s.Name == "" &&
		len(s.Content) == 0 &&
		len(s.OriginalContent) == 0 &&
		s.AdjustmentMilliseconds == 0
}

func (s *State) Adjust(adjustmentMilliseconds int64) (err error) {
	s.AdjustmentMilliseconds = adjustmentMilliseconds
	s.Content, err = Normalize(s.OriginalContent, path.Ext(s.Name), time.Duration(adjustmentMilliseconds)*time.Millisecond)
	return
}
