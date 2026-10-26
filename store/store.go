// Package store provides the in-memory key-value state machine applied by Raft.
package store

import (
	"sync"
)

type Op string

const (
	OpSet    Op = "SET"
	OpDelete Op = "DELETE"
)

// Command is a single mutation that gets replicated via Raft.
type Command struct {
	Op    Op     `json:"op"`
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

// Store is a thread-safe in-memory key-value store.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

func New() *Store {
	return &Store{data: make(map[string]string)}
}

func (s *Store) Apply(cmd Command) {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch cmd.Op {
	case OpSet:
		s.data[cmd.Key] = cmd.Value
	case OpDelete:
		delete(s.data, cmd.Key)
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}

// Snapshot returns a consistent copy of the entire store.
func (s *Store) Snapshot() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap := make(map[string]string, len(s.data))
	for k, v := range s.data {
		snap[k] = v
	}
	return snap
}

// Restore replaces the store content atomically (used for snapshot install).
func (s *Store) Restore(snap map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = snap
}
