package main

import (
	"context"
	"fmt"
)

type contextKey int

func WithStringValue(val string) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, contextKey(0), val)
	}
}

func StringValueFromContext(ctx context.Context) string {
	return ctx.Value(contextKey(0)).(string)
}

func main() {
	var withContext = WithContextCompose(
		WithStringValue("hello"),
		WithStringValue("ni hao"),
	)
	ctx := withContext(context.Background())

	fmt.Println(StringValueFromContext(ctx))
	fmt.Println(ctx)
}

type WithContext = func(ctx context.Context) context.Context

func WithContextCompose(withContexts ...WithContext) WithContext {
	return func(ctx context.Context) context.Context {
		for i := range withContexts {
			ctx = withContexts[i](ctx)
		}
		return ctx
	}
}
