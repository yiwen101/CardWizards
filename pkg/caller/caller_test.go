package caller

import (
	"testing"

	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/pkg/store"
)

func TestCallerClientsUpdatesMethods(t *testing.T) {
	store.InfoStore.Load("", "", "../../testing/idl")
	client, err := GetClient("arithmetic")
	// debug mode peek: client should be with weighted round robin lb
	test.Assert(t, err == nil)
	test.Assert(t, client != nil)
	store.InfoStore.SetLbType("arithmetic", "random")
	client, err = GetClient("arithmetic")
	test.Assert(t, err == nil)
	test.Assert(t, client != nil)
	// debug mode peek: client should be with random lb
	_, err = GetClient("false")
	test.Assert(t, err != nil)
	meta, err := store.InfoStore.GetServiceInfo("arithmetic")
	opts, err := getOptionsFor(meta)
	test.Assert(t, err == nil)
	test.Assert(t, len(opts) == 2)
	store.InfoStore.RemoveService("arithmetic")
	_, err = GetClient("arithmetic")
	test.Assert(t, err != nil)
}
