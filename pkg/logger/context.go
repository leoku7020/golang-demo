package logger

import (
	"context"

	"demo/pkg/ctxkit"
)

// ContextWithFields injects Fields into returned context
func ContextWithFields(c context.Context, fs Fields) context.Context {
	return ctxkit.HSet(c, ctxkit.KeyLogger, ctxkit.KV(fs))
}

func fields(c context.Context) (Fields, bool) {
	kv, ok := ctxkit.HGetAll(c, ctxkit.KeyLogger)
	if !ok {
		return nil, false
	}

	return Fields(kv), true
}
