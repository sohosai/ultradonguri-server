package handler

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	timetable_service "github.com/sohosai/ultradonguri-server/internal/core/service"
	"go.uber.org/fx"
)

func NewHTTPServer(lc fx.Lifecycle, mux *chi.Mux) *http.Server {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			l, err := net.Listen("tcp", ":8080")
			if err != nil {
				return err
			}
			go func() {
				if err := http.Serve(l, mux); err != nil {
					log.Println(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
	return &http.Server{
		Addr: ":8080",
	}
}

func Route(qs *timetable_service.TimetableService) *chi.Mux {

	api := chi.NewRouter()
	api.Get("/queue", func(w http.ResponseWriter, r *http.Request) {
		queue, err := qs.TimetableByID("1")
		log.Println(queue)
		if err != nil {
			http.Error(w, "Queue not found", http.StatusInternalServerError)
			return
		}
		w.Write(fmt.Appendf(nil, "Queue: %v", queue))
	})
	return api
}
