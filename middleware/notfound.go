package middleware

import (
	"net/http"
)

// NotFoundmiddleware simply returns a 404 response for any requests that
// happen to get this far down the middleware chain.
type NotFoundMiddleware struct{}

func NewNotFound() *NotFoundMiddleware {
	return &NotFoundMiddleware{}
}

func (n *NotFoundMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte("404 Not Found"))
}
