package rest

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/account/app"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/account/app/command"
)

type HttpServer struct {
	app app.App
}

func NewHttpServer(app app.App) HttpServer {
	return HttpServer{app}
}

// (POST /register)
func (h HttpServer) Register(w http.ResponseWriter, r *http.Request) {
	var registerReqBody RegisterReqBody
	if err := render.Decode(r, &registerReqBody); err != nil {
		w.WriteHeader(400)
		render.Respond(w, r, ErrorResp{
			HttpStatus: 400,
			Name:       "bad_request",
		})
		return
	}

	if err := h.app.Commands.RegisterAccount.Handle(r.Context(), command.RegisterAccount{
		ID:         uuid.New().String(),
		TelegramID: registerReqBody.TelegramId,
	}); err != nil {
		w.WriteHeader(500)
		render.Respond(w, r, ErrorResp{
			HttpStatus: 500,
			Name:       "internal_server_error",
		})
		return
	}

	w.WriteHeader(200)
	render.Respond(w, r, SuccessResp{
		Success: true,
	})
}
