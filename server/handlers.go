package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RichardKnop/voucher/service"
)

// VoucherHandler ...
type IndexHandler struct {
	service service.IFace
}

// VoucherHandler ...
type VoucherHandler struct {
	service service.IFace
}

func (h *VoucherHandler) handleGetIndex(w http.ResponseWriter, r *http.Request) {
	f := func(ctx context.Context) (interface{}, error) {
		return "index", nil
	}

	processHandler(f, w, r)
}

func (h *VoucherHandler) handleGetSpecific(w http.ResponseWriter, r *http.Request, voucherID int) {
	f := func(ctx context.Context) (interface{}, error) {
		return fmt.Sprintf("voucherID: %d", voucherID), nil
	}

	processHandler(f, w, r)
}

func (h *VoucherHandler) handlePost(w http.ResponseWriter, r *http.Request, data []byte) {
	f := func(ctx context.Context) (interface{}, error) {
		return fmt.Sprintf("POST data: %s", data), nil
	}

	processHandler(f, w, r)
}
