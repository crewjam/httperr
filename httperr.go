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

var (
	ErrContinue                     = Error{StatusCode: 100}
	ErrSwitchingProtocols           = Error{StatusCode: 101}
	ErrOK                           = Error{StatusCode: 200}
	ErrCreated                      = Error{StatusCode: 201}
	ErrAccepted                     = Error{StatusCode: 202}
	ErrNonAuthoritativeInfo         = Error{StatusCode: 203}
	ErrNoContent                    = Error{StatusCode: 204}
	ErrResetContent                 = Error{StatusCode: 205}
	ErrPartialContent               = Error{StatusCode: 206}
	ErrMultipleChoices              = Error{StatusCode: 300}
	ErrMovedPermanently             = Error{StatusCode: 301}
	ErrFound                        = Error{StatusCode: 302}
	ErrSeeOther                     = Error{StatusCode: 303}
	ErrNotModified                  = Error{StatusCode: 304}
	ErrUseProxy                     = Error{StatusCode: 305}
	ErrTemporaryRedirect            = Error{StatusCode: 307}
	ErrBadRequest                   = Error{StatusCode: 400}
	ErrUnauthorized                 = Error{StatusCode: 401}
	ErrPaymentRequired              = Error{StatusCode: 402}
	ErrForbidden                    = Error{StatusCode: 403}
	ErrNotFound                     = Error{StatusCode: 404}
	ErrMethodNotAllowed             = Error{StatusCode: 405}
	ErrNotAcceptable                = Error{StatusCode: 406}
	ErrProxyAuthRequired            = Error{StatusCode: 407}
	ErrRequestTimeout               = Error{StatusCode: 408}
	ErrConflict                     = Error{StatusCode: 409}
	ErrGone                         = Error{StatusCode: 410}
	ErrLengthRequired               = Error{StatusCode: 411}
	ErrPreconditionFailed           = Error{StatusCode: 412}
	ErrRequestEntityTooLarge        = Error{StatusCode: 413}
	ErrRequestURITooLong            = Error{StatusCode: 414}
	ErrUnsupportedMediaType         = Error{StatusCode: 415}
	ErrRequestedRangeNotSatisfiable = Error{StatusCode: 416}
	ErrExpectationFailed            = Error{StatusCode: 417}
	ErrTeapot                       = Error{StatusCode: 418}
	ErrInternalServerError          = Error{StatusCode: 500}
	ErrNotImplemented               = Error{StatusCode: 501}
	ErrBadGateway                   = Error{StatusCode: 502}
	ErrServiceUnavailable           = Error{StatusCode: 503}
	ErrGatewayTimeout               = Error{StatusCode: 504}
	ErrHTTPVersionNotSupported      = Error{StatusCode: 505}
)
