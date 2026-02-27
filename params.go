// Package fees provides zero-dependency, pure-function utilities for
// calculating Cardano transaction fees and minimum UTxO (minADA) values
// in the Babbage and Conway eras, per CIP-55.
//
// All calculations are deterministic and require no network connection.
// Supply your own ProtocolParams from any Cardano API (Blockfrost, Maestro,
// Ogmios, cardano-cli) and this library does the rest.
//
// CIP-55 reference: https://cips.cardano.org/cip/CIP-55
// Ledger spec:      https://github.com/intersectmbo/cardano-ledger
package fees

// ProtocolParams holds the subset of Cardano protocol parameters needed
// for fee and minUTxO calculations. All fields use Lovelace as the unit
// unless noted otherwise.
//
// Obtain current values from your preferred Cardano API. For mainnet as of
// early 2025, typical values are shown in DefaultMainnetParams().
type ProtocolParams struct {
	// MinFeeA is the coefficient applied to transaction size in bytes.
	// Also called txFeePerByte or a in the fee formula: fee = a*size + b.
	// Mainnet: 44
	MinFeeA uint64

	// MinFeeB is the constant term added to the fee.
	// Also called txFeeFixed or b in the fee formula: fee = a*size + b.
	// Mainnet: 155381
	MinFeeB uint64

	// CoinsPerUTxOByte is the cost per byte of UTxO storage (Babbage/Conway).
	// Replaces the deprecated coinsPerUTxOWord parameter.
	// The minUTxO formula is: (160 + serializedOutputBytes) * CoinsPerUTxOByte.
	// Mainnet: 4310
	CoinsPerUTxOByte uint64

	// MaxTxSize is the maximum allowed transaction size in bytes.
	// Mainnet: 16384
	MaxTxSize uint64
}

// DefaultMainnetParams returns ProtocolParams populated with typical Cardano
// mainnet values as of the Conway era (early 2025). Always fetch live params
// from a Cardano node or API for production use — these may change via governance.
//
// Example:
//
//	p := fees.DefaultMainnetParams()
//	f ee := fees.MinFee(p, 300)
func DefaultMainnetParams() ProtocolParams {
	return ProtocolParams{
		MinFeeA:          44,
		MinFeeB:          155381,
		CoinsPerUTxOByte: 4310,
		MaxTxSize:        16384,
	}
}

// DefaultPreviewParams returns ProtocolParams for the Cardano preview testnet.
// Values may differ from mainnet; always verify against live protocol parameters.
//
// Example:
//
//	p := fees.DefaultPreviewParams()
func DefaultPreviewParams() ProtocolParams {
	return ProtocolParams{
		MinFeeA:          44,
		MinFeeB:          155381,
		CoinsPerUTxOByte: 4310,
		MaxTxSize:        16384,
	}
}

// Validate checks that ProtocolParams contain plausible non-zero values.
// Returns a non-nil error if any required field is zero.
//
// Example:
//
//	p := fees.ProtocolParams{MinFeeA: 44}
//	if err := p.Validate(); err != nil {
//		log.Fatal(err)
//	}
func (p ProtocolParams) Validate() error {
	if p.MinFeeA == 0 {
		return &ParamError{Field: "MinFeeA", Message: "must be non-zero"}
	}
	if p.MinFeeB == 0 {
		return &ParamError{Field: "MinFeeB", Message: "must be non-zero"}
	}
	if p.CoinsPerUTxOByte == 0 {
		return &ParamError{Field: "CoinsPerUTxOByte", Message: "must be non-zero"}
	}
	if p.MaxTxSize == 0 {
		return &ParamError{Field: "MaxTxSize", Message: "must be non-zero"}
	}
	return nil
}

// ParamError is returned when a ProtocolParams field is invalid.
type ParamError struct {
	// Field is the name of the invalid parameter.
	Field string
	// Message describes why the value is invalid.
	Message string
}

func (e *ParamError) Error() string {
	return "fees: invalid protocol param " + e.Field + ": " + e.Message
}