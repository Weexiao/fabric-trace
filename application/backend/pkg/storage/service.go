package storage

import (
	"backend/pkg"
	"backend/settings"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"strings"
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
	if settings.Cfg.Crypto.Enabled {
		return s.uploadEncrypted(ctx, traceID, fileID, mime, r)
	}
	return s.uploadPlaintextStream(ctx, traceID, fileID, mime, r)
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

func shouldCompressForEvidence(size int64) bool {
	cfg := settings.Cfg.Compression
	if !cfg.Enabled {
		return false
	}
	minBytes := cfg.MinSizeBytes
	maxBytes := cfg.MaxSizeBytes
	if minBytes <= 0 {
		minBytes = 5 * 1024 * 1024
	}
	if maxBytes <= 0 {
		maxBytes = 5 * 1024 * 1024 * 1024
	}
	if maxBytes < minBytes {
		maxBytes = minBytes
	}
	return size >= minBytes && size <= maxBytes
}

func compressForEvidence(data []byte) ([]byte, string, error) {
	alg := strings.ToLower(strings.TrimSpace(settings.Cfg.Compression.Algorithm))
	if alg == "" {
		alg = "gzip"
	}
	switch alg {
	case "btae":
		// Model compression pipeline is not integrated yet; fallback keeps behavior deterministic.
		compressed, err := pkg.GzipCompress(data)
		return compressed, "btae_fallback_gzip", err
	default:
		compressed, err := pkg.GzipCompress(data)
		return compressed, "gzip", err
	}
}

func (s *Service) uploadEncrypted(ctx context.Context, traceID, fileID, mime string, r io.Reader) (*Manifest, error) {
	plain, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	sourceHash := pkg.SHA256Hex(plain)
	compressedHash := ""
	compressAlg := ""
	if shouldCompressForEvidence(int64(len(plain))) {
		compressed, alg, err := compressForEvidence(plain)
		if err != nil {
			return nil, err
		}
		compressedHash = pkg.SHA256Hex(compressed)
		compressAlg = alg
	}

	encRes, err := pkg.EncryptAndHash(bytes.NewReader(plain), s.key)
	if err != nil {
		return nil, err
	}
	toStore, err := io.ReadAll(encRes.Reader())
	if err != nil {
		return nil, err
	}

	cid, size, storedHash, err := s.ipfs.Put(ctx, bytes.NewReader(toStore), int64(len(toStore)))
	if err != nil {
		return nil, err
	}

	return &Manifest{
		TraceabilityCode: traceID,
		FileID:           fileID,
		CID:              cid,
		Hash:             storedHash,
		SourceHash:       sourceHash,
		CompressedHash:   compressedHash,
		CompressAlg:      compressAlg,
		Mime:             mime,
		Size:             size,
		Encrypted:        true,
		KeyVersion:       settings.Cfg.Crypto.KeyVersion,
	}, nil
}

func (s *Service) uploadPlaintextStream(ctx context.Context, traceID, fileID, mime string, r io.Reader) (*Manifest, error) {
	sourceHasher := sha256.New()
	writerForHash := io.Writer(sourceHasher)

	compressedHash := ""
	compressAlg := ""
	var gzipWriter *gzip.Writer
	var compressedHasher hash.Hash
	if settings.Cfg.Compression.Enabled {
		compressedHasher = sha256.New()
		gzipWriter = gzip.NewWriter(compressedHasher)
		writerForHash = io.MultiWriter(sourceHasher, gzipWriter)
		compressAlg = resolveCompressionAlgorithmLabel()
	}

	teed := io.TeeReader(r, writerForHash)
	cid, size, storedHash, err := s.ipfs.Put(ctx, teed, -1)
	if gzipWriter != nil {
		if closeErr := gzipWriter.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}
	if err != nil {
		return nil, err
	}

	sourceHash := hex.EncodeToString(sourceHasher.Sum(nil))
	if shouldCompressForEvidence(size) && compressedHasher != nil {
		compressedHash = hex.EncodeToString(compressedHasher.Sum(nil))
	} else {
		compressAlg = ""
	}

	return &Manifest{
		TraceabilityCode: traceID,
		FileID:           fileID,
		CID:              cid,
		Hash:             storedHash,
		SourceHash:       sourceHash,
		CompressedHash:   compressedHash,
		CompressAlg:      compressAlg,
		Mime:             mime,
		Size:             size,
		Encrypted:        false,
		KeyVersion:       "",
	}, nil
}

func resolveCompressionAlgorithmLabel() string {
	alg := strings.ToLower(strings.TrimSpace(settings.Cfg.Compression.Algorithm))
	switch alg {
	case "btae":
		return "btae_fallback_gzip"
	default:
		return "gzip"
	}
}
