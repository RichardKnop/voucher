package server

import (
	"net/http"
	"context"
	"io/ioutil"
	"strconv"
	"fmt"
)

// ServeHTTP ...
func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f := func(ctx context.Context) (interface{}, error) {
		return "Welcome to voucher service", nil
	}

	processHandler(f, w, r)
}

// ServeHTTP ...
func (h *VoucherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)

    var (
    	singleResource = len(head) > 0
    	voucherID int 
    	err error
    )

    if singleResource {
	    voucherID, err = strconv.Atoi(head)
	    if err != nil {
	        http.Error(w, fmt.Sprintf("Invalid voucher ID %s", head), http.StatusBadRequest)
	        return
	    }
    }
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