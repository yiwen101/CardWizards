package store

import (
	"testing"

	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/pkg/utils"
)

func init() {
	InfoStore.Load("test1", utils.PkgToIDL, "")
}
func TestPRoxyLevelCommands(t *testing.T) {
	m, err := InfoStore.GetAllServiceNames()
	test.Assert(t, err == nil, err)
	test.Assert(t, len(m) == 4, "no service found")
	isOn, err := InfoStore.CheckProxyStatus()
	test.Assert(t, err == nil, err)
	test.Assert(t, isOn, "proxy should be on")
	err = InfoStore.TurnOffProxy()
	test.Assert(t, err == nil, err)
	isOn, _ = InfoStore.CheckProxyStatus()
	test.Assert(t, !isOn, "proxy should be off")
	InfoStore.TurnOnProxy()
	isOn, _ = InfoStore.CheckProxyStatus()
	test.Assert(t, isOn, "proxy should be on again")
	a1, err := InfoStore.GetProxyAddress()
	test.Assert(t, err == nil, err)
	test.Assert(t, a1 == "test1", "proxy address should be test1")

}

/*
	GetAPIs(serviceName string) (map[string]*ApiMeta, error)
	AddService(idlFileName, clusterName string) error                 //validator, client, descriptor, router
	UpdateService(serviceName, idlFileName, clusterName string) error //validator, client, descriptor, router
	RemoveService(serviceName string) error                           //validator, client, descriptor, router
	TurnOnService(serviceName string) error
	TurnOffService(serviceName string) error
	GetServiceInfo(serviceName string) (*ServiceMeta, error)
*/

func TestServiceLevelCommands(t *testing.T) {
	_, err := InfoStore.GetAPIs("test")
	test.Assert(t, err != nil, err)
	m, err := InfoStore.GetAPIs("arithmetic")
	test.Assert(t, err == nil, err)
	test.Assert(t, len(m) == 5, "no api found")
	_, err = InfoStore.AddService("test", "test")
	test.Assert(t, err != nil, err)
	_, err = InfoStore.AddService("arithmetic.thrift", "l")
	test.Assert(t, err != nil, err)

	err = InfoStore.RemoveService("test")
	test.Assert(t, err != nil, err)

	_, err = InfoStore.UpdateService("test", "test", "test")
	test.Assert(t, err != nil, err)
	_, err = InfoStore.UpdateService("arithmetic", "arithmetic.thrift", "l")
	test.Assert(t, err == nil, err)
	err = InfoStore.RemoveService("arithmetic")
	test.Assert(t, err == nil, err)
	_, err = InfoStore.GetServiceInfo("arithmetic")
	test.Assert(t, err != nil, err)
	_, err = InfoStore.AddService("arithmetic.thrift", "l")
	test.Assert(t, err == nil, err)
	_, err = InfoStore.GetServiceInfo("arithmetic")
	test.Assert(t, err == nil, err)
}

/*
	    CheckAPIStatus(serviceName, methodName string) (*ApiMeta, error)
		TurnOnAPI(serviceName, methodName string) error                                       //APIGate
		TurnOffAPI(serviceName, methodName string) error                                      //APIGate
		TurnOnValidation(serviceName, methodName string) error                                //validator
		TurnOffValidation(serviceName, methodName string) error                               //validator
		AddRoute(serviceName, methodName, url, httpMethod string) error                       //router
		ModifyRoute(serviceName, methodName, url, httpMethod, newUrl, newMethod string) error //router
		RemoveRoute(serviceName, methodName, url, httpMethod string) error                    //router
		GetRoutes(serviceName, methodName string) (map[string]map[string]bool, error)
		GetLbType(serviceName string) (string, error) //lb
		SetLbType(serviceName, lbType string) error   //lb
*/
func TestAPILevelCommands(t *testing.T) {
	api, err := InfoStore.CheckAPIStatus("arithmetic", "Add")
	test.Assert(t, err == nil, err)
	test.Assert(t, api.IsOn == true, "api gate should be off")
	err = InfoStore.TurnOffAPI("arithmetic", "Add")
	test.Assert(t, err == nil, err)
	api, err = InfoStore.CheckAPIStatus("arithmetic", "Add")
	test.Assert(t, err == nil, err)
	test.Assert(t, api.IsOn == false, "api gate should be off")
	err = InfoStore.TurnOnAPI("arithmetic", "Add")
	test.Assert(t, err == nil, err)
	api, err = InfoStore.CheckAPIStatus("arithmetic", "Add")
	test.Assert(t, err == nil, err)
	test.Assert(t, api.IsOn == true, "api gate should be off")
	test.Assert(t, api.ValidationOn == false, "validation is off by default")
	err = InfoStore.TurnOnValidation("arithmetic", "Add")
	test.Assert(t, err == nil, err)
	api, err = InfoStore.CheckAPIStatus("arithmetic", "Add")
	test.Assert(t, err == nil, err)
	test.Assert(t, api.ValidationOn == true, "validation should be on")
	err = InfoStore.TurnOffValidation("arithmetic", "Add")
	test.Assert(t, err == nil, err)
	api, err = InfoStore.CheckAPIStatus("arithmetic", "Add")
	test.Assert(t, err == nil, err)
	test.Assert(t, api.ValidationOn == false, "validation should be off")

	routes, err := InfoStore.GetRoutes("arithmetic", "Add")
	test.Assert(t, err == nil, err)
	test.Assert(t, routes["GET"]["/arithmetic/Add"], "default")
	err = InfoStore.ModifyRoute("arithmetic", "Add", "GET", "/arithmetic/Add", "POST", "/test")
	test.Assert(t, err == nil, err)
	test.Assert(t, !routes["GET"]["/arithmetic/Add"], "deleted")
	test.Assert(t, routes["POST"]["/test"], "added")
}

type testHandeler struct {
	value bool
}

func (t *testHandeler) OnStatechanged(data ...interface{}) error {
	isOn := data[0].(bool)
	t.value = isOn
	return nil
}

func TestNotification(t *testing.T) {
	th := &testHandeler{value: false}

	InfoStore.RegisterProxyStateListener(th)
	InfoStore.TurnOnProxy()
	test.Assert(t, th.value == true, "proxy should be on")
	InfoStore.TurnOffProxy()
	test.Assert(t, th.value == false, "proxy should be off")
}
