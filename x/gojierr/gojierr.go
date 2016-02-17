package gojierr

import (
	"log"
	"net/http"

	"github.com/crewjam/httperr"
	"github.com/zenazn/goji/web"
	"golang.org/x/net/context"
)

type contextKeyType int

const (
	urlParamsKey contextKeyType = iota
)

// Param returns the named URL parameter or an empty string if it is not present
func Param(ctx context.Context, name string) string {
	rv, _ := ctx.Value(urlParamsKey).(map[string]string)[name]
	return rv
}

// NewContextFunc is invoked to renerate the new context for
// a request. The default implementation uses context.Background,
// but when running in appengine this must be replaced with
// appengine.NewContext(r).
var NewContextFunc = func(r *http.Request) context.Context {
	return context.Background()
}

type WrapRequest func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

func (f WrapRequest) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := NewContextFunc(r)
	ctx = context.WithValue(ctx, urlParamsKey, c.URLParams)
	err := f(ctx, w, r)
	if err != nil {
		if httpErr, ok := err.(httperr.Error); ok && httpErr.PrivateError != nil {
			log.Printf("ERROR: %s (%s)", httpErr.PrivateError, httpErr.StatusCode)
		} else {
			log.Printf("ERROR: %s", err)
		}
		httperr.Write(w, err)
	}
}
