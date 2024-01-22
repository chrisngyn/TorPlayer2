package subtitle

import (
	"path"
	"sync"
	"time"
)

type StateStorage struct {
	mu            sync.Mutex
	subtitleState map[string]State
}

func NewStateStorage() *StateStorage {
	return &StateStorage{
		subtitleState: make(map[string]State),
	}
}

func (s *StateStorage) GetSubtitleState(infoHash string) State {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.subtitleState[infoHash]
}

func (s *StateStorage) SetSubtitleState(infoHash string, state State) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.subtitleState[infoHash] = state
}

func (s *StateStorage) UnsetSubtitleState(infoHash string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.subtitleState, infoHash)
}

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
