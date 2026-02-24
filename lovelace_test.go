package fees_test

import (
	"testing"

	fees "github.com/njchilds90/go-cardano-fees"
)

func TestToLovelace(t *testing.T) {
	tests := []struct {
		name    string
		ada     float64
		want    uint64
		wantErr bool
	}{
		{"1 ADA", 1.0, 1_000_000, false},
		{"1.5 ADA", 1.5, 1_500_000, false},
		{"0.001 ADA", 0.001, 1_000, false},
		{"0 ADA", 0.0, 0, false},
		{"negative", -1.0, 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fees.ToLovelace(tc.ada)
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
				t.Errorf("ToLovelace(%f) = %d, want %d", tc.ada, got, tc.want)
			}
		})
	}
}

func TestToADA(t *testing.T) {
	tests := []struct {
		lovelace uint64
		want     float64
	}{
		{1_000_000, 1.0},
		{1_500_000, 1.5},
		{0, 0.0},
	}

	for _, tc := range tests {
		got := fees.ToADA(tc.lovelace)
		if got != tc.want {
			t.Errorf("ToADA(%d) = %f, want %f", tc.lovelace, got, tc.want)
		}
	}
}

func TestFormatADA(t *testing.T) {
	got := fees.FormatADA(1_310_000)
	want := "1.310000 ADA"
	if got != want {
		t.Errorf("FormatADA(1310000) = %q, want %q", got, want)
	}
}

func TestFormatLovelace(t *testing.T) {
	got := fees.FormatLovelace(1_310_000)
	want := "1310000 Lovelace"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAddLovelace(t *testing.T) {
	got, err := fees.AddLovelace(1_000_000, 500_000)
	if err != nil {
		t.Fatal(err)
	}
	if got != 1_500_000 {
		t.Errorf("got %d, want 1500000", got)
	}
}

func TestSumLovelace(t *testing.T) {
	vals := []uint64{1_000_000, 2_000_000, 500_000}
	got, err := fees.SumLovelace(vals)
	if err != nil {
		t.Fatal(err)
	}
	if got != 3_500_000 {
		t.Errorf("got %d, want 3500000", got)
	}
}

func TestSumLovelaceEmpty(t *testing.T) {
	got, err := fees.SumLovelace(nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != 0 {
		t.Errorf("empty sum should be 0, got %d", got)
	}
}

func TestLovelacePerADA(t *testing.T) {
	if fees.LovelacePerADA != 1_000_000 {
		t.Errorf("LovelacePerADA should be 1000000, got %d", fees.LovelacePerADA)
	}
}
