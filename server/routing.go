package server

import (
	"context"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/RichardKnop/voucher/server/response"
	"github.com/RichardKnop/voucher/service"
)

var (
	isAlpha = regexp.MustCompile(`^[A-Za-z]+$`).MatchString
)

// ServeHTTP ...
func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f := func(ctx context.Context) (interface{}, int64, error) {
		return "Welcome to voucher service", 0, nil
	}

	processHandler(f, w, r, 200)
}

// ServeHTTP ...
func (h *VoucherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)

	var (
		singleResource = len(head) > 0
		voucherID      = head
	)

	switch r.Method {
	case "GET":
		if singleResource {
			if err := service.ValidateVoucherID(voucherID); err != nil {
				response.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			h.handleGetSpecific(w, r, voucherID)
		} else {
			h.handleGetIndex(w, r)
		}
	case "POST":
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			response.Error(w, "failed to read request data", http.StatusInternalServerError)
			return
		}
		h.handlePost(w, r, data)
	case "DELETE":
		if singleResource {
			if err := service.ValidateVoucherID(voucherID); err != nil {
				response.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			h.handleDelete(w, r, voucherID)
		} else {
			response.Error(w, "HTTP method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		response.Error(w, "HTTP method not allowed", http.StatusMethodNotAllowed)
	}
}
