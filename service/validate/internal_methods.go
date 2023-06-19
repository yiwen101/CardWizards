package validate

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	desc "github.com/yiwen101/CardWizards/service/descriptor"
)

func validateType(t *descriptor.TypeDescriptor, json interface{}, fieldName string) error {
	// debug log: it seems that all number encoded in json format are float64

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

func isGenericRoute(serviceName, methodName string) error {
	return desc.DescsManager.ValidateServiceAndMethodName(serviceName, methodName)
}

func isAnnotatedRoute(req *descriptor.HTTPRequest) (string, error) {
	return desc.DescsManager.ValidateServiceAndMethodNameWithAnnotedRoutes(req)
}

func treatJsonBody(ctx context.Context, c *app.RequestContext) (map[string]interface{}, error) {
	/*
		if string(c.ContentType()) != "application/json" {
			return nil, fmt.Errorf("Invalid Content-Type, expected application/json")
		} */

	b, _ := c.Body()
	log.Println(b)
	//body2 := c.Request.BodyBytes()

	log.Println("gateway here")
	/*
		log.Println("c.Body() is:")
		log.Println(b)
		log.Println("c.ContentType is:")
		log.Println(c.ContentType())
		log.Println("c.Request.BodyBytes() is:")
		log.Println(body2)
		log.Println("c.Params")
		log.Println(c.Params)
		//log.Println("c.Request")
		//log.Println(c.Request)
		log.Println("c.Request.Body")
		log.Println(c.Request.Body())
		log.Println("string(c.Request.Body)")
		log.Println(string(c.Request.Body()))
		log.Println("c.Params")
		log.Println(c.Params)
		log.Printf("c.Keys")
		log.Println(c.Keys)
	*/
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
