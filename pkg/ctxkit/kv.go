package ctxkit

import (
	"context"
)

// KV is the data struct to store key-value pairs in the context.Context
type KV map[string]interface{}

func (kv *KV) clone() KV {
	m := KV{}
	for k, v := range *kv {
		m[k] = v
	}

	return m
}

// structKV wraps the KV insides. It's the hybird structure to keep performance and solve racing issue.
// * KV inside keeps each value in it can be queried by the associated key in O(1) time.
// * structKV outside keeps accessing it with Context.Value() would return you a copy of the struct.
type structKV struct {
	m map[string]interface{}
}

func (m *structKV) Get(key string) interface{} {
	return m.m[key]
}

// HSet sets KV into ctx.
func HSet(ctx context.Context, ctxKey Key, values KV) context.Context {
	kv, ok := HGetAll(ctx, ctxKey)
	if !ok {
		// not existed before, return a new one.
		return context.WithValue(ctx, ctxKey, structKV{m: values})
	}

	m := kv.clone()
	// override values in kv with `values`
	for k, v := range values {
		m[k] = v
	}

	return context.WithValue(ctx, ctxKey, structKV{m: m})
}

// HGetAll collects KV with specified ctxKey.
func HGetAll(ctx context.Context, ctxKey Key) (KV, bool) {
	m, ok := ctx.Value(ctxKey).(structKV)
	if !ok {
		return nil, false
	}

	return m.m, true
}

// HGet collects KV with specified ctxKey, and find the value with key
func HGet(ctx context.Context, ctxKey Key, key string) (interface{}, bool) {
	kv, ok := HGetAll(ctx, ctxKey)
	if !ok {
		return nil, false
	}

	v, ok := kv[key]
	if !ok {
		return nil, false
	}

	return v, true
}
