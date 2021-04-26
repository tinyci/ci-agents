package protoconv

import "context"

// ToProto converts a db/models type into a protobuf type. Pass a pointer. You'll get a pointer back.
func (c *Converter) ToProto(ctx context.Context, i interface{}) (interface{}, error) {
	return c.convert(ctx, toProto, i)
}

// FromProto converts in the other direction. Pass a pointer, get a pointer.
func (c *Converter) FromProto(ctx context.Context, i interface{}) (interface{}, error) {
	return c.convert(ctx, fromProto, i)
}
