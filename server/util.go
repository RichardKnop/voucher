package server

import (
	"context"
	"net/http"
	"time"

	"github.com/RichardKnop/voucher/server/response"
)

var timeout = 5 * time.Second

type handlerFunc func(ctx context.Context) (interface{}, int64, error)

func processHandler(f handlerFunc, w http.ResponseWriter, r *http.Request, successCode int) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, httpErrCode, err := f(ctx)
	if err != nil {
		response.Error(w, err.Error(), int(httpErrCode))
	} else {
		response.WriteJSON(w, result, successCode)
	}
}
