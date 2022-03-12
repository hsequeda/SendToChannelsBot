package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/administration/app"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/administration/app/command"
)

type HttpServer struct {
	app app.App
}

func NewHttpServer(app app.App) HttpServer {
	return HttpServer{app}
}

// (POST /addInput)
func (h HttpServer) AddInput(w http.ResponseWriter, r *http.Request) {
	var reqBody AddInputReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		render.Respond(w, r, ErrorResp{
			HttpStatus: http.StatusInternalServerError,
			Name:       "bad_request",
		})
		return
	}

	if err := h.app.Commands.AddInput.Handle(r.Context(), command.AddInput{
		UserID:      reqBody.UserId,
		ChatID:      reqBody.ChatId,
		InputType:   reqBody.InputType,
		Name:        reqBody.Name,
		Description: reqBody.Description,
	}); err != nil {
		var strRef = err.Error()
		render.Respond(w, r, ErrorResp{
			HttpStatus: http.StatusInternalServerError,
			Message:    &strRef,
			Name:       "internal_server_error",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	render.Respond(w, r, SuccessResp{
		Success: true,
	})
}
