package fees_test

import (
	"errors"
	"testing"

	fees "github.com/njchilds90/go-cardano-fees"
)

func TestMinFee(t *testing.T) {
	p := fees.DefaultMainnetParams()

	tests := []struct {
		name        string
		txSize      uint64
		wantFee     uint64
		wantErrType interface{}
	}{
		{
			name:    "typical simple tx 300 bytes",
			txSize:  300,
			wantFee: 44*300 + 155381, // 168581
		},
		{
			name:    "minimum 1 byte tx",
			txSize:  1,
			wantFee: 44*1 + 155381, // 155425
		},
		{
			name:    "max tx size",
			txSize:  16384,
			wantFee: 44*16384 + 155381,
		},
		{
			name:        "zero bytes",
			txSize:      0,
			wantErrType: &fees.FeeError{},
		},
		{
			name:        "exceeds max tx size",
			txSize:      20000,
			wantErrType: &fees.FeeError{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fees.MinFee(p, tc.txSize)
			if tc.wantErrType != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantFee {
				t.Errorf("MinFee(%d) = %d, want %d", tc.txSize, got, tc.wantFee)
			}
		})
	}
}

func TestMinFeeWithPadding(t *testing.T) {
	p := fees.DefaultMainnetParams()
	fee, err := fees.MinFeeWithPadding(p, 200, 100)
	if err != nil {
		t.Fatal(err)
	}
	expected, _ := fees.MinFee(p, 300)
	if fee != expected {
		t.Errorf("got %d, want %d", fee, expected)
	}
}

func TestEstimateFee(t *testing.T) {
	p := fees.DefaultMainnetParams()

	tests := []struct {
		name        string
		inputs      uint64
		outputs     uint64
		hasMeta     bool
		wantErr     bool
	}{
		{"1in 1out no meta", 1, 1, false, false},
		{"2in 2out with meta", 2, 2, true, false},
		{"0 inputs", 0, 1, false, true},
		{"0 outputs", 1, 0, false, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fee, err := fees.EstimateFee(p, tc.inputs, tc.outputs, tc.hasMeta)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if fee == 0 {
				t.Error("fee should not be zero")
			}
			// Sanity: fee should be at least MinFeeB
			if fee < p.MinFeeB {
				t.Errorf("fee %d is below MinFeeB %d", fee, p.MinFeeB)
			}
		})
	}
}

func TestInvalidParams(t *testing.T) {
	bad := fees.ProtocolParams{} // all zero
	_, err := fees.MinFee(bad, 300)
	if err == nil {
		t.Fatal("expected error for zero params")
	}
	var pe *fees.ParamError
	if !errors.As(err, &pe) {
		t.Errorf("expected *ParamError, got %T", err)
	}
}

func TestParamValidate(t *testing.T) {
	p := fees.DefaultMainnetParams()
	if err := p.Validate(); err != nil {
		t.Errorf("DefaultMainnetParams should be valid: %v", err)
	}

	p2 := fees.ProtocolParams{MinFeeA: 44} // missing rest
	if err := p2.Validate(); err == nil {
		t.Error("expected error for incomplete params")
	}
}

func TestDefaultPreviewParams(t *testing.T) {
	p := fees.DefaultPreviewParams()
	if err := p.Validate(); err != nil {
		t.Errorf("DefaultPreviewParams should be valid: %v", err)
	}
}
