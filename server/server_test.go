package server

import (
	"testing"
)

func TestShiftPath(t *testing.T) {
	t.Parallel()

	head, tail := shiftPath("")
	if head != "" {
		t.Fatalf("head expected \"\", instead got \"%s\"", head)
	}
	if tail != "/" {
		t.Fatalf("tail expected \"/\", instead got \"%s\"", tail)
	}

	head, tail = shiftPath("vouchers")
	if head != "vouchers" {
		t.Fatalf("head expected \"vouchers\", instead got \"%s\"", head)
	}
	if tail != "/" {
		t.Fatalf("tail expected \"/\", instead got \"%s\"", tail)
	}

	head, tail = shiftPath("vouchers/")
	if head != "vouchers" {
		t.Fatalf("head expected \"vouchers\", instead got \"%s\"", head)
	}
	if tail != "/" {
		t.Fatalf("tail expected \"/\", instead got \"%s\"", tail)
	}

	head, tail = shiftPath("/vouchers/")
	if head != "vouchers" {
		t.Fatalf("head expected \"vouchers\", instead got \"%s\"", head)
	}
	if tail != "/" {
		t.Fatalf("tail expected \"/\", instead got \"%s\"", tail)
	}

	head, tail = shiftPath("vouchers/foo")
	if head != "vouchers" {
		t.Fatalf("head expected \"vouchers\", instead got \"%s\"", head)
	}
	if tail != "/foo" {
		t.Fatalf("tail expected \"/foo\", instead got \"%s\"", tail)
	}

	head, tail = shiftPath("vouchers/foo/")
	if head != "vouchers" {
		t.Fatalf("head expected \"vouchers\", instead got \"%s\"", head)
	}
	if tail != "/foo" {
		t.Fatalf("tail expected \"/foo\", instead got \"%s\"", tail)
	}
}
