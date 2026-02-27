package fees

import "fmt"

// OutputSize describes a transaction output for minUTxO calculation purposes.
// It models the byte-size contribution of the output's components.
//
// For the Babbage/Conway era, the formula per CIP-55 is:
//
//	minUTxO = (160 + serializedOutputBytes) * CoinsPerUTxOByte
//
// where serializedOutputBytes is the CBOR-serialized size of the full TxOut.
// This library provides both an exact-byte variant (MinUTxOFromBytes) and a
// structural estimator (EstimateOutputBytes / MinUTxO).
type OutputSize struct {
	// AddressBytes is the byte length of the output address.
	// Shelley base addresses are typically 57 bytes.
	// Enterprise addresses are typically 29 bytes.
	AddressBytes uint64

	// NumPolicies is the number of distinct policy IDs in the token bundle.
	NumPolicies uint64

	// NumAssets is the total number of distinct asset names across all policies.
	NumAssets uint64

	// TotalAssetNameBytes is the sum of byte lengths of all asset names.
	// Asset names may be 0–32 bytes.
	TotalAssetNameBytes uint64

	// HasDatumHash indicates whether the output includes a datum hash (32 bytes).
	HasDatumHash bool

	// HasInlineDatum indicates whether the output includes an inline datum.
	// If true, InlineDatumBytes must be set.
	HasInlineDatum bool

	// InlineDatumBytes is the byte length of the inline datum, if present.
	InlineDatumBytes uint64

	// HasScriptRef indicates whether the output contains a reference script.
	// If true, ScriptRefBytes must be set.
	HasScriptRef bool

	// ScriptRefBytes is the byte length of the reference script, if present.
	ScriptRefBytes uint64
}

// MinUTxO calculates the minimum ADA (in Lovelace) that must be included
// in a transaction output for the Babbage/Conway era using CIP-55's formula:
//
//	minUTxO = (160 + serializedOutputBytes) * CoinsPerUTxOByte
//
// The serialized byte size is estimated from the OutputSize structure.
// For exact sizing, use MinUTxOFromBytes with the actual CBOR size.
//
// Example (ADA-only output, Shelley address):
//
//	p := fees.DefaultMainnetParams()
//	minADA, err := fees.MinUTxO(p, fees.OutputSize{AddressBytes: 57})
//	// minADA ≈ 969,750 Lovelace (≈ 0.97 ADA)
//
// Example (single NFT output):
//
//	minADA, err := fees.MinUTxO(p, fees.OutputSize{
//		AddressBytes:        57,
//		NumPolicies:         1,
//		NumAssets:           1,
//		TotalAssetNameBytes: 32,
//	})
//	// minADA ≈ 1,310,000 Lovelace (≈ 1.31 ADA)
func MinUTxO(p ProtocolParams, out OutputSize) (uint64, error) {
	if err := p.Validate(); err != nil {
		return 0, err
	}
	serialized := EstimateOutputBytes(out)
	return MinUTxOFromBytes(p, serialized)
}

// MinUTxOFromBytes calculates the minimum ADA in Lovelace using the
// exact serialized byte count of the TxOut. Use this when you have
// already CBOR-serialized the output and measured its length.
//
// Formula (CIP-55 / Babbage+):
//
//	minUTxO = (160 + serializedOutputBytes) * CoinsPerUTxOByte
//
// Example:
//
//	p := fees.DefaultMainnetParams()
//	minADA, err := fees.MinUTxOFromBytes(p, 125)
//	// minADA = (160 + 125) * 4310 = 1,228,350
func MinUTxOFromBytes(p ProtocolParams, serializedOutputBytes uint64) (uint64, error) {
	if err := p.Validate(); err != nil {
		return 0, err
	}
	if serializedOutputBytes == 0 {
		return 0, &MinUTxOError{Reason: "serializedOutputBytes must be greater than zero"}
	}
	return (160 + serializedOutputBytes) * p.CoinsPerUTxOByte, nil
}

