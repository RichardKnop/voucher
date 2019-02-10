package server

import (
	"context"
	"net/http"
	"time"

	"github.com/RichardKnop/voucher/server/response"
)

var timeout = 5 * time.Second

type handlerFunc func(ctx context.Context) (interface{}, error)

func processHandler(f handlerFunc, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := f(ctx)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		response.WriteJSON(w, result, 200)
	}
}
