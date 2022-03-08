package app

import "github.com/stdevHsequeda/SendToChannelsBot/pkgs/administration/app/command"

type App struct {
	Commands struct {
		AddInput command.AddInputHandler
	}
}

func NewApp(addInput command.AddInputHandler) App {
	return App{
		Commands: struct{ AddInput command.AddInputHandler }{
			AddInput: addInput,
		},
	}
}
