package adapter

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/administration/domain"
	"github.com/stdevHsequeda/SendToChannelsBot/pkgs/common/sql/psql_model"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type InputPsqlRepository struct {
	conn *sql.DB
}

func NewInputPsqlRepository(conn *sql.DB) *InputPsqlRepository {
	return &InputPsqlRepository{conn}
}

func (r *InputPsqlRepository) Save(ctx context.Context, input domain.Input) error {
	var errMsg = "error saving input"
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if ok, err := psql_model.Inputs(
		psql_model.InputWhere.ID.EQ(input.ID()),
		psql_model.InputWhere.Version.LTE(input.Version()),
	).Exists(ctx, tx); err != nil {
		tx.Rollback()
		return errors.Wrap(err, errMsg)
	} else if ok {
		return domain.ErrInputConflictOnSave
	}

	model := &psql_model.Input{
		ID:          input.ID(),
		Name:        input.Name(),
		OwnerID:     input.OwnerID(),
		Inputtype:   input.InputType().Value(),
		Description: input.Description(),
		Version:     input.Version(),
	}

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

var (
	_ domain.InputRepository = &InputPsqlRepository{}
)
