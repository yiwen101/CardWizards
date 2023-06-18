package temp

import (
	"net/http"

	"github.com/cloudwego/kitex/pkg/generic/descriptor"
)

const (
	RelativePathToIDL = "../IDL/"
	RefaultHttpMethod = http.MethodPost
	RefaultRoute      = "/:serviceName/:methodName/*furtherRoutes"
)

func validateType(t *descriptor.TypeDescriptor, json interface{}) bool {
	switch t.Elem.Type.ToThriftTType().String() {
	case "STRUCT":
		json, ok := json.(map[string]interface{})
		if !ok {
			return false
		}
		params := t.Struct.FieldsByID
		for _, p := range params {
			ok = validateName(p, json)

			if !ok {
				return false
			}
		}
		return true

	case "BOOL":
		_, ok := json.(bool)
		return ok
	case "BYTE":
		_, ok := json.(int8)
		return ok
	case "DOUBLE":
		_, ok := json.(float64)
		return ok
	case "I16":
		_, ok := json.(int16)
		return ok
	case "I32":
		_, ok := json.(int32)
		return ok
	case "I64":
		_, ok := json.(int64)
		return ok
	case "STRING":
		_, ok := json.(string)
		return ok
	case "MAP":
		keyType := t.Key
		valueType := t.Elem
		json, _ := json.(map[interface{}]interface{})
		for key, value := range json {
			ok := validateType(keyType, key)
			if !ok {
				return false
			}
			ok = validateType(valueType, value)
			if !ok {
				return false
			}
			return true
		}
	}
	return true
}

func validateName(p *descriptor.FieldDescriptor, json map[string]interface{}) bool {

	if p.Optional || p.DefaultValue != nil {
		return true
	}

	var v interface{}

	if json[p.Name] == nil {
		if json[p.Alias] == nil {
			return false
		} else {
			v = json[p.Alias]
		}
	} else {
		v = json[p.Name]
	}
	return validateType(p.Type, v)
}
func validateBody(fuc *descriptor.FunctionDescriptor, json interface{}) bool {
	params := fuc.Request.Struct.FieldsByID
	if params[2] != nil {
		// invalid number of fields, only parameter is request
		// benefit of doubt
		return true
	}
	return validateType(params[1].Type, json)
}

func validateMethod(desKeeper descriptorKeeper, methodName string) (*descriptor.FunctionDescriptor, bool) {
	sevDes := desKeeper.get()
	fuc, err := sevDes.LookupFunctionByMethod(methodName)
	return fuc, err == nil
}

func validateJson(methodName string, json map[string]interface{}) {
	descriptorKeeper, ok := serviceToDescriptorMap[methodName]

}
