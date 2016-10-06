// Package ctxtg provides Context which should be passed between all our services.
// Also ctxtg provides Token Signer, Parser and test helper implementations in ctxtgtest.
package ctxtg

import (
	"time"

	"context"
)

// Keys to get data from context.Context
type key int

// Values to access ctxtg values inside context.Context
const (
	_ key = iota

	TokenKey
	TracingIDKey
	DataKey
)

// Token represents JWT token
type Token string

// Context for microservices communication
type Context struct {
	Token     Token
	Deadline  int64
	TracingID string
	Data      map[string]interface{}
}

// ToContext convert to context.Context object
func (c *Context) ToContext() (context.Context, context.CancelFunc) {
	ctx, cancel := contextFromDeadline(c.Deadline)
	ctx = context.WithValue(ctx, TokenKey, c.Token)
	ctx = context.WithValue(ctx, TracingIDKey, c.TracingID)
	if c.Data != nil {
		ctx = context.WithValue(ctx, DataKey, c.Data)
	}
	return ctx, cancel
}

// FromContext convert context.Context to Context correctly extracting required fields
func FromContext(ctx context.Context) Context {
	return Context{
		Token:     tokenValue(ctx),
		Deadline:  unixDeadline(ctx),
		TracingID: stringValue(ctx, TracingIDKey),
		Data:      DataFromContext(ctx),
	}
}

// WithDataValue add key-value to Data map inside context.Context and return new context.Context
func WithDataValue(parent context.Context, key string, value interface{}) context.Context {
	if d := DataFromContext(parent); d != nil {
		d[key] = value
		return parent
	}
	return context.WithValue(parent, DataKey, map[string]interface{}{key: value})
}

// ValueFromData returns value from Data map inside ctx or nil
func ValueFromData(ctx context.Context, key string) interface{} {
	if m := DataFromContext(ctx); m != nil {
		return m[key]
	}
	return nil
}

// DataFromContext returns map from Data key from ctx or nil
func DataFromContext(ctx context.Context) map[string]interface{} {
	v := ctx.Value(DataKey)
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	return nil
}

func unixDeadline(ctx context.Context) int64 {
	if t, ok := ctx.Deadline(); ok {
		return t.Unix()
	}
	return 0
}

func tokenValue(ctx context.Context) Token {
	v := ctx.Value(TokenKey)
	if token, ok := v.(Token); ok {
		return token
	}
	return Token("")
}

func stringValue(ctx context.Context, key interface{}) string {
	v := ctx.Value(key)
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

func contextFromDeadline(deadline int64) (context.Context, context.CancelFunc) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if deadline <= 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithDeadline(context.Background(), time.Unix(deadline, 0))
	}
	return ctx, cancel
}
