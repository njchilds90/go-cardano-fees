package fees

import (
	"fmt"
	"math"
)

const (
	// LovelacePerADA is the number of Lovelace in one ADA.
	LovelacePerADA uint64 = 1_000_000
)

// ToLovelace converts an ADA amount (as a float64) to Lovelace (uint64),
// truncating any sub-Lovelace fractional remainder.
//
// Returns an error if ada is negative or would overflow uint64.
//
// Example:
//
//	lv, err := fees.ToLovelace(1.5)   // 1_500_000
//	lv, err := fees.ToLovelace(0.001) // 1_000
func ToLovelace(ada float64) (uint64, error) {
	if ada < 0 {
		return 0, fmt.Errorf("fees: ToLovelace: ada must be non-negative, got %f", ada)
	}
	// Compute the maximum ADA value that can be represented without overflow.
	// maxAda = MaxUint64 / LovelacePerADA
	maxAda := float64(math.MaxUint64) / float64(LovelacePerADA)
	if ada > maxAda {
		return 0, fmt.Errorf("fees: ToLovelace: value %f overflows uint64", ada)
	}
	result := ada * float64(LovelacePerADA)
	return uint64(result), nil
}

// ToADA converts a Lovelace amount to ADA as a float64.
//
// Example:
//
//	ada := fees.ToADA(1_500_000) // 1.5
func ToADA(lovelace uint64) float64 {
	return float64(lovelace) / float64(LovelacePerADA)
}

// FormatADA formats a Lovelace amount as a human-readable ADA string
// with 6 decimal places.
//
// Example:
//
//	fees.FormatADA(1_310_000) // "1.310000 ADA"
func FormatADA(lovelace uint64) string {
	return fmt.Sprintf("%.6f ADA", ToADA(lovelace))
}

// FormatLovelace formats a uint64 Lovelace value as a string with the
// unit suffix for display purposes.
//
// Example:
//
//	fees.FormatLovelace(1_310_000) // "1310000 Lovelace"
func FormatLovelace(lovelace uint64) string {
	return fmt.Sprintf("%d Lovelace", lovelace)
}

// AddLovelace safely adds two Lovelace values, returning an error on overflow.
//
// Example:
//
//	total, err := fees.AddLovelace(1_000_000, 500_000) // 1_500_000
func AddLovelace(a, b uint64) (uint64, error) {
	if a > math.MaxUint64-b {
		return 0, fmt.Errorf("fees: AddLovelace: overflow adding %d + %d", a, b)
	}
	return a + b, nil
}

// SumLovelace adds a slice of Lovelace values, returning an error on overflow.
//
// Example:
//
//	total, err := fees.SumLovelace([]uint64{1_000_000, 2_000_000, 500_000})
func SumLovelace(values []uint64) (uint64, error) {
	var total uint64
	for _, v := range values {
		var err error
		total, err = AddLovelace(total, v)
		if err != nil {
			return 0, err
		}
	}
	return total, nil
}

// IsAboveMinUTxO returns true if the given Lovelace amount meets or exceeds
// the calculated minUTxO for the described output.
//
// Example:
//
//	p := fees.DefaultMainnetParams()
//	ok, required, err := fees.IsAboveMinUTxO(p, 2_000_000, fees.OutputSize{AddressBytes: 57})
//	// ok=true, required≈969750
func IsAboveMinUTxO(p ProtocolParams, lovelace uint64, out OutputSize) (bool, uint64, error) {
	required, err := MinUTxO(p, out)
	if err != nil {
		return false, 0, err
	}
	return lovelace >= required, required, nil
}