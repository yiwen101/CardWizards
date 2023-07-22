package filter

/*
import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/yiwen101/CardWizards/pkg/store"
)

type handler func(ctx context.Context, c *app.RequestContext) bool

func (h handler) Handle(ctx context.Context, c *app.RequestContext) bool {
	return h(ctx, c)
}

type filter interface {
	Handle(ctx context.Context, c *app.RequestContext)
	AddNext(f filter)
}

type filterImple struct {
	Handler handler
	next    filter
}

func (f *filterImple) AddNext(next filter) {
	f.next = next
}

func (f *filterImple) Handle(ctx context.Context, c *app.RequestContext) {
	proceed := f.Handler.Handle(ctx, c)
	if proceed && f.next != nil {
		f.next.Handle(ctx, c)
	}
}

func newFilter(handler func(ctx context.Context, c *app.RequestContext)) filter {
	return &filterImple{
		Handler: handler,
		next:    nil,
	}
}

func newFilterChain(handlers ...func(ctx context.Context, c *app.RequestContext) bool) filter {
	if len(handlers) == 0 {
		return nil
	}
	filters := make([]filter, len(handlers))
	for i, handler := range handlers {
		filters[i] = newFilter(handler)
	}
	for i := 0; i < len(filters)-1; i++ {
		filters[i].AddNext(filters[i+1])
	}
	return filters[0]
}

func ProxyGate(ctx context.Context, c *app.RequestContext) bool {
	return store.ProxyIsOn()
}


func ApiGate(ctx context.Context, c *app.RequestContext) bool {
	return store.ApiIsOn(c.ServiceName, c.MethodName)
}
*/
