package storage

// Manifest mirrors the on-chain FileManifest for convenience in backend.
type Manifest struct {
	TraceabilityCode string `json:"traceabilityCode"`
	FileID           string `json:"fileID"`
	CID              string `json:"cid"`
	Hash             string `json:"hash"`
	SourceHash       string `json:"sourceHash,omitempty"`
	CompressedHash   string `json:"compressedHash"`
	CompressedBits   []int  `json:"compressedBits,omitempty"`
	CompressAlg      string `json:"compressAlg"`
	Mime             string `json:"mime"`
	Size             int64  `json:"size"`
	Encrypted        bool   `json:"encrypted"`
	KeyVersion       string `json:"keyVersion"`
	Role             string `json:"role"`
	Uploader         string `json:"uploader"`
	Timestamp        string `json:"timestamp"`
}
