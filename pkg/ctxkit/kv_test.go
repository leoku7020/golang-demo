package ctxkit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type kvSuite struct {
	suite.Suite
}

func (s *kvSuite) SetupSuite()    {}
func (s *kvSuite) TearDownSuite() {}
func (s *kvSuite) SetupTest()     {}
func (s *kvSuite) TearDownTest()  {}

func TestKVSuite(t *testing.T) {
	suite.Run(t, new(kvSuite))
}

func (s *kvSuite) TestGet() {
	tests := []struct {
		Desc      string
		SetupTest func() context.Context
		Ctx       context.Context
		CtxKey    Key
		ExpBool   bool
		ExpKV     KV
	}{
		{
			Desc: "not existed",
			SetupTest: func() context.Context {
				return context.TODO()
			},
			CtxKey:  KeyNone, // no such key
			ExpBool: false,
			ExpKV:   nil,
		},
		{
			Desc: "normal case",
			SetupTest: func() context.Context {
				return context.WithValue(context.TODO(), KeyLogger, structKV{m: map[string]interface{}{
					"foo":    "bar",
					"number": 123,
				}})
			},
			CtxKey:  KeyLogger,
			ExpBool: true,
			ExpKV:   KV{"foo": "bar", "number": 123},
		},
	}

	for _, t := range tests {
		if t.SetupTest != nil {
			t.Ctx = t.SetupTest()
		}

		kv, b := HGetAll(t.Ctx, t.CtxKey)
		s.Require().Equal(t.ExpBool, b, t.Desc)
		s.Require().Equal(t.ExpKV, kv, t.Desc)

		s.TearDownTest()
	}
}

func (s *kvSuite) TestGetValue() {
	tests := []struct {
		Desc      string
		SetupTest func() context.Context
		Ctx       context.Context
		CtxKey    Key
		Key       string
		ExpBool   bool
		ExpValue  interface{}
	}{
		{
			Desc: "not existed",
			SetupTest: func() context.Context {
				return context.TODO()
			},
			CtxKey:   KeyNone, // no such key
			ExpBool:  false,
			ExpValue: nil,
		},
		{
			Desc: "wrong key",
			SetupTest: func() context.Context {
				return context.WithValue(context.TODO(), KeyLogger, structKV{m: map[string]interface{}{
					"foo":    "bar",
					"number": 123,
				}})
			},
			CtxKey:   KeyLogger,
			Key:      "wrong",
			ExpBool:  false,
			ExpValue: nil,
		},
		{
			Desc: "normal case",
			SetupTest: func() context.Context {
				return context.WithValue(context.TODO(), KeyLogger, structKV{m: map[string]interface{}{
					"foo":    "bar",
					"number": 123,
				}})
			},
			CtxKey:   KeyLogger,
			Key:      "foo",
			ExpBool:  true,
			ExpValue: "bar",
		},
	}

	for _, t := range tests {
		if t.SetupTest != nil {
			t.Ctx = t.SetupTest()
		}

		v, b := HGet(t.Ctx, t.CtxKey, t.Key)
		s.Require().Equal(t.ExpBool, b, t.Desc)
		s.Require().Equal(t.ExpValue, v, t.Desc)

		s.TearDownTest()
	}
}

func (s *kvSuite) TestSet() {
	tests := []struct {
		Desc   string
		Ctx    context.Context
		CtxKey Key
		KV     KV
		ExpKV  KV
	}{
		{
			Desc:   "first time with nil",
			Ctx:    context.TODO(),
			CtxKey: KeyLogger,
			KV:     nil,
			ExpKV:  nil,
		},
		{
			Desc:   "first time with KV",
			Ctx:    context.TODO(),
			CtxKey: KeyLogger,
			KV:     KV{"foo": "bar", "number": 123},
			ExpKV:  KV{"foo": "bar", "number": 123},
		},
		{
			Desc: "second time with nil",
			Ctx: context.WithValue(context.TODO(), KeyLogger, structKV{m: map[string]interface{}{
				"foo":    "bar",
				"number": 123,
			}}),
			CtxKey: KeyLogger,
			KV:     nil,
			ExpKV:  KV{"foo": "bar", "number": 123},
		},
		{
			Desc: "second time with KV",
			Ctx: context.WithValue(context.TODO(), KeyLogger, structKV{m: map[string]interface{}{
				"foo":    "bar",
				"number": 123,
			}}),
			CtxKey: KeyLogger,
			KV:     KV{"foo": "bar", "bool": true},
			ExpKV:  KV{"foo": "bar", "number": 123, "bool": true}, // merged
		},
	}

	for _, t := range tests {
		c := HSet(t.Ctx, t.CtxKey, t.KV)
		kv, b := HGetAll(c, KeyLogger)
		s.Require().True(b)
		s.Require().Equal(t.ExpKV, kv, t.Desc)

		s.TearDownTest()
	}
}
