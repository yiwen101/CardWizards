package proxy

import (
	"bytes"
	"context"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/cloudwego/thriftgo/pkg/test"
	myRouter "github.com/yiwen101/CardWizards/pkg/router"
	"github.com/yiwen101/CardWizards/pkg/store"
)

// // turn on the kitex server and nacos server before run this test
func TestPerformRequest(t *testing.T) {
	store.InfoStore.Load("", "../../testing/idl", "")
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/*:test",
		func(ctx context.Context, c *app.RequestContext) {
			Proxy.Serve(ctx, c, nil)
		})

	body := bytes.NewBufferString("{}")
	len := body.Len()
	w := ut.PerformRequest(router, "GET", "/arithmetic/Add", &ut.Body{Body: body, Len: len},
		ut.Header{Key: "Content-Type", Value: "application/json"})
	resp := w.Result()
	var j map[string]interface{}
	err := sonic.Unmarshal(resp.Body(), &j)
	test.Assert(t, err == nil, err)
	fA, ok := j["firstArguement"].(float64)
	test.Assert(t, ok)
	test.Assert(t, fA == 0)
	w = ut.PerformRequest(router, "GET", "/arithmetic/Add", &ut.Body{Body: body, Len: len}, ut.Header{})
	resp = w.Result()

	assert.DeepEqual(t, 400, resp.StatusCode())
	realBody := string(resp.Body())
	assert.DeepEqual(t, "Invalid Content-Type: ", realBody)

	w = ut.PerformRequest(router, "GET", "/wrong/route", &ut.Body{Body: body, Len: len},
		ut.Header{Key: "Content-Type", Value: "application/json"})
	resp = w.Result()
	assert.DeepEqual(t, 404, resp.StatusCode())
	realBody = string(resp.Body())
	assert.DeepEqual(t, "404 page not found", realBody)
}

func TestRunTimeHandlerChange(t *testing.T) {
	store.InfoStore.Load("", "../../testing/idl", "")
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/*:test",
		func(ctx context.Context, c *app.RequestContext) {
			Proxy.Serve(ctx, c, nil)
		})

	store.InfoStore.TurnOnValidation("arithmetic", "Add")
	body := bytes.NewBufferString("{}")
	len := body.Len()
	w := ut.PerformRequest(router, "GET", "/arithmetic/Add", &ut.Body{Body: body, Len: len},
		ut.Header{Key: "Content-Type", Value: "application/json"})
	resp := w.Result()
	code := resp.StatusCode()
	test.Assert(t, code == 400)

	store.InfoStore.TurnOffAPI("arithmetic", "Add")
	w = ut.PerformRequest(router, "GET", "/arithmetic/Add", &ut.Body{Body: body, Len: len},
		ut.Header{Key: "Content-Type", Value: "application/json"})
	resp, code = w.Result(), w.Result().StatusCode()
	test.Assert(t, code == 404)

	store.InfoStore.TurnOnAPI("arithmetic", "Add")
	err := store.InfoStore.ModifyRoute("arithmetic", "Add", "GET", "/arithmetic/Add", "GET", "/test/test2")
	test.Assert(t, err == nil)
	_, ok := myRouter.GetRoute("GET", "/test/test2")
	test.Assert(t, ok)
	w = ut.PerformRequest(router, "GET", "/test/test2", &ut.Body{Body: body, Len: len},
		ut.Header{Key: "Content-Type", Value: "application/json"})
	resp, code = w.Result(), w.Result().StatusCode()
	test.Assert(t, code == 400)
}
