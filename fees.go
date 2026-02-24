package fees

import "fmt"

// MinFee calculates the minimum transaction fee in Lovelace using the
// Cardano linear fee formula:
//
//	fee = MinFeeA * txSizeBytes + MinFeeB
//
// txSizeBytes is the size of the serialized transaction in bytes.
// Use an estimated or actual CBOR-serialized size. A common safe approach
// is to build the transaction with a dummy fee, measure its byte length,
// then recalculate.
//
// Returns an error if params are invalid or txSizeBytes exceeds MaxTxSize.
//
// Example:
//
//	p := fees.DefaultMainnetParams()
//	lovelace, err := fees.MinFee(p, 350)
//	// lovelace ≈ 170,781
func MinFee(p ProtocolParams, txSizeBytes uint64) (uint64, error) {
	if err := p.Validate(); err != nil {
		return 0, err
	}
	if txSizeBytes == 0 {
		return 0, &FeeError{Reason: "txSizeBytes must be greater than zero"}
	}
	if txSizeBytes > p.MaxTxSize {
		return 0, &FeeError{
			Reason: fmt.Sprintf("txSizeBytes %d exceeds MaxTxSize %d", txSizeBytes, p.MaxTxSize),
		}
	}
	return p.MinFeeA*txSizeBytes + p.MinFeeB, nil
}

// MinFeeWithPadding calculates the minimum transaction fee and adds a safety
// buffer expressed as a number of additional bytes. This is useful when the
// true serialized size is not yet known.
//
// A common pattern is to add 100–200 bytes of padding for witness data.
//
// Example:
//
//	p := fees.DefaultMainnetParams()
//	fee, err := fees.MinFeeWithPadding(p, 300, 150)
//	// fee = MinFee(p, 450)
func MinFeeWithPadding(p ProtocolParams, txSizeBytes, paddingBytes uint64) (uint64, error) {
	return MinFee(p, txSizeBytes+paddingBytes)
}

// EstimateFee provides a quick fee estimate given the number of transaction
// inputs, outputs, and whether a metadata payload is present. It uses a
// simple byte-size model calibrated against mainnet transactions.
//
// This is an approximation useful for UI display and pre-flight checks.
// For exact fees, serialize the full transaction and use MinFee.
//
// Example:
//
//	p := fees.DefaultMainnetParams()
//	fee, err := fees.EstimateFee(p, 2, 2, true)
func EstimateFee(p ProtocolParams, numInputs, numOutputs uint64, hasMetadata bool) (uint64, error) {
	if err := p.Validate(); err != nil {
		return 0, err
	}
	if numInputs == 0 {
		return 0, &FeeError{Reason: "numInputs must be at least 1"}
	}
	if numOutputs == 0 {
		return 0, &FeeError{Reason: "numOutputs must be at least 1"}
	}

	// Empirically-derived byte model:
	//   base tx overhead:  ~200 bytes
	//   per input:         ~40 bytes (TxIn hash+index + Vkey witness ~100 bytes)
	//   per output:        ~65 bytes (address + value)
	//   metadata overhead: ~250 bytes estimate
	const (
		baseTxSize    uint64 = 200
		bytesPerInput uint64 = 140 // includes witness
		bytesPerOutput uint64 = 65
		metadataSize  uint64 = 250
	)

	estimated := baseTxSize + bytesPerInput*numInputs + bytesPerOutput*numOutputs
	if hasMetadata {
		estimated += metadataSize
	}

	return MinFee(p, estimated)
}

// FeeError is returned when a fee calculation cannot be completed.
type FeeError struct {
	// Reason describes why the calculation failed.
	Reason string
}

func (e *FeeError) Error() string {
	return "fees: " + e.Reason
}
