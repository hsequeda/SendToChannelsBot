package account

import (
	"database/sql"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/account/adapter"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/account/app"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/account/app/command"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/account/port/rest"
)

type AccountModule struct {
	httpServer           rest.HttpServer
	basicAuthCredentials struct {
		user string
		pass string
	}
}

type AccountModuleConfig struct {
	PsqlConn             *sql.DB
	BasicAuthCredentials struct {
		User string
		Pass string
	}
}

var (
	account *AccountModule
	once    sync.Once
)

func NewModule(conf AccountModuleConfig) *AccountModule {
	once.Do(func() {
		accountRepo := adapter.NewAccountRepositoryPsql(conf.PsqlConn)
		app := app.NewApp(command.NewRegisterAccountHandler(accountRepo))
		httpServer := rest.NewHttpServer(app)
		account = &AccountModule{
			httpServer: httpServer,
			basicAuthCredentials: struct {
				user string
				pass string
			}{
				user: conf.BasicAuthCredentials.User,
				pass: conf.BasicAuthCredentials.Pass,
			}}
	})

	return account
}

func (a AccountModule) BindRouter(router chi.Router) http.Handler {
	router.Use(middleware.BasicAuth("Production Realm", map[string]string{a.basicAuthCredentials.user: a.basicAuthCredentials.pass}))
	return rest.HandlerFromMux(a.httpServer, router)
}
