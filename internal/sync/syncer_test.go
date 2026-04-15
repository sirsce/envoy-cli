package sync_test

import (
	"testing"

	"github.com/envoy-cli/internal/crypto"
	"github.com/envoy-cli/internal/storage"
	"github.com/envoy-cli/internal/sync"
)

func newTestSyncer(t *testing.T) (*sync.Syncer, storage.Backend, storage.Backend) {
	t.Helper()

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	enc, err := crypto.NewEncryptor(key)
	if err != nil {
		t.Fatalf("new encryptor: %v", err)
	}

	local := storage.NewMemoryBackend()
	remote := storage.NewMemoryBackend()

	return sync.NewSyncer(local, remote, enc), local, remote
}

func TestSyncer_PushPull(t *testing.T) {
	syncer, local, remote := newTestSyncer(t)

	original := []byte("APP_ENV=production\nDEBUG=false\n")
	if err := local.Write("app", original); err != nil {
		t.Fatalf("write local: %v", err)
	}

	if err := syncer.Push("app"); err != nil {
		t.Fatalf("push: %v", err)
	}

	// Verify remote data is encrypted (not equal to original)
	remoteData, err := remote.Read("app")
	if err != nil {
		t.Fatalf("read remote: %v", err)
	}
	if string(remoteData) == string(original) {
		t.Error("expected remote data to be encrypted, got plaintext")
	}

	// Clear local and pull
	if err := local.Write("app", []byte{}); err != nil {
		t.Fatalf("clear local: %v", err)
	}

	if err := syncer.Pull("app"); err != nil {
		t.Fatalf("pull: %v", err)
	}

	pulled, err := local.Read("app")
	if err != nil {
		t.Fatalf("read local after pull: %v", err)
	}
	if string(pulled) != string(original) {
		t.Errorf("expected %q, got %q", original, pulled)
	}
}

func TestSyncer_Diff_NoDiff(t *testing.T) {
	syncer, local, _ := newTestSyncer(t)

	data := []byte("KEY=value\n")
	if err := local.Write("app", data); err != nil {
		t.Fatalf("write local: %v", err)
	}
	if err := syncer.Push("app"); err != nil {
		t.Fatalf("push: %v", err)
	}

	diffs, err := syncer.Diff("app")
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	if len(diffs) != 0 {
		t.Errorf("expected no diff, got %v", diffs)
	}
}

func TestSyncer_Diff_WithChanges(t *testing.T) {
	syncer, local, _ := newTestSyncer(t)

	original := []byte("KEY=value\nOTHER=foo\n")
	if err := local.Write("app", original); err != nil {
		t.Fatalf("write local: %v", err)
	}
	if err := syncer.Push("app"); err != nil {
		t.Fatalf("push: %v", err)
	}

	// Modify local
	if err := local.Write("app", []byte("KEY=changed\nNEW=bar\n")); err != nil {
		t.Fatalf("write local: %v", err)
	}

	diffs, err := syncer.Diff("app")
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	if len(diffs) == 0 {
		t.Error("expected diffs, got none")
	}
}

func TestSyncer_Pull_NotFound(t *testing.T) {
	syncer, _, _ := newTestSyncer(t)

	// Pulling an environment that has never been pushed should return an error.
	if err := syncer.Pull("nonexistent"); err == nil {
		t.Error("expected error when pulling nonexistent environment, got nil")
	}
}
