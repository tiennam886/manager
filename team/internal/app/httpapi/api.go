package httpapi

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net"
	"net/http"
	"time"
)

func v1(r *chi.Mux) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Timeout(30 * time.Second))
		r.Use(middleware.Recoverer)

		r.Route("/team", func(r chi.Router) {
			r.Post("/", TeamAdd)
			r.Get("/{uid}", TeamFindByUID)
			r.Delete("/{uid}", TeamDeleteByUID)

			r.Get("/notice", TeamNotice)
		})
	})
}

func Serve(ctx context.Context, addr string) (err error) {
	defer func() {
		log.Println("HTTP server stopped", err)
	}()

	r := chi.NewRouter()

	v1(r)

	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()

	srv := http.Server{
		Addr:    addr,
		Handler: r,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	errChan := make(chan error, 1)

	go func(ctx context.Context, errChan chan error) {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}(ctx, errChan)

	log.Printf("HTTP server started at %s\n", addr)

	select {
	case <-ctx.Done():
		return nil
	case err = <-errChan:
		return err
	}
}
