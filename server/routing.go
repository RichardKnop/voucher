package server

import (
	"context"
	"io/ioutil"
	"net/http"
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
		voucherID      = head // todo: validate voucher ID
	)

	switch r.Method {
	case "GET":
		if singleResource {
			h.handleGetSpecific(w, r, voucherID)
		} else {
			h.handleGetIndex(w, r)
		}
	case "POST":
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request data", http.StatusInternalServerError)
			return
		}
		h.handlePost(w, r, data)
	default:
		http.Error(w, "HTTP method not allowed", http.StatusMethodNotAllowed)
	}
}
