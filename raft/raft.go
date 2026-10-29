// Package raft implements the Raft consensus algorithm (leader election +
// log replication). It does NOT use any external library.
package raft

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/mgueldogdu/kvstore/internal/store"
)

// State machine interface that the Raft node drives.
type StateMachine interface {
	Apply(cmd store.Command)
	Snapshot() map[string]string
	Restore(snap map[string]string)
}

type role int

const (
	follower role = iota
	candidate
	leader
)

// Node is a single Raft participant.
type Node struct {
	mu sync.Mutex

	id    string
	peers []string // raft addresses of other nodes

	// persistent state (simplified: kept in memory only)
	currentTerm uint64
	votedFor    string
	log         *Log

	// volatile state
	commitIndex uint64
	lastApplied uint64

	// leader-only volatile state
	nextIndex  map[string]uint64
	matchIndex map[string]uint64

	role     role
	leaderID string

	electionTimer *time.Timer
	heartbeatTick *time.Ticker

	sm        StateMachine
	transport *Transport

	proposeCh chan proposeReq
	stopCh    chan struct{}

	cfg nodeConfig
}