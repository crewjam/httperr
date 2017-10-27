package httperrctx

import (
	"log"
	"net/http"

	"golang.org/x/net/context"

	"github.com/crewjam/httperr"
)

type Handler interface {
	ServeHTTPErr(context.Context, http.ResponseWriter, *http.Request) error
}

type HandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

func (f HandlerFunc) ServeHTTPErr(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return f(ctx, w, r)
}

type WrappedHandler interface {
	http.Handler
	ServeHTTPC(context.Context, http.ResponseWriter, *http.Request)
}

func Wrap(h Handler) WrappedHandler {
	return &wrapHandler{Handler: h}
}

func WrapFunc(f func(ctx context.Context, w http.ResponseWriter, r *http.Request) error) WrappedHandler {
	return &wrapHandler{Handler: HandlerFunc(f)}
}

type wrapHandler struct {
	Handler Handler
}

func (wh *wrapHandler) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := wh.Handler.ServeHTTPErr(ctx, w, r); err != nil {
		if httpErr, ok := err.(httperr.Error); ok && httpErr.PrivateError != nil {
			log.Printf("ERROR: %s (%s)", httpErr.PrivateError, httpErr.StatusCode)
		} else {
			log.Printf("ERROR: %s", err)
		}
		httperr.Write(w, err)
	}
}

func (wh *wrapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wh.ServeHTTPC(context.TODO(), w, r)
}
