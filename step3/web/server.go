package web

import (
	"context"
	"github.com/miaomiao3/log"
	"net/http"
	"time"
)

type APIServer struct {
	server *http.Server
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"a":1}`))
}


func TestParam(w http.ResponseWriter, r *http.Request) {
	name := GetParam(r, "name")
	id := GetParam(r, "id")
	log.Debug("name:", name)
	log.Debug("id:", id)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"a":1}`))
}

func NewAPIServer() *APIServer {
	r := NewRouter()

	r.GET("/", HealthHandler)
	r.GET("/test/:name/:id", TestParam)

	// log middle ware
	r.Use(loggerMiddleware)

	httpServer := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	APIServer := &APIServer{
		server: httpServer,
	}

	return APIServer
}

func (s *APIServer) Start() {
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Error("server start failed")
		}
	}()
}

func (s *APIServer) Stop() {
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	s.server.Shutdown(ctx)
}
