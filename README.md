# go-cardano-fees

[![CI](https://github.com/njchilds90/go-cardano-fees/actions/workflows/ci.yml/badge.svg)](https://github.com/njchilds90/go-cardano-fees/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/njchilds90/go-cardano-fees.svg)](https://pkg.go.dev/github.com/njchilds90/go-cardano-fees)
[![Go Report Card](https://goreportcard.com/badge/github.com/njchilds90/go-cardano-fees)](https://goreportcard.com/report/github.com/njchilds90/go-cardano-fees)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**Zero-dependency Go library for Cardano transaction fee and minimum UTxO (minADA) calculation — Babbage/Conway era, CIP-55 compliant.**

Built for backend minting services, AI agents, wallets, and tools that need to know exactly how much ADA to attach to an output — with no node connection required.

---

## Why This Exists

Every Cardano developer sending native tokens must answer two questions before submitting a transaction:

1. **What is the minimum transaction fee?** (`fee = a × txSizeBytes + b`)
2. **What is the minimum ADA (minUTxO) I must include with this output?** (`(160 + outputBytes) × coinsPerUTxOByte`)

Python has `PyCardano.utils.min_lovelace()`. JavaScript has `min_ada_required()` in the serialization lib. Go has nothing standalone — the math is buried in heavy SDKs with node dependencies.

This library is the Go answer: pure functions, zero dependencies, table-driven, fully testable offline.

---

## Companion Library

Pair with [`go-cardano-metadata`](https://github.com/njchilds90/go-cardano-metadata) to build **and** price your NFT outputs in one session:

- `go-cardano-metadata` → build CIP-25/CIP-68 compliant metadata payloads
- `go-cardano-fees` → calculate the ADA required to send them on-chain

---

## Installation
```bash
go get github.com/njchilds90/go-cardano-fees@v1.0.0
```

---

## Quick Start
```go
package main

import (
    "fmt"
    fees "github.com/njchilds90/go-cardano-fees"
)

func main() {
    p := fees.DefaultMainnetParams()

    // Minimum transaction fee for a 350-byte transaction
    fee, err := fees.MinFee(p, 350)
    if err != nil {
        panic(err)
    }
    fmt.Println("Min fee:", fees.FormatADA(fee)) // "0.170781 ADA"

    // Min ADA for an ADA-only output
    minADA, _ := fees.MinUTxOADAOnly(p)
    fmt.Println("Min ADA (ada-only):", fees.FormatADA(minADA))

    // Min ADA for a single NFT output (32-byte asset name)
    nftMin, _ := fees.MinUTxOForNFT(p, 32)
    fmt.Println("Min ADA (NFT):", fees.FormatADA(nftMin))

    // Quick fee estimate without serializing
    estimated, _ := fees.EstimateFee(p, 2, 2, true)
    fmt.Println("Estimated fee (2in/2out/metadata):", fees.FormatADA(estimated))
}
```

---

## API Reference

### Protocol Parameters
```go
// Use mainnet defaults (verify against live API in production)
p := fees.DefaultMainnetParams()

// Or supply your own from Blockfrost / Maestro / Ogmios
p := fees.ProtocolParams{
    MinFeeA:          44,
    MinFeeB:          155381,
    CoinsPerUTxOByte: 4310,
    MaxTxSize:        16384,
}

// Validate before use
if err := p.Validate(); err != nil {
    log.Fatal(err)
}
```

### Transaction Fee Calculation
```go
// Exact fee from serialized tx size
fee, err := fees.MinFee(p, txSizeBytes)

// Fee with padding for unknown witness size
fee, err := fees.MinFeeWithPadding(p, baseTxBytes, 150)

// Quick structural estimate (no serialization needed)
fee, err := fees.EstimateFee(p, numInputs, numOutputs, hasMetadata)
```

### Minimum UTxO (minADA)
```go
// ADA-only output (standard Shelley address)
minADA, err := fees.MinUTxOADAOnly(p)

// Single NFT (1 policy, 1 asset, N-byte name)
minADA, err := fees.MinUTxOForNFT(p, assetNameLenBytes)

// Token bundle
minADA, err := fees.MinUTxOForBundle(p, numPolicies, numAssets, totalAssetNameBytes)

// Custom output with full structural control
minADA, err := fees.MinUTxO(p, fees.OutputSize{
    AddressBytes:        57,
    NumPolicies:         2,
    NumAssets:           5,
    TotalAssetNameBytes: 80,
    HasDatumHash:        true,
})

// From exact CBOR-serialized byte count
minADA, err := fees.MinUTxOFromBytes(p, 225)
```

### Lovelace Utilities
```go
lv, err  := fees.ToLovelace(1.5)        // 1_500_000
ada      := fees.ToADA(1_500_000)       // 1.5
str      := fees.FormatADA(1_500_000)   // "1.500000 ADA"
str      := fees.FormatLovelace(1_500_000) // "1500000 Lovelace"
sum, err := fees.SumLovelace([]uint64{...})

// Guard: check if an output has enough ADA before submission
ok, required, err := fees.IsAboveMinUTxO(p, myLovelace, outputSize)
```

---

## CIP-55 Formula

The Babbage/Conway era minUTxO formula (CIP-55):
```
minUTxO = (160 + |serialized TxOut in bytes|) × coinsPerUTxOByte
```

This replaced the older `coinsPerUTxOWord` formula from the Mary era. `go-cardano-fees` implements both the structural estimator and the exact-byte variant so you can use whichever fits your pipeline.

---

## For AI Agents

This library is designed to be called deterministically in automated minting pipelines:
```go
type MintRequest struct {
    NumInputs   uint64
    NumOutputs  uint64
    AssetName   string
    HasMetadata bool
}

func preflight(req MintRequest, p fees.ProtocolParams) error {
    fee, err := fees.EstimateFee(p, req.NumInputs, req.NumOutputs, req.HasMetadata)
    if err != nil {
        return fmt.Errorf("fee estimation failed: %w", err)
    }

    minADA, err := fees.MinUTxOForNFT(p, uint64(len(req.AssetName)))
    if err != nil {
        return fmt.Errorf("minUTxO calculation failed: %w", err)
    }

    total, err := fees.SumLovelace([]uint64{fee, minADA})
    if err != nil {
        return fmt.Errorf("overflow: %w", err)
    }

    fmt.Printf("Required: %s fee + %s minADA = %s total\n",
        fees.FormatADA(fee),
        fees.FormatADA(minADA),
        fees.FormatADA(total),
    )
    return nil
}
```

---

## Error Types

All errors are typed and unwrappable:

| Type | When |
|---|---|
| `*ParamError` | Invalid `ProtocolParams` field |
| `*FeeError` | Invalid input to fee calculation |
| `*MinUTxOError` | Invalid input to minUTxO calculation |

---

## Versioning

This library follows [Semantic Versioning](https://semver.org):

- `v1.x.x` — Babbage/Conway era support (current)
- `v2.x.x` — future: exact CBOR sizing built-in (if needed)

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## Changelog

See [CHANGELOG.md](CHANGELOG.md).

---

## License

MIT — see [LICENSE](LICENSE).