// Code generated by SQLBoiler 4.5.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

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

// Ref is an object representing the database table.
type Ref struct {
	ID           int64  `boil:"id" json:"id" toml:"id" yaml:"id"`
	RepositoryID int64  `boil:"repository_id" json:"repository_id" toml:"repository_id" yaml:"repository_id"`
	Ref          string `boil:"ref" json:"ref" toml:"ref" yaml:"ref"`
	Sha          string `boil:"sha" json:"sha" toml:"sha" yaml:"sha"`

	R *refR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L refL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var RefColumns = struct {
	ID           string
	RepositoryID string
	Ref          string
	Sha          string
}{
	ID:           "id",
	RepositoryID: "repository_id",
	Ref:          "ref",
	Sha:          "sha",
}

// Generated where

var RefWhere = struct {
	ID           whereHelperint64
	RepositoryID whereHelperint64
	Ref          whereHelperstring
	Sha          whereHelperstring
}{
	ID:           whereHelperint64{field: "\"refs\".\"id\""},
	RepositoryID: whereHelperint64{field: "\"refs\".\"repository_id\""},
	Ref:          whereHelperstring{field: "\"refs\".\"ref\""},
	Sha:          whereHelperstring{field: "\"refs\".\"sha\""},
}

// RefRels is where relationship names are stored.
var RefRels = struct {
	Repository         string
	BaseRefSubmissions string
	HeadRefSubmissions string
}{
	Repository:         "Repository",
	BaseRefSubmissions: "BaseRefSubmissions",
	HeadRefSubmissions: "HeadRefSubmissions",
}

// refR is where relationships are stored.
type refR struct {
	Repository         *Repository     `boil:"Repository" json:"Repository" toml:"Repository" yaml:"Repository"`
	BaseRefSubmissions SubmissionSlice `boil:"BaseRefSubmissions" json:"BaseRefSubmissions" toml:"BaseRefSubmissions" yaml:"BaseRefSubmissions"`
	HeadRefSubmissions SubmissionSlice `boil:"HeadRefSubmissions" json:"HeadRefSubmissions" toml:"HeadRefSubmissions" yaml:"HeadRefSubmissions"`
}

// NewStruct creates a new relationship struct
func (*refR) NewStruct() *refR {
	return &refR{}
}

// refL is where Load methods for each relationship are stored.
type refL struct{}

var (
	refAllColumns            = []string{"id", "repository_id", "ref", "sha"}
	refColumnsWithoutDefault = []string{"repository_id", "ref", "sha"}
	refColumnsWithDefault    = []string{"id"}
	refPrimaryKeyColumns     = []string{"id"}
)

type (
	// RefSlice is an alias for a slice of pointers to Ref.
	// This should generally be used opposed to []Ref.
	RefSlice []*Ref
	// RefHook is the signature for custom Ref hook methods
	RefHook func(context.Context, boil.ContextExecutor, *Ref) error

	refQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	refType                 = reflect.TypeOf(&Ref{})
	refMapping              = queries.MakeStructMapping(refType)
	refPrimaryKeyMapping, _ = queries.BindMapping(refType, refMapping, refPrimaryKeyColumns)
	refInsertCacheMut       sync.RWMutex
	refInsertCache          = make(map[string]insertCache)
	refUpdateCacheMut       sync.RWMutex
	refUpdateCache          = make(map[string]updateCache)
	refUpsertCacheMut       sync.RWMutex
	refUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var refBeforeInsertHooks []RefHook
var refBeforeUpdateHooks []RefHook
var refBeforeDeleteHooks []RefHook
var refBeforeUpsertHooks []RefHook

var refAfterInsertHooks []RefHook
var refAfterSelectHooks []RefHook
var refAfterUpdateHooks []RefHook
var refAfterDeleteHooks []RefHook
var refAfterUpsertHooks []RefHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Ref) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range refBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Ref) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range refBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Ref) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range refBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Ref) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range refBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Ref) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range refAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Ref) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range refAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Ref) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range refAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Ref) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range refAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Ref) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range refAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddRefHook registers your hook function for all future operations.
func AddRefHook(hookPoint boil.HookPoint, refHook RefHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		refBeforeInsertHooks = append(refBeforeInsertHooks, refHook)
	case boil.BeforeUpdateHook:
		refBeforeUpdateHooks = append(refBeforeUpdateHooks, refHook)
	case boil.BeforeDeleteHook:
		refBeforeDeleteHooks = append(refBeforeDeleteHooks, refHook)
	case boil.BeforeUpsertHook:
		refBeforeUpsertHooks = append(refBeforeUpsertHooks, refHook)
	case boil.AfterInsertHook:
		refAfterInsertHooks = append(refAfterInsertHooks, refHook)
	case boil.AfterSelectHook:
		refAfterSelectHooks = append(refAfterSelectHooks, refHook)
	case boil.AfterUpdateHook:
		refAfterUpdateHooks = append(refAfterUpdateHooks, refHook)
	case boil.AfterDeleteHook:
		refAfterDeleteHooks = append(refAfterDeleteHooks, refHook)
	case boil.AfterUpsertHook:
		refAfterUpsertHooks = append(refAfterUpsertHooks, refHook)
	}
}

