package adapter

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/account/domain"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/common/sql/psql_model"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type AccountPsqlRepository struct {
	conn *sql.DB
}

// NewAccountRepository TODO
func NewAccountRepositoryPsql(conn *sql.DB) *AccountPsqlRepository {
	return &AccountPsqlRepository{conn: conn}
}

// Save create or update an account in the repository. Verify optimistic locking.
func (r *AccountPsqlRepository) Save(ctx context.Context, account *domain.Account) error {
	var errMsg = "error saving account"
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	// if exist an account with the same ID and a lesser or equal version then returns OK = true otherwise return OK = false
	if ok, err := psql_model.Accounts(psql_model.AccountWhere.TelegramID.EQ(account.TelegramID().Value()),
		psql_model.AccountWhere.Version.LTE(int64(account.Version())),
	).Exists(ctx, tx); err != nil {
		tx.Rollback()
		return errors.Wrap(err, errMsg)
	} else if ok {
		tx.Rollback()
		return domain.ErrAccountConflictOnSave
	}

	model := r.unmarshalAccount(account)
	if err := model.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		tx.Rollback()
		return errors.Wrap(err, errMsg)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.Wrap(err, errMsg)
	}

	return nil
}

// GetOneByTgID returns an account by a TelegramID.
func (r *AccountPsqlRepository) GetOneByTgID(ctx context.Context, tgID domain.TelegramID) (*domain.Account, error) {

	model, err := psql_model.Accounts(psql_model.AccountWhere.TelegramID.EQ(tgID.Value())).One(ctx, r.conn)
	if err != nil {
		return nil, errors.Wrap(err, "error getting one account by its Telegram ID")
	}

	return r.marshalAccount(model), nil
}

// GetOneByID returns an account by an ID.
func (r *AccountPsqlRepository) GetOneByID(ctx context.Context, id string) (*domain.Account, error) {
	model, err := psql_model.Accounts(psql_model.AccountWhere.ID.EQ(id)).One(ctx, r.conn)
	if err != nil {
		return nil, errors.Wrap(err, "error getting one account by its ID")
	}

	return r.marshalAccount(model), nil
}

func (r *AccountPsqlRepository) marshalAccount(model *psql_model.Account) *domain.Account {
	telegramID, err := domain.NewTelegramID(model.TelegramID)
	if err != nil {
		panic(errors.Wrap(err, "error creating a telegramID from the repository"))
	}

	return domain.CreateAccountFromRepository(model.ID, telegramID, uint(model.Version))
}

func (r *AccountPsqlRepository) unmarshalAccount(account *domain.Account) *psql_model.Account {
	return &psql_model.Account{
		ID:         account.ID(),
		TelegramID: account.TelegramID().Value(),
		Version:    int64(account.Version()),
	}
}
