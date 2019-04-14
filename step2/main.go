package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// get url param by GetParam
	work := GetParam(r, "work")
	fmt.Println("work:", work)

	fmt.Fprint(w, "Welcome!\n")
}

type MiddlewareFunc func(w http.ResponseWriter, req *http.Request)

type simpleRouter struct {
	BeforeMiddleware MiddlewareFunc
	AfterMiddleware  MiddlewareFunc
	r                *httprouter.Router
}

func NewSimpleRouter() (s *simpleRouter) {
	s = &simpleRouter{
		r: httprouter.New(),
	}
	return s
}

func (s *simpleRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if s.BeforeMiddleware != nil {
		s.BeforeMiddleware(w, req)
	}

	s.r.ServeHTTP(w, req)

	if s.AfterMiddleware != nil {
		s.AfterMiddleware(w, req)
	}
}

// wrapper for httprouter GET
func (s *simpleRouter) GET(path string, handle http.HandlerFunc) {
	s.r.GET(path, wrap(handle))
}

func wrap(handler http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		SetParam(r, params)
		handler(w, r)
	}
}

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

func main() {
	router := NewSimpleRouter()
	router.BeforeMiddleware = func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("handle before")
		w.Header().Set("Before", "hello")
	}
	router.AfterMiddleware = func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("handle after")
		w.Header().Set("After", "world") // 注意，因为实际处理函数在处理的时候一般会先写resp code, 这里是不会生效的
	}
	router.GET("/step2/:work", Index)

	log.Fatal(http.ListenAndServe(":8082", router))
}