// One returns a single ref record from the query.
func (q refQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Ref, error) {
	o := &Ref{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for refs")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Ref records from the query.
func (q refQuery) All(ctx context.Context, exec boil.ContextExecutor) (RefSlice, error) {
	var o []*Ref

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Ref slice")
	}

	if len(refAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Ref records in the query.
func (q refQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count refs rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q refQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if refs exists")
	}

	return count > 0, nil
}

// Repository pointed to by the foreign key.
func (o *Ref) Repository(mods ...qm.QueryMod) repositoryQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.RepositoryID),
	}

	queryMods = append(queryMods, mods...)

	query := Repositories(queryMods...)
	queries.SetFrom(query.Query, "\"repositories\"")

	return query
}

// BaseRefSubmissions retrieves all the submission's Submissions with an executor via base_ref_id column.
func (o *Ref) BaseRefSubmissions(mods ...qm.QueryMod) submissionQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"submissions\".\"base_ref_id\"=?", o.ID),
	)

	query := Submissions(queryMods...)
	queries.SetFrom(query.Query, "\"submissions\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"submissions\".*"})
	}

	return query
}

// HeadRefSubmissions retrieves all the submission's Submissions with an executor via head_ref_id column.
func (o *Ref) HeadRefSubmissions(mods ...qm.QueryMod) submissionQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"submissions\".\"head_ref_id\"=?", o.ID),
	)

	query := Submissions(queryMods...)
	queries.SetFrom(query.Query, "\"submissions\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"submissions\".*"})
	}

	return query
}

// LoadRepository allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (refL) LoadRepository(ctx context.Context, e boil.ContextExecutor, singular bool, maybeRef interface{}, mods queries.Applicator) error {
	var slice []*Ref
	var object *Ref

	if singular {
		object = maybeRef.(*Ref)
	} else {
		slice = *maybeRef.(*[]*Ref)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &refR{}
		}
		args = append(args, object.RepositoryID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &refR{}
			}

			for _, a := range args {
				if a == obj.RepositoryID {
					continue Outer
				}
			}

			args = append(args, obj.RepositoryID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`repositories`),
		qm.WhereIn(`repositories.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Repository")
	}

	var resultSlice []*Repository
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Repository")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for repositories")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for repositories")
	}

	if len(refAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Repository = foreign
		if foreign.R == nil {
			foreign.R = &repositoryR{}
		}
		foreign.R.Refs = append(foreign.R.Refs, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.RepositoryID == foreign.ID {
				local.R.Repository = foreign
				if foreign.R == nil {
					foreign.R = &repositoryR{}
				}
				foreign.R.Refs = append(foreign.R.Refs, local)
				break
			}
		}
	}

	return nil
}

// LoadBaseRefSubmissions allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (refL) LoadBaseRefSubmissions(ctx context.Context, e boil.ContextExecutor, singular bool, maybeRef interface{}, mods queries.Applicator) error {
	var slice []*Ref
	var object *Ref

	if singular {
		object = maybeRef.(*Ref)
	} else {
		slice = *maybeRef.(*[]*Ref)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &refR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &refR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`submissions`),
		qm.WhereIn(`submissions.base_ref_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load submissions")
	}

	var resultSlice []*Submission
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice submissions")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on submissions")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for submissions")
	}

	if len(submissionAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.BaseRefSubmissions = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &submissionR{}
			}
			foreign.R.BaseRef = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.BaseRefID {
				local.R.BaseRefSubmissions = append(local.R.BaseRefSubmissions, foreign)
				if foreign.R == nil {
					foreign.R = &submissionR{}
				}
				foreign.R.BaseRef = local
				break
			}
		}
	}

	return nil
}

