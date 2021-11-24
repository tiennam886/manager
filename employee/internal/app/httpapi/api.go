package httpapi

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/tiennam886/manager/pkg/logger"
)

func v1(r *chi.Mux) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Timeout(30 * time.Second))
		r.Use(middleware.Recoverer)

		r.Route("/employee", func(r chi.Router) {
			r.Post("/", EmployeeAdd)
			r.Get("/{uid}", EmployeeFindByUID)
			r.Get("/", EmployeeGetAll)

			r.Delete("/{uid}", EmployeeDeleteByUID)
			r.Patch("/{uid}", EmployeeUpdateByUID)

			r.Get("/event/{event}", EventHandler)

		})
	})
}

func Serve(ctx context.Context, addr string) (err error) {
	sugarLogger = logger.ConfigZap()
	defer func() {
		sugarLogger.Errorf("HTTP server stopped" + err.Error())
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
	sugarLogger.Errorf("HTTP server started at %s\n", addr)

	select {
	case <-ctx.Done():
		return nil
	case err = <-errChan:
		return err
	}
}
