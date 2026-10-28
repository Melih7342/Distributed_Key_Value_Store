package raft

// RPC message types exchanged between Raft peers.

// VoteRequest is sent by a candidate to request a vote.
type VoteRequest struct {
	Term         uint64 `json:"term"`
	CandidateID  string `json:"candidate_id"`
	LastLogIndex uint64 `json:"last_log_index"`
	LastLogTerm  uint64 `json:"last_log_term"`
}

// VoteResponse is the reply to a VoteRequest.
type VoteResponse struct {
	Term        uint64 `json:"term"`
	VoteGranted bool   `json:"vote_granted"`
}

// AppendRequest is sent by the leader to replicate log entries (also used as heartbeat).
type AppendRequest struct {
	Term         uint64  `json:"term"`
	LeaderID     string  `json:"leader_id"`
	PrevLogIndex uint64  `json:"prev_log_index"`
	PrevLogTerm  uint64  `json:"prev_log_term"`
	Entries      []Entry `json:"entries"`
	LeaderCommit uint64  `json:"leader_commit"`
}

// AppendResponse is the reply to an AppendRequest.
type AppendResponse struct {
	Term    uint64 `json:"term"`
	Success bool   `json:"success"`
	// ConflictIndex allows the leader to skip back quickly on failure.
	ConflictIndex uint64 `json:"conflict_index"`
}

// InstallSnapshotRequest carries a full state-machine snapshot.
type InstallSnapshotRequest struct {
	Term      uint64            `json:"term"`
	LeaderID  string            `json:"leader_id"`
	LastIndex uint64            `json:"last_index"`
	LastTerm  uint64            `json:"last_term"`
	Data      map[string]string `json:"data"`
}

// InstallSnapshotResponse is the reply to an InstallSnapshotRequest.
type InstallSnapshotResponse struct {
	Term uint64 `json:"term"`
}
