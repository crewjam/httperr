package httperr

import (
	"fmt"
	"net/http"
)

// Response is an alias for http.Response that implements
// the error interface. Example:
//
//   resp, err := http.Get("http://www.example.com")
//   if err != nil {
//   	return err
//   }
//   if resp.StatusCode != http.StatusOK {
//   	return httperr.Response(*resp)
//   }
//   // ...
//
type Response http.Response

func (re Response) Error() string {
	return fmt.Sprintf("%d %s", re.StatusCode, re.Status)
}
