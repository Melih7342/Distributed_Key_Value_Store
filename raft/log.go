package raft

import (
	"encoding/json"
)

// Entry is a single record in the Raft replicated log.
type Entry struct {
	Index   uint64          `json:"index"`
	Term    uint64          `json:"term"`
	Command json.RawMessage `json:"command"` // opaque payload – decoded by the state machine
}

// Log is an in-memory append-only log of entries.
// Index is 1-based; index 0 is a sentinel "empty" entry.
type Log struct {
	snapshotIndex uint64
	entries       []Entry
}

func newLog() *Log {
	// entry at index 0 acts as a virtual "previous" entry before any real log
	return &Log{entries: []Entry{{Index: 0, Term: 0}}}
}

func (l *Log) lastIndex() uint64 {
	return l.entries[len(l.entries)-1].Index
}

func (l *Log) lastTerm() uint64 {
	return l.entries[len(l.entries)-1].Term
}

// entry returns the entry at logical index i (1-based).
// Returns zero value if i == 0.
func (l *Log) entry(i uint64) Entry {
	if i == 0 || i >= uint64(len(l.entries)) {
		return Entry{}
	}
	physical := i - l.snapshotIndex
	return l.entries[physical]
}

// termOf returns the term of the entry at index i.
func (l *Log) termOf(i uint64) uint64 {
	return l.entry(i).Term
}

// append adds entries, truncating any conflicting suffix first.
func (l *Log) append(prevIndex uint64, entries []Entry) {
	// Truncate any log entries past prevIndex that conflict.
	l.entries = l.entries[:prevIndex+1]
	l.entries = append(l.entries, entries...)
}

// entriesFrom returns all entries with index > after.
func (l *Log) entriesFrom(after uint64) []Entry {
	physical := after - l.snapshotIndex
	if physical+1 >= uint64(len(l.entries)) {
		return nil
	}
	return l.entries[physical+1:]
}

// slice returns entries in range [lo, hi) (1-based).
func (l *Log) slice(lo, hi uint64) []Entry {
	physical_lo := lo - l.snapshotIndex
	physical_hi := hi - l.snapshotIndex

	if lo >= hi || physical_lo >= uint64(len(l.entries)) {
		return nil
	}
	end := physical_hi
	if int(end) > len(l.entries) {
		end = uint64(len(l.entries))
	}
	return l.entries[physical_lo:end]
}
