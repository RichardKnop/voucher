package service

import (
	"testing"
)

func TestValidateVoucherID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		voucherID string
		valid     bool
	}{
		{
			voucherID: "foo",
			valid:     true,
		},
		{
			voucherID: "foo123",
			valid:     true,
		},
		{
			voucherID: "foo_123",
			valid:     false,
		},
	}

	for _, testCase := range testCases {
		err := ValidateVoucherID(testCase.voucherID)
		if testCase.valid && err != nil {
			t.Fatalf("expected no error but got \"%v\"", err)
		}
		if !testCase.valid && err == nil {
			t.Fatalf("expected error but got nil")
		}
	}
}
