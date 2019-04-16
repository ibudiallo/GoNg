package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ibudiallo/gong"
)

type Server struct {
	Env *gong.Env
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	s.setDefaultHeaders(w)
	if err := s.serveFile(w, r); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "page not found")
	}
	log.Println("Request", r.URL, t.Format(time.RFC3339))
}

func (s *Server) serveFile(w http.ResponseWriter, r *http.Request) error {
	uri := strings.Replace(r.URL.Path, r.URL.RawQuery, "", -1)
	path := filepath.Join(s.Env.Root, uri)
	http.ServeFile(w, r, path)
	return nil
}

func (s *Server) setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Powered-by", "GONG")
}

func Init(env *gong.Env) {
	srv := &http.Server{
		Addr:    env.Port,
		Handler: &Server{Env: env},
	}
	go func() {
		log.Printf("Server is ready to listen and serve on port %s.\n", env.Port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Failed: %v\n", err)
			os.Exit(1)
		}
	}()
	graceful(srv, 10*time.Second)
}

func graceful(srv *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGKILL, syscall.SIGHUP)
	<-stop
	log.Println("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Failed: %v\n", err)
		os.Exit(1)
	}
	log.Println("shut down complete")
}
