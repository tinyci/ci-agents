package protoconv

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/db/models"
)

var (
	toProto   int = 0
	fromProto int = 1
)

type conversionFunc func(context.Context, *sql.DB, interface{}) (interface{}, error)

// Converter is a handle into our conversion system. Requires a database object.
type Converter struct {
	converters convertMap
	db         *sql.DB
}

type (
	reflectMap map[reflect.Type]conversionFunc
	convertMap [2]reflectMap
)

// New creates a new *Converter.
func New(db *sql.DB) *Converter {
	c := &Converter{db: db}
	c.registerConversion(toProto, &models.OAuth{}, oauthToProto)
	c.registerConversion(fromProto, &types.Ref{}, refFromProto)
	c.registerConversion(toProto, &models.Ref{}, refToProto)
	c.registerConversion(fromProto, &types.Repository{}, repoFromProto)
	c.registerConversion(toProto, &models.Repository{}, repoToProto)
	c.registerConversion(fromProto, &types.User{}, userFromProto)
	c.registerConversion(toProto, &models.User{}, userToProto)
	c.registerConversion(fromProto, &types.Submission{}, subFromProto)
	c.registerConversion(toProto, &models.Submission{}, subToProto)
	c.registerConversion(fromProto, &types.QueueItem{}, queueItemFromProto)
	c.registerConversion(toProto, &models.QueueItem{}, queueItemToProto)
	c.registerConversion(fromProto, &types.Task{}, taskFromProto)
	c.registerConversion(toProto, &models.Task{}, taskToProto)
	c.registerConversion(fromProto, &types.Run{}, runFromProto)
	c.registerConversion(toProto, &models.Run{}, runToProto)
	c.registerConversion(fromProto, &types.UserError{}, userErrorFromProto)
	c.registerConversion(toProto, &models.UserError{}, userErrorToProto)
	c.registerConversion(fromProto, &types.Session{}, sessionFromProto)
	c.registerConversion(toProto, &models.Session{}, sessionToProto)
	return c
}

func (c *Converter) registerConversion(direction int, i interface{}, fun conversionFunc) {
	if c.converters[direction] == nil {
		c.converters[direction] = reflectMap{}
	}

	c.converters[direction][reflect.TypeOf(i)] = fun
}

func (c *Converter) convert(ctx context.Context, direction int, i interface{}) (interface{}, error) {
	fun, ok := c.converters[direction][reflect.TypeOf(i)]
	if !ok {
		return nil, fmt.Errorf("invalid conversion for type %T", i)
	}

	return fun(ctx, c.db, i)
}
