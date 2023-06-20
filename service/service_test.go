package service

import (
	"context"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/thriftgo/pkg/test"
)

func handlerCacheTest(t *testing.T) {
	h := handlerCache{}
	f, ok := h.get("test", "test")
	test.Assert(t, ok == false)
	test.Assert(t, f == nil)
	testEffect := 0
	theFunc := func(ctx context.Context, c *app.RequestContext) { testEffect++ }
	h.save("test", "test", theFunc)
	f, ok = h.get("test", "test")
	test.Assert(t, ok == true)
	test.Assert(t, f != nil)
	f(context.Background(), nil)
	test.Assert(t, testEffect == 1)
}


