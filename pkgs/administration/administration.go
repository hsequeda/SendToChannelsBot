package administration

import (
	"database/sql"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/administration/adapter"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/administration/app"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/administration/app/command"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/administration/domain"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/administration/port/rest"
)

type ModuleConfiguration struct {
	PsqlConn             *sql.DB
	TgBotAPI             *tgbotapi.BotAPI
	BasicAuthCredentials struct {
		User string
		Pass string
	}
}

type Module struct {
	httpServer           rest.HttpServer
	basicAuthCredentials struct {
		user string
		pass string
	}
}

var (
	administration *Module
	once           sync.Once
)

func NewModule(conf ModuleConfiguration) *Module {
	once.Do(func() {
		telegramService := adapter.NewTelegramServiceTgBotApi(conf.TgBotAPI)
		administration = &Module{
			httpServer: rest.NewHttpServer(app.NewApp(command.NewAddInputHandler(domain.NewInputFactory(telegramService), adapter.NewInputPsqlRepository(conf.PsqlConn)))),
			basicAuthCredentials: struct {
				user string
				pass string
			}{
				user: conf.BasicAuthCredentials.User,
				pass: conf.BasicAuthCredentials.Pass,
			},
		}
	})

	return administration
}

func (m Module) BindRouter(router chi.Router) http.Handler {
	router.Use(
		middleware.BasicAuth("Production Realm", map[string]string{
			m.basicAuthCredentials.user: m.basicAuthCredentials.pass,
		}),
	)

	return rest.HandlerFromMux(m.httpServer, router)
}
