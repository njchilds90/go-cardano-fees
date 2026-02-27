# Changelog

All notable changes to `go-cardano-fees` are documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] — 2026-02-24

### Added

- `ProtocolParams` struct with `Validate()`, `DefaultMainnetParams()`, `DefaultPreviewParams()`
- `ParamError` structured error type
- `MinFee(params, txSizeBytes)` — exact linear fee formula
- `MinFeeWithPadding(params, txSizeBytes, paddingBytes)` — fee with witness padding
- `EstimateFee(params, numInputs, numOutputs, hasMetadata)` — structural estimate without serialization
- `FeeError` structured error type
- `OutputSize` struct for structural minUTxO estimation
- `MinUTxO(params, OutputSize)` — structural minUTxO estimator
- `MinUTxOFromBytes(params, serializedBytes)` — exact CIP-55 formula
- `EstimateOutputBytes(OutputSize)` — byte-size estimator for TxOut components
- `MinUTxOADAOnly(params)` — quick helper for ADA-only outputs
- `MinUTxOForNFT(params, assetNameLen)` — quick helper for single NFT outputs
- `MinUTxOForBundle(params, numPolicies, numAssets, totalAssetNameBytes)` — multi-asset bundle
- `IsAboveMinUTxO(params, lovelace, OutputSize)` — preflight guard
- `MinUTxOError` structured error type
- `LovelacePerADA` constant
- `ToLovelace(ada float64)` — ADA → Lovelace conversion
- `ToADA(lovelace uint64)` — Lovelace → ADA conversion
- `FormatADA(lovelace uint64)` — human-readable ADA string
- `FormatLovelace(lovelace uint64)` — human-readable Lovelace string
- `AddLovelace(a, b uint64)` — overflow-safe addition
- `SumLovelace([]uint64)` — overflow-safe sum
- Full table-driven test suite with race detector support
- GitHub Actions CI across Go 1.21, 1.22, 1.23
- GoDoc examples on all exported functions
- MIT license