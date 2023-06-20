package validate

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
)

type validatorImplement struct {
	functionDescriptor *descriptor.FunctionDescriptor
}

func validateType(t *descriptor.TypeDescriptor, json interface{}, fieldName string) error {
	//  it seems that all number encoded in json format are float64

	switch t.Type.ToThriftTType().String() {
	case "STRUCT":
		json, ok := json.(map[string]interface{})
		if !ok {
			return fmt.Errorf("type mismatch, expected struct, got %T for field %s", json, fieldName)
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
			return fmt.Errorf("type mismatch, expected bool, got %T for field %s", json, fieldName)
		}
	case "BYTE":
		_, ok := json.(byte)
		if ok {
			return nil
		} else {
			return fmt.Errorf("type mismatch, expected byte, got %T for field %s", json, fieldName)
		}
	case "DOUBLE":
		_, ok := json.(float64)
		if ok {
			return nil
		} else {
			return fmt.Errorf("type mismatch, expected number, got %T for field %s", json, fieldName)
		}
	case "I16":
		_, ok := json.(float64)
		if ok {
			return nil
		} else {
			return fmt.Errorf("type mismatch, expected number, got %T for field %s", json, fieldName)
		}
	case "I32":
		_, ok := json.(float64)
		if ok {
			return nil
		} else {
			return fmt.Errorf("type mismatch, expected number, got %T for field %s", json, fieldName)
		}
	case "I64":
		_, ok := json.(float64)
		if ok {
			return nil
		} else {
			return fmt.Errorf("type mismatch, expected number, got %T for field %s", json, fieldName)
		}
	case "STRING":
		_, ok := json.(string)
		if ok {
			return nil
		} else {
			return fmt.Errorf("type mismatch, expected string, got %T for field %s", json, fieldName)
		}
	case "MAP":
		keyType := t.Key
		valueType := t.Elem
		json, _ := json.(map[interface{}]interface{})
		for key, value := range json {
			err := validateType(keyType, key, fieldName+" key")
			if err != nil {
				return err
			}
			err = validateType(valueType, value, fieldName+" value")
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
	return validateType(p.Type, v, p.Name)
}

func validateBody(fuc *descriptor.FunctionDescriptor, json map[string]interface{}) error {
	params := fuc.Request.Struct.FieldsByID
	if params[2] != nil {
		// invalid number of fields, only parameter is request
		// benefit of doubt
		return nil
	}
	return validateType(params[1].Type, json, "request")
}

func treatJsonBody(c *app.RequestContext) (map[string]interface{}, error) {
	b, _ := c.Body()
	var j map[string]interface{}

	err := json.Unmarshal(b, &j)
	if err != nil {
		log.Println("error unmarshalling json body, here")
		return nil, err
	} else {
		log.Println("json body unmarshalled successfully")
	}
	return j, nil
}
