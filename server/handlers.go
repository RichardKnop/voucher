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
		offsetParam := r.URL.Query().Get("offset")
		offset, _ := strconv.Atoi(offsetParam)

		countParam := r.URL.Query().Get("count")
		count, _ := strconv.Atoi(countParam)

		vouchers, nextOffset, httpErrCode, err := h.svc.FindAll(int64(offset), int64(count))
		if err != nil {
			return nil, httpErrCode, err
		}
		return map[string]interface{}{
			"offset":     offset,
			"nextOffset": nextOffset,
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
