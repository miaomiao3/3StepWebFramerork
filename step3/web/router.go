package web

import (
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)

type QRouter struct {
	middleware []MiddlewareFunc
	router     *httprouter.Router
}

func NewRouter() *QRouter {
	r := &QRouter{
		router: httprouter.New(),
	}
	return r
}

// add middleware
func (r *QRouter) Use(middlewareFunc ...MiddlewareFunc) {
	r.middleware = append(r.middleware, middlewareFunc...)
}

func (r *QRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler := http.Handler(r.router)
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}
	handler.ServeHTTP(w, req)
}
func (r *QRouter) GET(path string, handler http.HandlerFunc) {
	r.router.GET(path, wrap(handler))
}

func (r *QRouter) HEAD(path string, handler http.HandlerFunc) {
	r.router.HEAD(path, wrap(handler))
}

func (r *QRouter) POST(path string, handler http.HandlerFunc) {
	r.router.POST(path, wrap(handler))
}

func (r *QRouter) OPTIONS(path string, handler http.HandlerFunc) {
	r.router.OPTIONS(path, wrap(handler))
}

func (r *QRouter) PUT(path string, handler http.HandlerFunc) {
	r.router.PUT(path, wrap(handler))
}

func (r *QRouter) DELETE(path string, handler http.HandlerFunc) {
	r.router.DELETE(path, wrap(handler))
}

func wrap(handler http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		SetParam(r,  params)
		handler(w, r)
	}
}

// MiddlewareFunc is a function which receives an http.Handler and returns another http.Handler.
type MiddlewareFunc func(http.Handler) http.Handler


type contextReadCloser struct {
	io.ReadCloser
	params httprouter.Params
}

func SetParam(req *http.Request, params httprouter.Params) {
	crc := getContextReadCloser(req)
	crc.params = params
}

func GetParam(req *http.Request, key string) string {
	crc := getContextReadCloser(req)
	for _, v := range crc.params {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}

// add a simple context
func getContextReadCloser(req *http.Request) *contextReadCloser {
	crc, ok := req.Body.(*contextReadCloser)
	if !ok {
		crc = &contextReadCloser{
			ReadCloser: req.Body,
		}
		req.Body = crc
	}
	return crc
}
