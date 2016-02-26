package model

import (
	//"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/client"

	"reflect"
	"strconv"
	"strings"
)

const (
	schemaTagName = "schema"
)

func AddType(s *client.Schemas, schemaName string, obj interface{}) *client.Schema {
	t := reflect.TypeOf(obj)
	schema := client.Schema{
		Resource: client.Resource{
			Id:    schemaName,
			Type:  "schema",
			Links: map[string]string{},
		},
		PluralName:        guessPluralName(schemaName),
		ResourceFields:    typeToFields(t),
		CollectionMethods: []string{"GET"},
		ResourceMethods:   []string{"GET"},
	}

	if s.Data == nil {
		s.Data = []client.Schema{}
	}

	s.Data = append(s.Data, schema)

	return &s.Data[len(s.Data)-1]
}

func guessPluralName(name string) string {
	if name == "" {
		return ""
	}

	if strings.HasSuffix(name, "s") ||
		strings.HasSuffix(name, "ch") ||
		strings.HasSuffix(name, "x") {
		return name + "es"
	}
	return name + "s"
}

func parseTagString(tag string) string {
	splitted := strings.Split(tag, "=")
	if len(splitted) == 2 {
		return splitted[1]
	}
	return ""
}

func parseTagBool(tag string) bool {
	splitted := strings.Split(tag, "=")
	if len(splitted) == 2 {
		if val, err := strconv.ParseBool(splitted[1]); err == nil {
			return val
		}
	}
	return false
}

func parseTagInterface(tag string) interface{} {
	splitted := strings.Split(tag, "=")
	if len(splitted) == 2 {
		return splitted[1]
	}
	return nil
}

func typeToFields(t reflect.Type) map[string]client.Field {
	result := map[string]client.Field{}

	for i := 0; i < t.NumField(); i++ {
		schemaField := client.Field{}

		typeField := t.Field(i)
		if typeField.Anonymous && typeField.Type.Kind() == reflect.Struct {
			parentFields := typeToFields(typeField.Type)
			for k, v := range result {
				parentFields[k] = v
			}
			result = parentFields
			continue
		} else if typeField.Anonymous {
			continue
		}

		create := false
		update := false
		nullable := false
		required := false
		schemaType := ""
		var defaultValue interface{}
		tagParts := strings.Split(typeField.Tag.Get(schemaTagName), ",")
		for _, tag := range tagParts[0:] {
			switch {
			case strings.HasPrefix(tag, "create"):
				create = parseTagBool(tag)
			case strings.HasPrefix(tag, "update"):
				update = parseTagBool(tag)
			case strings.HasPrefix(tag, "nullable"):
				nullable = parseTagBool(tag)
			case strings.HasPrefix(tag, "required"):
				required = parseTagBool(tag)
			case strings.HasPrefix(tag, "type"):
				schemaType = parseTagString(tag)
			case strings.HasPrefix(tag, "default"):
				defaultValue = parseTagInterface(tag)
			}
		}
		schemaField.Create = create
		schemaField.Update = update
		schemaField.Nullable = nullable
		schemaField.Required = required
		schemaField.Default = defaultValue

		if schemaType != "" {
			schemaField.Type = schemaType
		} else {
			fieldString := strings.ToLower(typeField.Type.Kind().String())

			switch {
			case strings.HasPrefix(fieldString, "int") || strings.HasPrefix(fieldString, "uint"):
				schemaField.Type = "int"
			case fieldString == "bool":
				schemaField.Type = "boolean"
			case fieldString == "float32" || fieldString == "float64":
				schemaField.Type = "float"
			case fieldString == "string":
				schemaField.Type = "string"
			case fieldString == "map":
				schemaField.Type = "map[string]"
			case fieldString == "slice":
				schemaField.Type = "array[string]"
			}
		}

		name := strings.Split(typeField.Tag.Get("json"), ",")[0]
		if name == "" && len(typeField.Name) > 1 {
			name = strings.ToLower(typeField.Name[0:1]) + typeField.Name[1:]
		} else if name == "" {
			name = typeField.Name
		}

		if schemaField.Type != "" {
			result[name] = schemaField
		}
	}

	return result
}

func GetTransitioning(state string, trans string) string {
	if trans == "error" {
		return trans
	}
	if strings.HasSuffix(state, "ing") && strings.ToLower(state) != "running" {
		return "yes"
	}
	return "no"
}
