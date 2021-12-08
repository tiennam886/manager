package httpapi

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/tiennam886/manager/employee/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/tiennam886/manager/pkg/logger"
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
		r.Route("/employee", func(r chi.Router) {
			r.Post("/", EmployeeAdd)
			r.Post("/{eid}/team/{tid}", EmployeeAddToTeam)

			r.Options("/", EmployeeOption)
			r.Options("/{param}", EmployeeOption)
			r.Options("/{eid}/team/{tid}", EmployeeOption)

			r.Get("/{uid}", EmployeeFindByUID)
			r.Get("/list/{uid}", EmployeeFindTeams)
			r.Get("/", EmployeeGetAll)

			r.Delete("/{uid}", EmployeeDeleteByUID)
			r.Delete("/{eid}/team/{tid}", EmployeeRemoveFromTeam)

			r.Patch("/{uid}", EmployeeUpdateByUID)

			r.Get("/event/{event}", EventHandler)

		})
	})
}

func Serve(ctx context.Context, addr string) (err error) {
	sugarLogger = logger.ConfigZap(config.Get().LogLevel)
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
	sugarLogger.Infow("HTTP server started at " + addr)

	select {
	case <-ctx.Done():
		return nil
	case err = <-errChan:
		return err
	}
}
