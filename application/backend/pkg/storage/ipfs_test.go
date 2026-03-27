package storage

import (
	"backend/pkg"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIPFSPutStreamSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/v0/add") {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if _, err := io.Copy(io.Discard, r.Body); err != nil {
			t.Fatalf("read request body failed: %v", err)
		}
		_ = r.Body.Close()

		_ = json.NewEncoder(w).Encode(map[string]string{
			"Hash": "QmTestCID",
			"Size": "999",
		})
	}))
	defer ts.Close()

	payload := []byte("streamed ipfs payload")
	client := NewIPFSClient(ts.URL, 1024)

	cid, size, hashHex, err := client.Put(context.Background(), strings.NewReader(string(payload)), int64(len(payload)))
	if err != nil {
		t.Fatalf("put failed: %v", err)
	}
	if cid != "QmTestCID" {
		t.Fatalf("cid=%s, want QmTestCID", cid)
	}
	if size != int64(len(payload)) {
		t.Fatalf("size=%d, want %d", size, len(payload))
	}
	if hashHex != pkg.SHA256Hex(payload) {
		t.Fatalf("hash=%s, want %s", hashHex, pkg.SHA256Hex(payload))
	}
}

func TestIPFSPutSizeLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"Hash": "QmIgnored"})
	}))
	defer ts.Close()

	client := NewIPFSClient(ts.URL, 4)
	_, _, _, err := client.Put(context.Background(), strings.NewReader("12345"), 5)
	if err == nil {
		t.Fatal("expected size limit error, got nil")
	}
	if !strings.Contains(err.Error(), "exceeds limit") {
		t.Fatalf("unexpected error: %v", err)
	}
}
