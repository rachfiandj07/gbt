package middleware

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type (
	Chain func(httprouter.Handle) httprouter.Handle
)

// ChainReq is a middleware for.....
func ChainReq(endHandler httprouter.Handle, chains ...Chain) httprouter.Handle {
	if len(chains) == 0 {
		return endHandler
	}

	return chains[0](ChainReq(endHandler, chains[1:]...))
}

var InitContext = func(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		r = r.WithContext(context.Background())
		next(w, r, p)
	}
}

// SetHeader is for add response header for common JSON api
var SetHeader = func(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,OPTION")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie, Source-Type, Origin, Content-Filename")
		next(w, r, p)
	}
}
