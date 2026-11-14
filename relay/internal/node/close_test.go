package node_test

import (
	"testing"
	"time"

	rt "github.com/fjaeckel/seriesoftubes/relay/internal/relaytest"
)

// Close returns only once the goroutines sessions run in have ended, and starts none after.
func TestCloseWaitsForSessionGoroutines(t *testing.T) {
	n := rt.New(t)
	x := n.Found("closing", "open").Nodes[0]
	release, ended := make(chan struct{}), make(chan struct{})
	if !x.R.Spawn(func() { <-release }) {
		t.Fatal("nothing ran on a live relay")
	}
	go func() { x.Stop(); close(ended) }()
	select {
	case <-ended:
		t.Fatal("Close returned while a session goroutine still ran")
	case <-time.After(200 * time.Millisecond):
	}
	close(release)
	<-ended
	if x.R.Spawn(func() { t.Error("ran after Close") }) {
		t.Fatal("a closed relay started a goroutine")
	}
}
