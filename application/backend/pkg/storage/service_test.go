package storage

import (
	"backend/settings"
	"testing"
)

func TestShouldCompressForEvidenceBoundaries(t *testing.T) {
	oldCfg := settings.Cfg
	defer func() { settings.Cfg = oldCfg }()

	settings.Cfg.Compression.Enabled = true
	settings.Cfg.Compression.MinSizeBytes = 5 * 1024 * 1024
	settings.Cfg.Compression.MaxSizeBytes = 5 * 1024 * 1024 * 1024

	cases := []struct {
		name string
		size int64
		want bool
	}{
		{name: "below min", size: 5*1024*1024 - 1, want: false},
		{name: "at min", size: 5 * 1024 * 1024, want: true},
		{name: "middle", size: 100 * 1024 * 1024, want: true},
		{name: "at max", size: 5 * 1024 * 1024 * 1024, want: true},
		{name: "above max", size: 5*1024*1024*1024 + 1, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := shouldCompressForEvidence(tc.size); got != tc.want {
				t.Fatalf("shouldCompressForEvidence(%d)=%v, want %v", tc.size, got, tc.want)
			}
		})
	}
}

func TestResolveCompressionAlgorithmLabel(t *testing.T) {
	oldCfg := settings.Cfg
	defer func() { settings.Cfg = oldCfg }()

	settings.Cfg.Compression.Algorithm = "gzip"
	if got := resolveCompressionAlgorithmLabel(); got != "gzip" {
		t.Fatalf("got %s, want gzip", got)
	}

	settings.Cfg.Compression.Algorithm = "btae"
	if got := resolveCompressionAlgorithmLabel(); got != "btae_fallback_gzip" {
		t.Fatalf("got %s, want btae_fallback_gzip", got)
	}
}
