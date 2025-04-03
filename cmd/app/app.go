package app

import (
	"net/http"

	timetable_service "github.com/sohosai/ultradonguri-server/internal/core/service"
	"github.com/sohosai/ultradonguri-server/internal/handler"
	"github.com/sohosai/ultradonguri-server/internal/infra"
	"go.uber.org/fx"
)

func Run() {
	fx.New(
		fx.Provide(
			handler.NewHTTPServer,
			timetable_service.New,

			infra.NewQueueRepository,
			infra.NewTimetableRepository,
			infra.NewMusicRepository,

			handler.Route,
		),
		fx.Invoke(func(srv *http.Server) {}),
	).Run()
}
