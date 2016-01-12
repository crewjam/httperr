// package httperr implements an error object that speaks HTTP.
package httperr

import (
	"fmt"
	"net/http"
)

// Error represents an error that can be modeled as an
// http status code.
type Error struct {
	StatusCode   int    // the HTTP status code. If not supplied, http.StatusInternalServerError is used.
	Status       string // the HTTP status text. If not supplied, http.StatusText(http.StatusCode) is used.
	PrivateError error  // an additional error that is not displayed to the user, but may be logged
}

func (h Error) Error() string {
	if h.StatusCode == 0 {
		h.StatusCode = http.StatusInternalServerError
	}
	if h.Status == "" {
		h.Status = http.StatusText(h.StatusCode)
	}
	return fmt.Sprintf("%d %s", h.StatusCode, h.Status)
}

// WriteResponse writes an error response to w using the specified status code.
func (h Error) WriteResponse(w http.ResponseWriter) {
	if h.StatusCode == 0 {
		h.StatusCode = http.StatusInternalServerError
	}
	if h.Status == "" {
		h.Status = http.StatusText(h.StatusCode)
	}
	http.Error(w, h.Status, h.StatusCode)
}

// ResponseWriter is an interface for structs that know how to write themselves
// to a response. This interface is implemented by Error.
type ResponseWriter interface {
	WriteResponse(w http.ResponseWriter)
}

// Write writes the specified error to w. If err is a ResponseWriter the
// WriteResponse method is invoked to produce the response. Otherwise a
// generic 500 Internal Server Error is written.
func Write(w http.ResponseWriter, err error) {
	if wr, ok := err.(ResponseWriter); ok {
		wr.WriteResponse(w)
	} else {
		wr := Error{PrivateError: err}
		wr.WriteResponse(w)
	}
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := f(w, r); err != nil {
		Write(w, err)
	}
}

var (
	BadRequest                   = Error{StatusCode: 400}
	Unauthorized                 = Error{StatusCode: 401}
	PaymentRequired              = Error{StatusCode: 402}
	Forbidden                    = Error{StatusCode: 403}
	NotFound                     = Error{StatusCode: 404}
	MethodNotAllowed             = Error{StatusCode: 405}
	NotAcceptable                = Error{StatusCode: 406}
	ProxyAuthRequired            = Error{StatusCode: 407}
	RequestTimeout               = Error{StatusCode: 408}
	Conflict                     = Error{StatusCode: 409}
	Gone                         = Error{StatusCode: 410}
	LengthRequired               = Error{StatusCode: 411}
	PreconditionFailed           = Error{StatusCode: 412}
	RequestEntityTooLarge        = Error{StatusCode: 413}
	RequestURITooLong            = Error{StatusCode: 414}
	UnsupportedMediaType         = Error{StatusCode: 415}
	RequestedRangeNotSatisfiable = Error{StatusCode: 416}
	ExpectationFailed            = Error{StatusCode: 417}
	Teapot                       = Error{StatusCode: 418}
	InternalServerError          = Error{StatusCode: 500}
	NotImplemented               = Error{StatusCode: 501}
	BadGateway                   = Error{StatusCode: 502}
	ServiceUnavailable           = Error{StatusCode: 503}
	GatewayTimeout               = Error{StatusCode: 504}
	HTTPVersionNotSupported      = Error{StatusCode: 505}
)
