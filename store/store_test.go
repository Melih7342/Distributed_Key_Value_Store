package store

import (
	"testing"
)

func TestSetGet(t *testing.T) {
	s := New()
	s.Apply(Command{Op: OpSet, Key: "hello", Value: "world"})
	v, ok := s.Get("hello")
	if !ok || v != "world" {
		t.Fatalf("expected world got %q ok=%v", v, ok)
	}
}

func TestDelete(t *testing.T) {
	s := New()
	s.Apply(Command{Op: OpSet, Key: "k", Value: "v"})
	s.Apply(Command{Op: OpDelete, Key: "k"})
	_, ok := s.Get("k")
	if ok {
		t.Fatal("key should have been deleted")
	}
}

func TestSnapshot(t *testing.T) {
	s := New()
	s.Apply(Command{Op: OpSet, Key: "a", Value: "1"})
	snap := s.Snapshot()
	s.Apply(Command{Op: OpSet, Key: "b", Value: "2"})

	if _, ok := snap["b"]; ok {
		t.Fatal("snapshot should not include mutations after it was taken")
	}
}

func TestRestore(t *testing.T) {
	s := New()
	s.Restore(map[string]string{"x": "42"})
	v, ok := s.Get("x")
	if !ok || v != "42" {
		t.Fatalf("expected 42 got %q", v)
	}
}
