package registrar

import (
	"context"
	"testing"
	"time"
)

func TestRegisterLookupAndExpire(t *testing.T) {
	s := NewInMemory()
	ctx := context.Background()
	if err := s.Register(ctx, "alice", "sip:alice@1.2.3.4", 1); err != nil {
		t.Fatalf("register: %v", err)
	}
	got, err := s.Lookup(ctx, "alice")
	if err != nil || len(got) != 1 {
		t.Fatalf("lookup after register: %v len=%d", err, len(got))
	}
	time.Sleep(1500 * time.Millisecond)
	got, err = s.Lookup(ctx, "alice")
	if err != nil || len(got) != 0 {
		t.Fatalf("lookup after expire: %v len=%d", err, len(got))
	}
}

func TestDeregister(t *testing.T) {
	s := NewInMemory()
	ctx := context.Background()
	_ = s.Register(ctx, "bob", "sip:bob@1.2.3.4", 60)
	_ = s.Register(ctx, "bob", "sip:bob@5.6.7.8", 60)
	if err := s.Deregister(ctx, "bob", "sip:bob@1.2.3.4"); err != nil {
		t.Fatalf("deregister: %v", err)
	}
	got, _ := s.Lookup(ctx, "bob")
	if len(got) != 1 || got[0].Contact != "sip:bob@5.6.7.8" {
		t.Fatalf("unexpected bindings: %+v", got)
	}
}