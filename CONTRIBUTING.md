# Contributing to go-cardano-fees

Thank you for your interest in contributing. This library aims to stay
small, focused, and zero-dependency. Please read this before submitting.

---

## Scope

This library does **one thing**: calculate Cardano transaction fees and
minimum UTxO values from protocol parameters. Contributions should stay
within this scope. Serialization, signing, and node interaction belong
in other libraries.

**In scope:**
- Bug fixes in fee/minUTxO calculations
- Accuracy improvements with ledger spec citations
- New era support (e.g., future parameter changes via governance)
- Additional helper functions for common output patterns
- Documentation and example improvements

**Out of scope:**
- External dependencies of any kind
- CBOR serialization (bring your own)
- Network or node interaction
- Wallet or key management

---

## How to Contribute

1. **Open an issue first** for any non-trivial change to discuss the approach.
2. Fork the repository and create a branch: `git checkout -b fix/my-fix`
3. Write or update tests. All changes must maintain or improve coverage.
4. Run `go test -race ./...` and `go vet ./...` locally (or verify via CI).
5. Submit a pull request with a clear description referencing the issue.

---

## Standards

- Zero external runtime dependencies — this is a hard constraint
- All exported functions must have GoDoc comments with an Example
- Tests must be table-driven where multiple cases apply
- Error types must be structured (`*ParamError`, `*FeeError`, `*MinUTxOError`)
- Ledger spec or CIP references required for any formula changes

---

## Reporting Calculation Errors

If you find a discrepancy between this library's output and `cardano-cli transaction calculate-min-fee` or on-chain behavior, please open an issue with:

- The exact protocol parameters used
- The input values
- The expected result (with source: CLI output, ledger spec section, etc.)
- The actual result from this library

---

## License

By contributing, you agree your contributions are licensed under the MIT License.
```

---

## Release & Verification Instructions
```
RELEASE STEPS — go-cardano-fees v1.0.0
========================================

1. VERIFY CI IS GREEN
   - Go to: https://github.com/njchilds90/go-cardano-fees/actions
   - Confirm the latest CI run on main shows all green checks across Go 1.21/1.22/1.23
   - Do not release if any check is failing

2. CREATE THE TAG
   - Go to: https://github.com/njchilds90/go-cardano-fees/releases/new
   - In "Choose a tag" → type:  v1.0.0  → click "Create new tag: v1.0.0 on publish"
   - Target branch: main

3. FILL IN THE RELEASE FORM
   - Title:  v1.0.0 — Initial Release
   - Description (paste this):

     ## go-cardano-fees v1.0.0

     Zero-dependency Go library for Cardano transaction fee and minimum UTxO
     calculation (Babbage/Conway era, CIP-55 compliant).

     ### What's included
     - `MinFee` — exact linear fee formula (a × size + b)
     - `MinFeeWithPadding` — fee with witness byte padding
     - `EstimateFee` — structural estimate, no serialization needed
     - `MinUTxO` / `MinUTxOFromBytes` — CIP-55 minADA formula
     - `MinUTxOADAOnly` / `MinUTxOForNFT` / `MinUTxOForBundle` — common-case helpers
     - `IsAboveMinUTxO` — preflight submission guard
     - `ToLovelace` / `ToADA` / `FormatADA` / `SumLovelace` — Lovelace utilities
     - Structured error types: `*ParamError`, `*FeeError`, `*MinUTxOError`
     - Full table-driven test suite with race detector
     - CI across Go 1.21, 1.22, 1.23

     ### Companion library
     Pair with [go-cardano-metadata](https://github.com/njchilds90/go-cardano-metadata)
     for a complete NFT minting toolkit.

     ### Install
```
     go get github.com/njchilds90/go-cardano-fees@v1.0.0
```

4. CHECK "Set as the latest release"
5. Click "Publish release"

6. VERIFY pkg.go.dev INDEXING
   - Wait approximately 10 minutes
   - Visit: https://pkg.go.dev/github.com/njchilds90/go-cardano-fees
   - You should see the full GoDoc with all exported types and functions
   - If it hasn't appeared, manually trigger: https://pkg.go.dev/github.com/njchilds90/go-cardano-fees@v1.0.0

SEMANTIC VERSIONING GUIDANCE
==============================
v1.0.x  Bug fixes, calculation accuracy improvements (no API changes)
v1.1.0  New helper functions, new era support, new error types (backwards compatible)
v2.0.0  Only if API breaks (e.g., ProtocolParams fields renamed, return types changed)
