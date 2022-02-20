package app

import "github.com/stdevHsequeda/SendToChannelsBot/pkgs/account/app/command"

type App struct {
	Commands struct {
		RegisterAccount command.RegisterAccountHandler
	}
}

// NewApp TODO
func NewApp(registerAccount command.RegisterAccountHandler) App {
	return App{
		Commands: struct {
			RegisterAccount command.RegisterAccountHandler
		}{
			RegisterAccount: registerAccount,
		},
	}
}
