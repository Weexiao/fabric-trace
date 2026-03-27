package storage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// IPFSClient uploads/downloads files via IPFS HTTP API.
type IPFSClient struct {
	apiURL       string
	maxSizeBytes int64
	httpClient   *http.Client
}

func NewIPFSClient(apiURL string, maxSizeBytes int64) *IPFSClient {
	return &IPFSClient{
		apiURL:       strings.TrimRight(apiURL, "/"),
		maxSizeBytes: maxSizeBytes,
		httpClient:   &http.Client{},
	}
}

var errSizeLimitExceeded = errors.New("ipfs upload size limit exceeded")

type hashCountWriter struct {
	hash         hash.Hash
	maxSizeBytes int64
	size         int64
}

func (w *hashCountWriter) Write(p []byte) (int, error) {
	w.size += int64(len(p))
	if w.maxSizeBytes > 0 && w.size > w.maxSizeBytes {
		return 0, fmt.Errorf("%w: file exceeds limit %d bytes", errSizeLimitExceeded, w.maxSizeBytes)
	}
	_, _ = w.hash.Write(p)
	return len(p), nil
}

// Put uploads data to IPFS add endpoint in streaming mode.
// Hash and size are calculated on the uploaded byte stream.
func (c *IPFSClient) Put(ctx context.Context, r io.Reader, sizeHint int64) (cid string, size int64, hashHex string, err error) {
	_ = sizeHint
	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	counter := &hashCountWriter{hash: sha256.New(), maxSizeBytes: c.maxSizeBytes}

	go func() {
		defer func() {
			_ = mw.Close()
			_ = pw.Close()
		}()

		fw, e := mw.CreateFormFile("file", "blob")
		if e != nil {
			_ = pw.CloseWithError(e)
			return
		}

		tee := io.TeeReader(r, counter)
		if _, e = io.Copy(fw, tee); e != nil {
			_ = pw.CloseWithError(e)
			return
		}
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL+"/api/v0/add?pin=true", pr)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, errSizeLimitExceeded) {
			err = fmt.Errorf("file exceeds limit %d bytes", c.maxSizeBytes)
		}
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("ipfs add failed: status %d, body %s", resp.StatusCode, string(body))
		return
	}

	var addResp struct {
		Hash string `json:"Hash"`
		Size string `json:"Size"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&addResp); err != nil {
		return
	}
	cid = addResp.Hash
	size = counter.size
	hashHex = hex.EncodeToString(counter.hash.Sum(nil))
	return
}

// Get fetches content by CID via IPFS cat and returns buffered bytes reader for hash校验/解密。
func (c *IPFSClient) Get(ctx context.Context, cid string) (io.ReadCloser, int64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL+"/api/v0/cat?arg="+cid, nil)
	if err != nil {
		return nil, 0, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return nil, 0, fmt.Errorf("ipfs cat failed: status %d, body %s", resp.StatusCode, string(body))
	}
	data, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, 0, err
	}
	return io.NopCloser(bytes.NewReader(data)), int64(len(data)), nil
}

// LoadAESKey loads a 32-byte key from env.
func LoadAESKey(envName string) ([]byte, error) {
	key := os.Getenv(envName)
	if len(key) == 0 {
		return nil, fmt.Errorf("env %s is empty", envName)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("env %s must be 32 bytes for AES-256-GCM", envName)
	}
	return []byte(key), nil
}
