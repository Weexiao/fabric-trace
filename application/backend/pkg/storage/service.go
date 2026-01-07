package storage

import (
	"backend/pkg"
	"backend/settings"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// Service combines crypto and IPFS storage.
type Service struct {
	ipfs *IPFSClient
	key  []byte
}

func NewService() (*Service, error) {
	cfg := settings.Cfg
	if cfg.Storage.Type != "ipfs" {
		return nil, fmt.Errorf("unsupported storage type %s", cfg.Storage.Type)
	}

	var key []byte
	if cfg.Crypto.Enabled {
		k, err := LoadAESKey(cfg.Crypto.KeyEnv)
		if err != nil {
			return nil, err
		}
		key = k
	}

	if cfg.Storage.IPFS.APIURL == "" {
		return nil, fmt.Errorf("ipfs api_url is empty")
	}

	ipfs := NewIPFSClient(cfg.Storage.IPFS.APIURL, cfg.Storage.MaxSizeMB*1024*1024)
	return &Service{ipfs: ipfs, key: key}, nil
}

// Upload encrypts, uploads to IPFS, and returns manifest fields.
func (s *Service) Upload(ctx context.Context, traceID, fileID, mime string, r io.Reader) (*Manifest, error) {
	if !settings.Cfg.Crypto.Enabled {
		// No encryption: store plaintext directly
		plain, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		cid, size, _, err := s.ipfs.Put(ctx, bytes.NewReader(plain), int64(len(plain)))
		if err != nil {
			return nil, err
		}
		m := &Manifest{
			TraceabilityCode: traceID,
			FileID:           fileID,
			CID:              cid,
			Hash:             pkg.SHA256Hex(plain),
			Mime:             mime,
			Size:             size,
			Encrypted:        false,
			KeyVersion:       "",
		}
		return m, nil
	}

	encRes, err := pkg.EncryptAndHash(r, s.key)
	if err != nil {
		return nil, err
	}
	cid, size, hashHex, err := s.ipfs.Put(ctx, encRes.Reader(), encRes.Size())
	if err != nil {
		return nil, err
	}
	m := &Manifest{
		TraceabilityCode: traceID,
		FileID:           fileID,
		CID:              cid,
		Hash:             hashHex,
		Mime:             mime,
		Size:             size,
		Encrypted:        true,
		KeyVersion:       settings.Cfg.Crypto.KeyVersion,
	}
	return m, nil
}

// Download fetches from IPFS, validates hash, and decrypts.
func (s *Service) Download(ctx context.Context, manifest Manifest) (io.ReadCloser, int64, error) {
	rc, sizeCipher, err := s.ipfs.Get(ctx, manifest.CID)
	if err != nil {
		return nil, 0, err
	}
	data, err := io.ReadAll(rc)
	_ = rc.Close()
	if err != nil {
		return nil, 0, err
	}
	if manifest.Hash != "" {
		// verify hash
		if hashHex := pkg.SHA256Hex(data); hashHex != manifest.Hash {
			return nil, 0, fmt.Errorf("hash mismatch: expected %s got %s", manifest.Hash, hashHex)
		}
	}

	if !manifest.Encrypted || !settings.Cfg.Crypto.Enabled {
		// no encryption
		return io.NopCloser(bytes.NewReader(data)), sizeCipher, nil
	}

	dec, err := pkg.DecryptReader(bytes.NewReader(data), s.key)
	if err != nil {
		return nil, 0, err
	}
	return io.NopCloser(dec.Reader()), sizeCipher, nil
}

// MarshalManifestJSON helper for chaincode invocation.
func MarshalManifestJSON(m *Manifest) (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
