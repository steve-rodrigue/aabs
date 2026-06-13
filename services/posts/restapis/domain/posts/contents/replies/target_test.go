package replies

import "testing"

func TestTargetWithReply(t *testing.T) {
	reply := NewMockReply("Parent reply")

	target := &target{
		reply: reply,
	}

	if !target.IsReply() {
		t.Fatalf("expected target to be reply")
	}

	if target.Reply() != reply {
		t.Fatalf("expected reply")
	}

	if target.IsThread() {
		t.Fatalf("expected target not to be thread")
	}

	if target.Thread() != nil {
		t.Fatalf("expected nil thread")
	}
}

func TestTargetWithThread(t *testing.T) {
	thread := newMockThread()

	target := &target{
		thread: thread,
	}

	if target.IsReply() {
		t.Fatalf("expected target not to be reply")
	}

	if target.Reply() != nil {
		t.Fatalf("expected nil reply")
	}

	if !target.IsThread() {
		t.Fatalf("expected target to be thread")
	}

	if target.Thread() != thread {
		t.Fatalf("expected thread")
	}
}
