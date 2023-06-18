package validate

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
   desc "github.com/yiwen101/CardWizards/service/descriptor"
)

func validateType(t *descriptor.TypeDescriptor, json interface{}) error {

	switch t.Elem.Type.ToThriftTType().String() {
	case "STRUCT":
		json, ok := json.(map[string]interface{})
		if !ok {
			return fmt.Errorf("type mismatch, expected struct, got %T", json)
		}
		params := t.Struct.FieldsByID
		for _, p := range params {
			err := validateName(p, json)

			if err != nil {
				return err
			}
		}
		return nil

	case "BOOL":
		_, ok := json.(bool)
		if ok {
			return nil
		} else {
			return fmt.Errorf("type mismatch, expected bool, got %T", json)
		}
	case "BYTE":
		_, ok := json.(int8)
		if ok {
			return nil
		} else {
			return fmt.Errorf("type mismatch, expected byte, got %T", json)
		}
	case "DOUBLE":
		_, ok := json.(float64)
		if ok {
			return nil
		} else {
			return fmt.Errorf("type mismatch, expected double, got %T", json)
		}
	case "I16":
		_, ok := json.(int16)
		if ok {
			return nil
		} else {
			return fmt.Errorf("type mismatch, expected i16, got %T", json)
		}
	case "I32":
		_, ok := json.(int32)
		if ok {
			return nil
		} else {
			return fmt.Errorf("type mismatch, expected i32, got %T", json)
		}
	case "I64":
		_, ok := json.(int64)
		if ok {
			return nil
		} else {
			return fmt.Errorf("type mismatch, expected i64, got %T", json)
		}
	case "STRING":
		_, ok := json.(string)
		if ok {
			return nil
		} else {
			return fmt.Errorf("type mismatch, expected string, got %T", json)
		}
	case "MAP":
		keyType := t.Key
		valueType := t.Elem
		json, _ := json.(map[interface{}]interface{})
		for key, value := range json {
			err := validateType(keyType, key)
			if err != nil {
				return err
			}
			err = validateType(valueType, value)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func validateName(p *descriptor.FieldDescriptor, json map[string]interface{}) error {

	if p.Optional || p.DefaultValue != nil {
		return nil
	}

	var v interface{}

	if json[p.Name] == nil {
		if json[p.Alias] == nil {
			return fmt.Errorf("missing required field %s", p.Name)
		} else {
			v = json[p.Alias]
		}
	} else {
		v = json[p.Name]
	}
	return validateType(p.Type, v)
}

func validateBody(fuc *descriptor.FunctionDescriptor, json interface{}) error {
	params := fuc.Request.Struct.FieldsByID
	if params[2] != nil {
		// invalid number of fields, only parameter is request
		// benefit of doubt
		return nil
	}
	return validateType(params[1].Type, json)
}

func isGenericRoute(serviceName, methodName string) error {
	return desc.DescsManager.ValidateServiceAndMethodName(serviceName, methodName)
}

func isAnnotatedRoute(req *descriptor.HTTPRequest) (string, error) {
	return desc.DescsManager.ValidateServiceAndMethodNameWithAnnotedRoutes(req)
}

func treatJsonBody(ctx context.Context, c *app.RequestContext) (map[string]interface{}, error) {
	if string(c.ContentType()) != "application/json" {
		return nil, fmt.Errorf("Invalid Content-Type, expected application/json")
	}

	body, err := c.Body()
	if err != nil {
		return nil, err
	}

	var j map[string]interface{}

	err = json.Unmarshal(body, &j)
	if err != nil {
		return nil, err
	}
	return j, nil
}