// LoadHeadRefSubmissions allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (refL) LoadHeadRefSubmissions(ctx context.Context, e boil.ContextExecutor, singular bool, maybeRef interface{}, mods queries.Applicator) error {
	var slice []*Ref
	var object *Ref

	if singular {
		object = maybeRef.(*Ref)
	} else {
		slice = *maybeRef.(*[]*Ref)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &refR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &refR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.ID) {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`submissions`),
		qm.WhereIn(`submissions.head_ref_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load submissions")
	}

	var resultSlice []*Submission
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice submissions")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on submissions")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for submissions")
	}

	if len(submissionAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.HeadRefSubmissions = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &submissionR{}
			}
			foreign.R.HeadRef = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.HeadRefID) {
				local.R.HeadRefSubmissions = append(local.R.HeadRefSubmissions, foreign)
				if foreign.R == nil {
					foreign.R = &submissionR{}
				}
				foreign.R.HeadRef = local
				break
			}
		}
	}

	return nil
}

// SetRepository of the ref to the related item.
// Sets o.R.Repository to related.
// Adds o to related.R.Refs.
func (o *Ref) SetRepository(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Repository) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"refs\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"repository_id"}),
		strmangle.WhereClause("\"", "\"", 2, refPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.RepositoryID = related.ID
	if o.R == nil {
		o.R = &refR{
			Repository: related,
		}
	} else {
		o.R.Repository = related
	}

	if related.R == nil {
		related.R = &repositoryR{
			Refs: RefSlice{o},
		}
	} else {
		related.R.Refs = append(related.R.Refs, o)
	}

	return nil
}

// AddBaseRefSubmissions adds the given related objects to the existing relationships
// of the ref, optionally inserting them as new records.
// Appends related to o.R.BaseRefSubmissions.
// Sets related.R.BaseRef appropriately.
func (o *Ref) AddBaseRefSubmissions(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Submission) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.BaseRefID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"submissions\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"base_ref_id"}),
				strmangle.WhereClause("\"", "\"", 2, submissionPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.BaseRefID = o.ID
		}
	}

	if o.R == nil {
		o.R = &refR{
			BaseRefSubmissions: related,
		}
	} else {
		o.R.BaseRefSubmissions = append(o.R.BaseRefSubmissions, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &submissionR{
				BaseRef: o,
			}
		} else {
			rel.R.BaseRef = o
		}
	}
	return nil
}

// AddHeadRefSubmissions adds the given related objects to the existing relationships
// of the ref, optionally inserting them as new records.
// Appends related to o.R.HeadRefSubmissions.
// Sets related.R.HeadRef appropriately.
func (o *Ref) AddHeadRefSubmissions(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Submission) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.HeadRefID, o.ID)
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"submissions\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"head_ref_id"}),
				strmangle.WhereClause("\"", "\"", 2, submissionPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.HeadRefID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &refR{
			HeadRefSubmissions: related,
		}
	} else {
		o.R.HeadRefSubmissions = append(o.R.HeadRefSubmissions, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &submissionR{
				HeadRef: o,
			}
		} else {
			rel.R.HeadRef = o
		}
	}
	return nil
}

// SetHeadRefSubmissions removes all previously related items of the
// ref replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.HeadRef's HeadRefSubmissions accordingly.
// Replaces o.R.HeadRefSubmissions with related.
// Sets related.R.HeadRef's HeadRefSubmissions accordingly.
func (o *Ref) SetHeadRefSubmissions(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Submission) error {
	query := "update \"submissions\" set \"head_ref_id\" = null where \"head_ref_id\" = $1"
	values := []interface{}{o.ID}
	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, query)
		fmt.Fprintln(writer, values)
	}
	_, err := exec.ExecContext(ctx, query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.HeadRefSubmissions {
			queries.SetScanner(&rel.HeadRefID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.HeadRef = nil
		}

		o.R.HeadRefSubmissions = nil
	}
	return o.AddHeadRefSubmissions(ctx, exec, insert, related...)
}

// RemoveHeadRefSubmissions relationships from objects passed in.
// Removes related items from R.HeadRefSubmissions (uses pointer comparison, removal does not keep order)
// Sets related.R.HeadRef.
func (o *Ref) RemoveHeadRefSubmissions(ctx context.Context, exec boil.ContextExecutor, related ...*Submission) error {
	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.HeadRefID, nil)
		if rel.R != nil {
			rel.R.HeadRef = nil
		}
		if _, err = rel.Update(ctx, exec, boil.Whitelist("head_ref_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.HeadRefSubmissions {
			if rel != ri {
				continue
			}

			ln := len(o.R.HeadRefSubmissions)
			if ln > 1 && i < ln-1 {
				o.R.HeadRefSubmissions[i] = o.R.HeadRefSubmissions[ln-1]
			}
			o.R.HeadRefSubmissions = o.R.HeadRefSubmissions[:ln-1]
			break
		}
	}

	return nil
}

// Refs retrieves all the records using an executor.
func Refs(mods ...qm.QueryMod) refQuery {
	mods = append(mods, qm.From("\"refs\""))
	return refQuery{NewQuery(mods...)}
}

// FindRef retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindRef(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*Ref, error) {
	refObj := &Ref{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"refs\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, refObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from refs")
	}

	return refObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Ref) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no refs provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(refColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	refInsertCacheMut.RLock()
	cache, cached := refInsertCache[key]
	refInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			refAllColumns,
			refColumnsWithDefault,
			refColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(refType, refMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(refType, refMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"refs\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"refs\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into refs")
	}

	if !cached {
		refInsertCacheMut.Lock()
		refInsertCache[key] = cache
		refInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Ref.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Ref) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	refUpdateCacheMut.RLock()
	cache, cached := refUpdateCache[key]
	refUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			refAllColumns,
			refPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update refs, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"refs\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, refPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(refType, refMapping, append(wl, refPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update refs row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for refs")
	}

	if !cached {
		refUpdateCacheMut.Lock()
		refUpdateCache[key] = cache
		refUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q refQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for refs")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for refs")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o RefSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), refPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"refs\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, refPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in ref slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all ref")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Ref) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no refs provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(refColumnsWithDefault, o)

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

	refUpsertCacheMut.RLock()
	cache, cached := refUpsertCache[key]
	refUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			refAllColumns,
			refColumnsWithDefault,
			refColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			refAllColumns,
			refPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert refs, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(refPrimaryKeyColumns))
			copy(conflict, refPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"refs\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(refType, refMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(refType, refMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert refs")
	}

	if !cached {
		refUpsertCacheMut.Lock()
		refUpsertCache[key] = cache
		refUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Ref record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Ref) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Ref provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), refPrimaryKeyMapping)
	sql := "DELETE FROM \"refs\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from refs")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for refs")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q refQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no refQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from refs")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for refs")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o RefSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(refBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), refPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"refs\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, refPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from ref slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for refs")
	}

	if len(refAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Ref) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindRef(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *RefSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := RefSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), refPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"refs\".* FROM \"refs\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, refPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in RefSlice")
	}

	*o = slice

	return nil
}

// RefExists checks if the Ref row exists.
func RefExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"refs\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if refs exists")
	}

	return exists, nil
}
