package httpapi

import (
	"context"
	"fmt"
	"github.com/tiennam886/manager/pkg/logger"
	"github.com/tiennam886/manager/team/internal/config"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func v1(r *chi.Mux) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Timeout(30 * time.Second))
		r.Use(middleware.Recoverer)
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "*")

				next.ServeHTTP(w, r)
			})

		})
		r.Route("/team", func(r chi.Router) {
			r.Post("/", TeamAdd)
			r.Get("/{uid}", TeamFindByUID)
			r.Get("/", TeamGetAll)

			r.Patch("/{uid}", TeamUpdateByUID)
			r.Delete("/{uid}", TeamDeleteByUID)

			r.Options("/", TeamOption)
			r.Options("/{param}", TeamOption)

			r.Get("/notice", TeamNotice)
		})
	})
}

func Serve(ctx context.Context, addr string) (err error) {
	sugarLogger = logger.ConfigZap(config.Get().LogLevel)
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
	sugarLogger.Infow(fmt.Sprintf("HTTP server started at %s", addr))

	select {
	case <-ctx.Done():
		return nil
	case err = <-errChan:
		return err
	}
}
