package fees_test

import (
	"testing"

	fees "github.com/njchilds90/go-cardano-fees"
)

func TestMinUTxOFromBytes(t *testing.T) {
	p := fees.DefaultMainnetParams()

	tests := []struct {
		name    string
		bytes   uint64
		want    uint64
		wantErr bool
	}{
		{
			name:  "base shelley output ~100 bytes",
			bytes: 100,
			want:  (160 + 100) * 4310, // 1_120_600
		},
		{
			name:  "larger output 200 bytes",
			bytes: 200,
			want:  (160 + 200) * 4310, // 1_551_600
		},
		{
			name:    "zero bytes",
			bytes:   0,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fees.MinUTxOFromBytes(p, tc.bytes)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestMinUTxOADAOnly(t *testing.T) {
	p := fees.DefaultMainnetParams()
	got, err := fees.MinUTxOADAOnly(p)
	if err != nil {
		t.Fatal(err)
	}
	// ADA-only Shelley output is roughly 0.9–1.1 ADA
	adaVal := fees.ToADA(got)
	if adaVal < 0.8 || adaVal > 1.5 {
		t.Errorf("MinUTxOADAOnly = %.6f ADA, expected between 0.8 and 1.5", adaVal)
	}
}

func TestMinUTxOForNFT(t *testing.T) {
	p := fees.DefaultMainnetParams()

	tests := []struct {
		name         string
		assetNameLen uint64
		wantErr      bool
	}{
		{"empty asset name", 0, false},
		{"short name 8 bytes", 8, false},
		{"max name 32 bytes", 32, false},
		{"over max 33 bytes", 33, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fees.MinUTxOForNFT(p, tc.assetNameLen)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			// NFT minUTxO should be above ADA-only minimum
			adaOnly, _ := fees.MinUTxOADAOnly(p)
			if got <= adaOnly {
				t.Errorf("NFT minUTxO %d should exceed ADA-only %d", got, adaOnly)
			}
		})
	}
}

func TestMinUTxOForBundle(t *testing.T) {
	p := fees.DefaultMainnetParams()

	tests := []struct {
		name                string
		numPolicies         uint64
		numAssets           uint64
		totalAssetNameBytes uint64
		wantErr             bool
	}{
		{"1 policy 1 asset", 1, 1, 10, false},
		{"2 policies 5 assets", 2, 5, 80, false},
		{"0 policies", 0, 1, 10, true},
		{"0 assets", 1, 0, 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fees.MinUTxOForBundle(p, tc.numPolicies, tc.numAssets, tc.totalAssetNameBytes)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got == 0 {
				t.Error("minUTxO should not be zero")
			}
		})
	}
}

func TestMinUTxOBundleGrowth(t *testing.T) {
	// More assets = higher minUTxO
	p := fees.DefaultMainnetParams()

	min1, _ := fees.MinUTxOForBundle(p, 1, 1, 10)
	min5, _ := fees.MinUTxOForBundle(p, 1, 5, 50)
	min10, _ := fees.MinUTxOForBundle(p, 2, 10, 100)

	if min1 >= min5 {
		t.Errorf("expected min1 < min5, got %d >= %d", min1, min5)
	}
	if min5 >= min10 {
		tt.Errorf("expected min5 < min10, got %d >= %d", min5, min10)
	}
}

func TestIsAboveMinUTxO(t *testing.T) {
	p := fees.DefaultMainnetParams()
	out := fees.OutputSize{AddressBytes: 57}

	// Should be above min with 2 ADA
	two_ada := uint64(2_000_000)
	ok, required, err := fees.IsAboveMinUTxO(p, two_ada, out)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("2 ADA (%d) should be above minUTxO (%d)", two_ada, required)
	}

	// Should be below min with 100 Lovelace
	ok2, _, err := fees.IsAboveMinUTxO(p, 100, out)
	if err != nil {
		t.Fatal(err)
	}
	if ok2 {
		t.Error("100 Lovelace should NOT be above minUTxO")
	}
}

func TestEstimateOutputBytes(t *testing.T) {
	// ADA-only output should be smaller than NFT output
	adaOnly := fees.EstimateOutputBytes(fees.OutputSize{AddressBytes: 57})
	nft := fees.EstimateOutputBytes(fees.OutputSize{
		AddressBytes:        57,
		NumPolicies:         1,
		NumAssets:           1,
		TotalAssetNameBytes: 16,
	})
	if adaOnly >= nft {
		t.Errorf("ADA-only (%d bytes) should be smaller than NFT (%d bytes)", adaOnly, nft)
	}
}