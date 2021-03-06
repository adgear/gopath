package path

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func JsonSchema(typ reflect.Type) string {
	mp := map[string]interface{}{}
	js := jsonSchema(typ, mp)
	if j, err := json.MarshalIndent(js, "", "    "); err != nil {
		return fmt.Sprint(err)
	} else {
		return strings.Replace(string(j), "\\n", "\n", -1)
	}
}

func jsonSchema(typ reflect.Type, mp map[string]interface{}) interface{} {

	switch typ.Kind() {

	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)

			var name string
			list := strings.Split(field.Tag.Get("json"), ",")
			if len(list) > 0 && list[0] != "" {
				name = list[0]
			} else {
				name = field.Name
			}

			mp[name] = jsonSchema(field.Type, map[string]interface{}{})
		}

	case reflect.Ptr:
		return jsonSchema(typ.Elem(), map[string]interface{}{})

	case reflect.Slice:
		return []interface{}{jsonSchema(typ.Elem(), map[string]interface{}{})}

	case reflect.Map:
		m := map[string]interface{}{}
		m[typ.Key().String()] = jsonSchema(typ.Elem(), map[string]interface{}{})
		return m

	default:
		return typ.Kind().String()
	}
	return mp
}
