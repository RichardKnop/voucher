package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/RichardKnop/voucher/service"
)

// VoucherHandler ...
type IndexHandler struct {
	svc service.IFace
}

// VoucherHandler ...
type VoucherHandler struct {
	svc service.IFace
}

func (h *VoucherHandler) handleGetIndex(w http.ResponseWriter, r *http.Request) {
	f := func(ctx context.Context) (interface{}, int64, error) {
		cursorParam := r.URL.Query().Get("cursor")
		cursor, _ := strconv.Atoi(cursorParam)
		vouchers, nextCursor, httpErrCode, err := h.svc.FindAll(uint64(cursor))
		if err != nil {
			return nil, httpErrCode, err
		}
		return map[string]interface{}{
			"cursor":     cursor,
			"nextCursor": nextCursor,
			"vouchers":   vouchers,
		}, 0, nil
	}

	processHandler(f, w, r, 200)
}

func (h *VoucherHandler) handleGetSpecific(w http.ResponseWriter, r *http.Request, voucherID string) {
	f := func(ctx context.Context) (interface{}, int64, error) {
		return h.svc.FindByID(voucherID)
	}

	processHandler(f, w, r, 200)
}

func (h *VoucherHandler) handlePost(w http.ResponseWriter, r *http.Request, data []byte) {
	f := func(ctx context.Context) (interface{}, int64, error) {
		return h.svc.Create(data)
	}

	processHandler(f, w, r, 201)
}
