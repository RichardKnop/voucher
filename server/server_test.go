package server

import (
	"testing"
)

func TestShiftPath(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		path         string
		expectedHead string
		expectedTail string
	}{
		{
			path:         "",
			expectedHead: "",
			expectedTail: "/",
		},
		{
			path:         "vouchers",
			expectedHead: "vouchers",
			expectedTail: "/",
		},
		{
			path:         "vouchers/",
			expectedHead: "vouchers",
			expectedTail: "/",
		},
		{
			path:         "/vouchers/",
			expectedHead: "vouchers",
			expectedTail: "/",
		},
		{
			path:         "vouchers/foo",
			expectedHead: "vouchers",
			expectedTail: "/foo",
		},
		{
			path:         "vouchers/foo/",
			expectedHead: "vouchers",
			expectedTail: "/foo",
		},
	}

	for _, testCase := range testCases {
		head, tail := shiftPath(testCase.path)
		if head != testCase.expectedHead {
			t.Fatalf("head expected \"%s\", instead got \"%s\"", testCase.expectedHead, head)
		}
		if tail != testCase.expectedTail {
			t.Fatalf("tail expected \"%s\", instead got \"%s\"", testCase.expectedTail, tail)
		}
	}
}