// EstimateOutputBytes estimates the serialized CBOR byte size of a TxOut
// from its structural components. This is based on the Mary/Babbage era
// ledger size model used in the minUTxO specification.
//
// The estimate uses the following model:
//   - TxOut envelope overhead:    ~10 bytes
//   - Address:                    AddressBytes
//   - ADA value (Lovelace only):  ~9 bytes (CBOR uint64)
//   - Per-policy overhead:        28 bytes (policy hash) per policy
//   - Per-asset-name:             12 bytes overhead per asset
//   - Asset name content:         TotalAssetNameBytes
//   - Integer value per asset:    ~5 bytes per asset
//   - Token bundle structure:     ~5 bytes if any assets
//   - Datum hash:                 32 bytes if HasDatumHash
//   - Inline datum:               InlineDatumBytes if HasInlineDatum
//   - Script reference:           ScriptRefBytes if HasScriptRef
//
// Example:
//
//	size := fees.EstimateOutputBytes(fees.OutputSize{
//		AddressBytes:        57,
//		NumPolicies:         1,
//		NumAssets:           1,
//		TotalAssetNameBytes: 9,
//	})
func EstimateOutputBytes(out OutputSize) uint64 {
	const (
		envelopeOverhead  uint64 = 10
		adaValueBytes     uint64 = 9
		policyHashBytes   uint64 = 28
		perAssetOverhead  uint64 = 12
		perAssetIntBytes  uint64 = 5
		tokenBundleFixed  uint64 = 5
		datumHashBytes    uint64 = 32
	)

	total := envelopeOverhead + out.AddressBytes + adaValueBytes

	if out.NumAssets > 0 || out.NumPolicies > 0 {
		total += tokenBundleFixed
		total += policyHashBytes * out.NumPolicies
		total += perAssetOverhead * out.NumAssets
		total += out.TotalAssetNameBytes
		total += perAssetIntBytes * out.NumAssets
	}

	if out.HasDatumHash {
		total += datumHashBytes
	}

	if out.HasInlineDatum {
		total += out.InlineDatumBytes
	}

	if out.HasScriptRef {
		total += out.ScriptRefBytes
	}

	return total
}

// MinUTxOADAOnly returns the minimum Lovelace for a simple ADA-only output
// with a standard Shelley base address (57 bytes). This is the most common
// case and a useful quick reference.
//
// Example:
//
//	p := fees.DefaultMainnetParams()
//	minADA, err := fees.MinUTxOADAOnly(p)
func MinUTxOADAOnly(p ProtocolParams) (uint64, error) {
	return MinUTxO(p, OutputSize{AddressBytes: 57})
}

// MinUTxOForNFT returns the minimum Lovelace for a UTxO holding a single
// NFT (one policy, one asset) with a standard Shelley base address.
//
// assetNameLen is the byte length of the asset name (0–32).
//
// Example:
//
//	p := fees.DefaultMainnetParams()
//	minADA, err := fees.MinUTxOForNFT(p, 32) // asset name is 32 bytes
func MinUTxOForNFT(p ProtocolParams, assetNameLen uint64) (uint64, error) {
	if assetNameLen > 32 {
		return 0, &MinUTxOError{
			Reason: fmt.Sprintf("assetNameLen %d exceeds maximum of 32 bytes", assetNameLen),
		}
	}
	return MinUTxO(p, OutputSize{
		AddressBytes:        57,
		NumPolicies:         1,
		NumAssets:           1,
		TotalAssetNameBytes: assetNameLen,
	})
}

// MinUTxOForBundle returns the minimum Lovelace for a UTxO holding a
// bundle of native tokens across one or more policies.
//
// Example:
//
//	// 2 policies, 5 total assets, asset names summing to 80 bytes
//	p := fees.DefaultMainnetParams()
//	minADA, err := fees.MinUTxOForBundle(p, 2, 5, 80)
func MinUTxOForBundle(p ProtocolParams, numPolicies, numAssets, totalAssetNameBytes uint64) (uint64, error) {
	if numPolicies == 0 {
		return 0, &MinUTxOError{Reason: "numPolicies must be at least 1"}
	}
	if numAssets == 0 {
		return 0, &MinUTxOError{Reason: "numAssets must be at least 1"}
	}
	return MinUTxO(p, OutputSize{
		AddressBytes:        57,
		NumPolicies:         numPolicies,
		NumAssets:           numAssets,
		TotalAssetNameBytes: totalAssetNameBytes,
	})
}

// MinUTxOError is returned when a minUTxO calculation cannot be completed.
type MinUTxOError struct {
	// Reason describes why the calculation failed.
	Reason string
}

func (e *MinUTxOError) Error() string {
	return "fees: minUTxO: " + e.Reason
}