package arrowipc

import (
	"net/http"
	"sync"

	"laba14-health-pipeline/internal/health"
)

type Store struct {
	mu         sync.RWMutex
	aggregates []health.Aggregate
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Append(items ...health.Aggregate) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.aggregates = append(s.aggregates, items...)
}

func (s *Store) Snapshot() []health.Aggregate {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]health.Aggregate, len(s.aggregates))
	copy(out, s.aggregates)
	return out
}

func (s *Store) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/arrow/aggregates", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.apache.arrow.stream")
		if err := WriteAggregates(w, s.Snapshot()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	return mux
}
