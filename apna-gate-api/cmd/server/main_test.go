package main

import (
	"net"
	"testing"
)

type internalJobStarterFake struct {
	calls int
	start bool
}

func (f *internalJobStarterFake) StartInternalJobs() bool {
	f.calls++
	return f.start
}

func TestBindFailureDoesNotStartInternalJobs(t *testing.T) {
	var listenConfig net.ListenConfig
	occupied, err := listenConfig.Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = occupied.Close() }()

	starter := &internalJobStarterFake{start: true}
	listener, err := bindAndStartInternalJobs(occupied.Addr().String(), starter)
	if err == nil {
		if listener != nil {
			_ = listener.Close()
		}
		t.Fatal("binding an occupied port unexpectedly succeeded")
	}
	if starter.calls != 0 {
		t.Fatalf("StartInternalJobs calls = %d, want 0", starter.calls)
	}
}

func TestSuccessfulBindStartsInternalJobs(t *testing.T) {
	starter := &internalJobStarterFake{start: true}
	listener, err := bindAndStartInternalJobs("127.0.0.1:0", starter)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = listener.Close() }()
	if starter.calls != 1 {
		t.Fatalf("StartInternalJobs calls = %d, want 1", starter.calls)
	}
}
