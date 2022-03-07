// Code generated by SQLBoiler 4.6.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package psql_model

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Input is an object representing the database table.
type Input struct {
	Ref         string `boil:"ref" json:"ref" toml:"ref" yaml:"ref"`
	Name        string `boil:"name" json:"name" toml:"name" yaml:"name"`
	Owner       string `boil:"owner" json:"owner" toml:"owner" yaml:"owner"`
	Inputtype   string `boil:"inputtype" json:"inputtype" toml:"inputtype" yaml:"inputtype"`
	Description string `boil:"description" json:"description" toml:"description" yaml:"description"`

	R *inputR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L inputL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var InputColumns = struct {
	Ref         string
	Name        string
	Owner       string
	Inputtype   string
	Description string
}{
	Ref:         "ref",
	Name:        "name",
	Owner:       "owner",
	Inputtype:   "inputtype",
	Description: "description",
}

var InputTableColumns = struct {
	Ref         string
	Name        string
	Owner       string
	Inputtype   string
	Description string
}{
	Ref:         "input.ref",
	Name:        "input.name",
	Owner:       "input.owner",
	Inputtype:   "input.inputtype",
	Description: "input.description",
}

// Generated where

var InputWhere = struct {
	Ref         whereHelperstring
	Name        whereHelperstring
	Owner       whereHelperstring
	Inputtype   whereHelperstring
	Description whereHelperstring
}{
	Ref:         whereHelperstring{field: "\"input\".\"ref\""},
	Name:        whereHelperstring{field: "\"input\".\"name\""},
	Owner:       whereHelperstring{field: "\"input\".\"owner\""},
	Inputtype:   whereHelperstring{field: "\"input\".\"inputtype\""},
	Description: whereHelperstring{field: "\"input\".\"description\""},
}

// InputRels is where relationship names are stored.
var InputRels = struct {
}{}

// inputR is where relationships are stored.
type inputR struct {
}

// NewStruct creates a new relationship struct
func (*inputR) NewStruct() *inputR {
	return &inputR{}
}

// inputL is where Load methods for each relationship are stored.
type inputL struct{}

var (
	inputAllColumns            = []string{"ref", "name", "owner", "inputtype", "description"}
	inputColumnsWithoutDefault = []string{"ref", "name", "owner", "inputtype", "description"}
	inputColumnsWithDefault    = []string{}
	inputPrimaryKeyColumns     = []string{"ref"}
)

type (
	// InputSlice is an alias for a slice of pointers to Input.
	// This should almost always be used instead of []Input.
	InputSlice []*Input

	inputQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	inputType                 = reflect.TypeOf(&Input{})
	inputMapping              = queries.MakeStructMapping(inputType)
	inputPrimaryKeyMapping, _ = queries.BindMapping(inputType, inputMapping, inputPrimaryKeyColumns)
	inputInsertCacheMut       sync.RWMutex
	inputInsertCache          = make(map[string]insertCache)
	inputUpdateCacheMut       sync.RWMutex
	inputUpdateCache          = make(map[string]updateCache)
	inputUpsertCacheMut       sync.RWMutex
	inputUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// OneG returns a single input record from the query using the global executor.
func (q inputQuery) OneG(ctx context.Context) (*Input, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single input record from the query.
func (q inputQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Input, error) {
	o := &Input{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "psql_model: failed to execute a one query for input")
	}

	return o, nil
}

// AllG returns all Input records from the query using the global executor.
func (q inputQuery) AllG(ctx context.Context) (InputSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all Input records from the query.
func (q inputQuery) All(ctx context.Context, exec boil.ContextExecutor) (InputSlice, error) {
	var o []*Input

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "psql_model: failed to assign all query results to Input slice")
	}

	return o, nil
}

// CountG returns the count of all Input records in the query, and panics on error.
func (q inputQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all Input records in the query.
func (q inputQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "psql_model: failed to count input rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q inputQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q inputQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "psql_model: failed to check if input exists")
	}

	return count > 0, nil
}

// Inputs retrieves all the records using an executor.
func Inputs(mods ...qm.QueryMod) inputQuery {
	mods = append(mods, qm.From("\"input\""))
	return inputQuery{NewQuery(mods...)}
}

// FindInputG retrieves a single record by ID.
func FindInputG(ctx context.Context, ref string, selectCols ...string) (*Input, error) {
	return FindInput(ctx, boil.GetContextDB(), ref, selectCols...)
}

// FindInput retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindInput(ctx context.Context, exec boil.ContextExecutor, ref string, selectCols ...string) (*Input, error) {
	inputObj := &Input{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"input\" where \"ref\"=$1", sel,
	)

	q := queries.Raw(query, ref)

	err := q.Bind(ctx, exec, inputObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "psql_model: unable to select from input")
	}

	return inputObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Input) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Input) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("psql_model: no input provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(inputColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	inputInsertCacheMut.RLock()
	cache, cached := inputInsertCache[key]
	inputInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			inputAllColumns,
			inputColumnsWithDefault,
			inputColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(inputType, inputMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(inputType, inputMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"input\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"input\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "psql_model: unable to insert into input")
	}

	if !cached {
		inputInsertCacheMut.Lock()
		inputInsertCache[key] = cache
		inputInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single Input record using the global executor.
// See Update for more documentation.
func (o *Input) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the Input.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Input) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	inputUpdateCacheMut.RLock()
	cache, cached := inputUpdateCache[key]
	inputUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			inputAllColumns,
			inputPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("psql_model: unable to update input, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"input\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, inputPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(inputType, inputMapping, append(wl, inputPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "psql_model: unable to update input row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "psql_model: failed to get rows affected by update for input")
	}

	if !cached {
		inputUpdateCacheMut.Lock()
		inputUpdateCache[key] = cache
		inputUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (q inputQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q inputQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "psql_model: unable to update all for input")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "psql_model: unable to retrieve rows affected for input")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o InputSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o InputSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("psql_model: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), inputPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"input\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, inputPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "psql_model: unable to update all in input slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "psql_model: unable to retrieve rows affected all in update all input")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Input) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Input) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("psql_model: no input provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(inputColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	inputUpsertCacheMut.RLock()
	cache, cached := inputUpsertCache[key]
	inputUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			inputAllColumns,
			inputColumnsWithDefault,
			inputColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			inputAllColumns,
			inputPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("psql_model: unable to upsert input, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(inputPrimaryKeyColumns))
			copy(conflict, inputPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"input\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(inputType, inputMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(inputType, inputMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "psql_model: unable to upsert input")
	}

	if !cached {
		inputUpsertCacheMut.Lock()
		inputUpsertCache[key] = cache
		inputUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteG deletes a single Input record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Input) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single Input record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Input) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("psql_model: no Input provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), inputPrimaryKeyMapping)
	sql := "DELETE FROM \"input\" WHERE \"ref\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "psql_model: unable to delete from input")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "psql_model: failed to get rows affected by delete for input")
	}

	return rowsAff, nil
}

func (q inputQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q inputQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("psql_model: no inputQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "psql_model: unable to delete all from input")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "psql_model: failed to get rows affected by deleteall for input")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o InputSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o InputSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), inputPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"input\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, inputPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "psql_model: unable to delete all from input slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "psql_model: failed to get rows affected by deleteall for input")
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Input) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("psql_model: no Input provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Input) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindInput(ctx, exec, o.Ref)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *InputSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("psql_model: empty InputSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *InputSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := InputSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), inputPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"input\".* FROM \"input\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, inputPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "psql_model: unable to reload all in InputSlice")
	}

	*o = slice

	return nil
}

// InputExistsG checks if the Input row exists.
func InputExistsG(ctx context.Context, ref string) (bool, error) {
	return InputExists(ctx, boil.GetContextDB(), ref)
}

// InputExists checks if the Input row exists.
func InputExists(ctx context.Context, exec boil.ContextExecutor, ref string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"input\" where \"ref\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, ref)
	}
	row := exec.QueryRowContext(ctx, sql, ref)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "psql_model: unable to check if input exists")
	}

	return exists, nil
}