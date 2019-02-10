package server

import (
	"net/http"
	"path"
	"strings"
)

// Server ...
type Server struct {
	indexHandler   *IndexHandler
	voucherHandler *VoucherHandler
}

// New ...
func New() *Server {
	return &Server{
		indexHandler:   new(IndexHandler),
		voucherHandler: new(VoucherHandler),
	}
}

// ServeHTTP ...
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "" {
		srv.indexHandler.ServeHTTP(w, r)
		return
	} else if head == "vouchers" {
		srv.voucherHandler.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Not Found", http.StatusNotFound)
}

// shiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
