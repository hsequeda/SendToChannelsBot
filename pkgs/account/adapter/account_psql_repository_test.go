package adapter_test

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/account/adapter"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/account/domain"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/common/sql"
	"github.com/stretchr/testify/require"
)

func TestAccountRepositoryPsql_Save(t *testing.T) {
	t.Parallel()
	accountRepositoryPsql := newAccountRepositoryPsql(t)
	t.Run("Ok case: Save an Account once", func(t *testing.T) {
		telegramId, _ := domain.NewTelegramID(int64(uuid.New().ID()))
		accountId := uuid.New().String()
		account, err := domain.NewAccount(accountId, telegramId)
		require.NoError(t, err)
		err = accountRepositoryPsql.Save(context.Background(), account)
		require.NoError(t, err)

		retrievedAccount, err := accountRepositoryPsql.GetOneByTgID(context.Background(), telegramId)
		require.NoError(t, err)
		reflect.DeepEqual(retrievedAccount, account)
	})

	t.Run("Ok case: Save an Account twice", func(t *testing.T) {
		telegramId, _ := domain.NewTelegramID(int64(uuid.New().ID()))
		accountId := uuid.New().String()
		account, err := domain.NewAccount(accountId, telegramId)
		require.NoError(t, err)
		err = accountRepositoryPsql.Save(context.Background(), account)
		require.NoError(t, err)

		retrievedAccount, err := accountRepositoryPsql.GetOneByTgID(context.Background(), telegramId)
		require.NoError(t, err)
		reflect.DeepEqual(retrievedAccount, account)
	})
}

func newAccountRepositoryPsql(t *testing.T) *adapter.AccountPsqlRepository {
	t.Helper()
	conn, err := sql.NewPostgresConnPoolFromConf(sql.PsqlDatabaseConfiguration{
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Host:     os.Getenv("DB_HOST"),
		SSLmode:  os.Getenv("DB_SSLMODE"),
	})

	require.NoError(t, err)
	return adapter.NewAccountRepositoryPsql(conn)
}
