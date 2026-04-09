package sync

import (
	"fmt"

	"github.com/envoy-cli/internal/crypto"
	"github.com/envoy-cli/internal/env"
	"github.com/envoy-cli/internal/storage"
)

// Syncer handles pushing and pulling .env files between local and remote backends.
type Syncer struct {
	local     storage.Backend
	remote    storage.Backend
	encryptor *crypto.Encryptor
	parser    *env.Parser
	writer    *env.Writer
}

// NewSyncer creates a new Syncer with the provided backends and encryptor.
func NewSyncer(local, remote storage.Backend, encryptor *crypto.Encryptor) *Syncer {
	return &Syncer{
		local:     local,
		remote:    remote,
		encryptor: encryptor,
		parser:    env.NewParser(),
		writer:    env.NewWriter(),
	}
}

// Push reads the local .env file, encrypts it, and writes it to the remote backend.
func (s *Syncer) Push(name string) error {
	data, err := s.local.Read(name)
	if err != nil {
		return fmt.Errorf("push: read local: %w", err)
	}

	encrypted, err := s.encryptor.Encrypt(data)
	if err != nil {
		return fmt.Errorf("push: encrypt: %w", err)
	}

	if err := s.remote.Write(name, encrypted); err != nil {
		return fmt.Errorf("push: write remote: %w", err)
	}

	return nil
}

// Pull reads the remote .env file, decrypts it, and writes it to the local backend.
func (s *Syncer) Pull(name string) error {
	data, err := s.remote.Read(name)
	if err != nil {
		return fmt.Errorf("pull: read remote: %w", err)
	}

	decrypted, err := s.encryptor.Decrypt(data)
	if err != nil {
		return fmt.Errorf("pull: decrypt: %w", err)
	}

	if err := s.local.Write(name, decrypted); err != nil {
		return fmt.Errorf("pull: write local: %w", err)
	}

	return nil
}

// Diff returns keys that differ between the local and remote (decrypted) env files.
func (s *Syncer) Diff(name string) ([]string, error) {
	localData, err := s.local.Read(name)
	if err != nil {
		return nil, fmt.Errorf("diff: read local: %w", err)
	}

	remoteEncrypted, err := s.remote.Read(name)
	if err != nil {
		return nil, fmt.Errorf("diff: read remote: %w", err)
	}

	remoteData, err := s.encryptor.Decrypt(remoteEncrypted)
	if err != nil {
		return nil, fmt.Errorf("diff: decrypt remote: %w", err)
	}

	localMap, err := s.parser.ParseBytes(localData)
	if err != nil {
		return nil, fmt.Errorf("diff: parse local: %w", err)
	}

	remoteMap, err := s.parser.ParseBytes(remoteData)
	if err != nil {
		return nil, fmt.Errorf("diff: parse remote: %w", err)
	}

	var diffKeys []string
	for k, lv := range localMap {
		if rv, ok := remoteMap[k]; !ok || rv != lv {
			diffKeys = append(diffKeys, k)
		}
	}
	for k := range remoteMap {
		if _, ok := localMap[k]; !ok {
			diffKeys = append(diffKeys, k)
		}
	}

	return diffKeys, nil
}
