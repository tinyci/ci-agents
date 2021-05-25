package protoconv

import (
	"context"
	"fmt"
	"reflect"

	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/ci-gen/openapi/services/uisvc"
)

var (
	toProto   int = 0
	fromProto int = 1
)

type conversionFunc func(context.Context, interface{}) (interface{}, error)

// APIConverter is a handle into our conversion system.
type APIConverter struct {
	converters convertMap
}

type (
	reflectMap map[reflect.Type]conversionFunc
	convertMap [2]reflectMap
)

// New creates a new *Converter.
func New() *APIConverter {
	c := &APIConverter{}
	c.registerConversion(toProto, &uisvc.ModelSubmission{}, subToProto)
	c.registerConversion(fromProto, &types.Submission{}, subFromProto)
	c.registerConversion(toProto, &uisvc.User{}, userToProto)
	c.registerConversion(fromProto, &types.User{}, userFromProto)
	c.registerConversion(toProto, &uisvc.Ref{}, refToProto)
	c.registerConversion(fromProto, &types.Ref{}, refFromProto)
	c.registerConversion(toProto, &uisvc.Repository{}, repoToProto)
	c.registerConversion(fromProto, &types.Repository{}, repoFromProto)
	c.registerConversion(toProto, &uisvc.Run{}, runToProto)
	c.registerConversion(fromProto, &types.Run{}, runFromProto)
	c.registerConversion(toProto, &uisvc.Task{}, taskToProto)
	c.registerConversion(fromProto, &types.Task{}, taskFromProto)
	c.registerConversion(toProto, &uisvc.UserError{}, ueToProto)
	c.registerConversion(fromProto, &types.UserError{}, ueFromProto)
	return c
}

func (c *APIConverter) registerConversion(direction int, i interface{}, fun conversionFunc) {
	if c.converters[direction] == nil {
		c.converters[direction] = reflectMap{}
	}

	c.converters[direction][reflect.TypeOf(i)] = fun
}

func (c *APIConverter) convert(ctx context.Context, direction int, i interface{}) (interface{}, error) {
	fun, ok := c.converters[direction][reflect.TypeOf(i)]
	if !ok {
		return nil, fmt.Errorf("invalid conversion for type %T", i)
	}

	return fun(ctx, i)
}
