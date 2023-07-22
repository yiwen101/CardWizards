package caller

import (
	"testing"

	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/pkg/store"
)

func TestCallerClientsUpdatesMethods(t *testing.T) {
	store.InfoStore.Load("", "", "../../testing/idl")
	client, ok := GetClient("arithmetic")
	// debug mode peek: client should be with weighted round robin lb
	test.Assert(t, ok)
	test.Assert(t, client != nil)
	store.InfoStore.SetLbType("arithmetic", "random")
	client, ok = GetClient("arithmetic")
	test.Assert(t, ok)
	test.Assert(t, client != nil)
	// debug mode peek: client should be with random lb
	_, ok = GetClient("false")
	test.Assert(t, !ok)
	meta, err := store.InfoStore.GetServiceInfo("arithmetic")
	opts, err := getOptionsFor(meta)
	test.Assert(t, err == nil)
	test.Assert(t, len(opts) == 2)
	store.InfoStore.RemoveService("arithmetic")
	_, ok = GetClient("arithmetic")
	test.Assert(t, !ok)
}
