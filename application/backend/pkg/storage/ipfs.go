package storage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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

// Put uploads data to IPFS add endpoint. Payload is buffered (<=限制)，hash/size基于完整字节。
func (c *IPFSClient) Put(ctx context.Context, r io.Reader, sizeHint int64) (cid string, size int64, hashHex string, err error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return
	}
	size = int64(len(data))
	if c.maxSizeBytes > 0 && size > c.maxSizeBytes {
		err = fmt.Errorf("file exceeds limit %d bytes", c.maxSizeBytes)
		return
	}
	h := sha256.Sum256(data)
	hashHex = hex.EncodeToString(h[:])

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL+"/api/v0/add?pin=true", bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

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
	if addResp.Size != "" {
		if parsed, e := strconv.ParseInt(addResp.Size, 10, 64); e == nil {
			size = parsed
		}
	}
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
		resp.Body.Close()
		return nil, 0, fmt.Errorf("ipfs cat failed: status %d, body %s", resp.StatusCode, string(body))
	}
	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
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
