package temp

import "net/http"

const (
	RelativePathToIDL = "../IDL/"
	RefaultHttpMethod = http.MethodPost
	RefaultRoute      = "/:serviceName/:methodName/*furtherRoutes"
)

func validateType (t *descriptor.TypeDescriptor, json map[string]interface) bool {
	switch t.Type {
	case descriptor.TypeStruct:
		params := t.Struct.FieldsByID
		for _, p := range paramLs {
			if p.Optional || p.DefaultValue != nil {
				continue
			}
			var v interface{}
			if json[p.Name] == nil {
				if json[p.Alias] == nil {
					// invalid number of fields
				} else {
					v = json[p.Alias]
				}		
			} else {
				v = json[p.Name]
			}
			ok = validateType(p.Type, v)
			if !ok { return false}

			case descriptor.TypeI64:
}

func validateName (p *descriptor.FieldDescriptor, jsonField interface) bool {
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